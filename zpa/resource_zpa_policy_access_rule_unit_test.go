package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

// v1 policy rules take conditions as one operand per value (lhs/rhs).
func clientTypeConditionV1() []any {
	return []any{map[string]any{
		"operator": "OR",
		"operands": []any{map[string]any{"object_type": "CLIENT_TYPE", "lhs": "id", "rhs": "zpn_client_type_exporter"}},
	}}
}

func clientTypeConditionsResponseV1() []policysetcontroller.Conditions {
	return []policysetcontroller.Conditions{{
		Operator: "OR",
		Operands: []policysetcontroller.Operands{{ObjectType: "CLIENT_TYPE", LHS: "id", RHS: "zpn_client_type_exporter"}},
	}}
}

var clientTypeConditionBody = []any{map[string]any{
	"operator": "OR",
	"operands": []any{map[string]any{"objectType": "CLIENT_TYPE", "lhs": "id", "rhs": "zpn_client_type_exporter"}},
}}

func TestResourcePolicyAccessRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("ACCESS_POLICY", false)
	spec.resource = resourcePolicyAccessRule()
	spec.config = map[string]any{
		"name":                 "Allow exporter",
		"action":               "ALLOW",
		"operator":             "AND",
		"conditions":           clientTypeConditionV1(),
		"app_connector_groups": []any{map[string]any{"id": []any{"acg1"}}},
	}
	spec.response = policysetcontroller.PolicyRule{
		ID: "r1", Name: "Allow exporter", Action: "ALLOW", Operator: "AND", PolicySetID: "ps-ACCESS_POLICY",
		Conditions:         clientTypeConditionsResponseV1(),
		AppConnectorGroups: []appconnectorgroup.AppConnectorGroup{{ID: "acg1"}},
	}
	spec.wantBody = map[string]any{
		"name": "Allow exporter", "action": "ALLOW", "policySetId": "ps-ACCESS_POLICY",
		"conditions":         clientTypeConditionBody,
		"appConnectorGroups": []any{map[string]any{"id": "acg1"}},
	}
	spec.wantState = map[string]any{
		"policy_set_id":        "ps-ACCESS_POLICY",
		"app_connector_groups": []any{map[string]any{"id": []any{"acg1"}}},
	}
	testCRUD(t, spec)
}

// An APP operand must reference an existing application segment; the rule is not created otherwise.
func TestResourcePolicyAccessRule_UnknownAppRejected(t *testing.T) {
	fx := newUnitFixture(t)
	fx.mockPolicySet("ACCESS_POLICY")
	fx.on("GET", fx.path("application", "404"), common.NotFoundResponse())

	d := fx.resourceData(resourcePolicyAccessRule(), "", map[string]any{
		"name":   "Allow app",
		"action": "ALLOW",
		"conditions": []any{map[string]any{
			"operator": "OR",
			"operands": []any{map[string]any{"object_type": "APP", "lhs": "id", "rhs": "404"}},
		}},
	})
	if msg := fx.runErr(resourcePolicyAccessCreate, d); !strings.Contains(strings.ToLower(msg), "application segment") {
		t.Fatalf("expected an unknown application segment to be rejected, got %q", msg)
	}
	if n := fx.calls("POST", fx.path("policySet", "ps-ACCESS_POLICY", "rule")); n != 0 {
		t.Errorf("no create request may be sent, got %d", n)
	}
}

// v1 operands are checked before the rule is created: the referenced object must
// exist and lhs/rhs must have the shape of the object type.
func TestResourcePolicyAccessRule_OperandValidation(t *testing.T) {
	// a tenant where every referenced object exists (as "x1")
	tenant := func(fx *unitFixture) {
		fx.catchAll("x1", map[string]any{
			"id": "x1", "name": "x1", "idpId": "x1", "postureUdid": "udid1", "networkId": "net1",
		})
	}
	empty := func(fx *unitFixture) { fx.mockPolicySet("ACCESS_POLICY") }
	// SCIM group IDs are numeric
	scimTenant := func(fx *unitFixture) {
		ps := fx.mockPolicySet("ACCESS_POLICY")
		fx.on("GET", fx.path("idp", "x1"), common.SuccessResponse(map[string]any{"id": "x1"}))
		fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimgroup", "101"), common.SuccessResponse(map[string]any{"id": 101}))
		fx.on("POST", fx.path("policySet", ps, "rule"), common.CreatedResponse(map[string]any{"id": "r1"}))
		fx.on("GET", fx.path("policySet", ps, "rule", "r1"), common.SuccessResponse(map[string]any{"id": "r1"}))
	}

	cases := []struct {
		name            string
		objectType, lhs string
		rhs, idpID      string
		setup           func(fx *unitFixture)
		wantErr         string // "" when the rule must be created
	}{
		{"app", "APP", "id", "x1", "", tenant, ""},
		{"app wrong lhs", "APP", "name", "x1", "", empty, "LHS must be"},
		{"app missing", "APP", "id", "x1", "", empty, "application segment ID"},
		{"segment group", "APP_GROUP", "id", "x1", "", tenant, ""},
		{"idp", "IDP", "id", "x1", "", tenant, ""},
		{"cloud connector group", "EDGE_CONNECTOR_GROUP", "id", "x1", "", tenant, ""},
		{"machine group", "MACHINE_GRP", "id", "x1", "", tenant, ""},
		{"client type", "CLIENT_TYPE", "id", "zpn_client_type_zapp", "", tenant, ""},
		{"client type unknown", "CLIENT_TYPE", "id", "zpn_client_type_x", "", empty, "RHS values must be one of"},
		{"posture", "POSTURE", "udid1", "true", "", tenant, ""},
		{"posture without udid", "POSTURE", "", "true", "", empty, "valid posture network ID"},
		{"posture unknown udid", "POSTURE", "udid9", "true", "", tenant, "valid posture network ID"},
		{"posture bad rhs", "POSTURE", "udid1", "yes", "", tenant, "RHS must be"},
		{"trusted network", "TRUSTED_NETWORK", "net1", "true", "", tenant, ""},
		{"trusted network without id", "TRUSTED_NETWORK", "", "true", "", empty, "valid trusted network ID"},
		{"trusted network bad rhs", "TRUSTED_NETWORK", "net1", "false", "", tenant, "RHS must be"},
		{"platform", "PLATFORM", "linux", "true", "", tenant, ""},
		{"platform without lhs", "PLATFORM", "", "true", "", empty, "valid platform ID"},
		{"platform bad rhs", "PLATFORM", "linux", "false", "", tenant, "RHS must be"},
		{"saml", "SAML", "x1", "jdoe@acme.com", "", tenant, ""},
		{"saml without attribute", "SAML", "", "jdoe@acme.com", "", empty, "valid SAML Attribute ID"},
		{"saml without value", "SAML", "x1", "", "", tenant, "no RHS is provided"},
		{"scim without idp", "SCIM", "x1", "Engineering", "", empty, "Idp ID must be set"},
		{"scim without attribute", "SCIM", "", "Engineering", "x1", empty, "valid SCIM Attribute ID"},
		{"scim group", "SCIM_GROUP", "x1", "101", "", scimTenant, ""},
		{"scim group unknown", "SCIM_GROUP", "x1", "999", "", scimTenant, "SCIM Group ID"},
		{"scim group without idp", "SCIM_GROUP", "", "x1", "", empty, "valid IDP Controller ID"},
		{"scim group without group", "SCIM_GROUP", "x1", "", "", tenant, "no RHS is provided"},
		{"country code", "COUNTRY_CODE", "US", "true", "", tenant, ""},
		{"country code invalid", "COUNTRY_CODE", "USA", "true", "", empty, "ISO-3166 Alpha-2"},
		{"risk factor", "RISK_FACTOR_TYPE", "ZIA", "HIGH", "", tenant, ""},
		{"risk factor bad lhs", "RISK_FACTOR_TYPE", "ZPA", "HIGH", "", empty, "LHS must be"},
		{"risk factor bad rhs", "RISK_FACTOR_TYPE", "ZIA", "SEVERE", "", empty, "RHS must be"},
		{"chrome enterprise", "CHROME_ENTERPRISE", "managed", "true", "", tenant, ""},
		{"chrome enterprise bad lhs", "CHROME_ENTERPRISE", "unmanaged", "true", "", empty, "LHS must be"},
		{"chrome enterprise bad rhs", "CHROME_ENTERPRISE", "managed", "maybe", "", empty, "RHS must be"},
		{"unknown object type", "NOT_A_TYPE", "id", "x1", "", empty, "invalid operand object type"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fx := newUnitFixture(t)
			tc.setup(fx)
			operand := map[string]any{"object_type": tc.objectType, "lhs": tc.lhs, "rhs": tc.rhs}
			if tc.idpID != "" {
				operand["idp_id"] = tc.idpID
			}
			d := fx.resourceData(resourcePolicyAccessRule(), "", map[string]any{
				"name": "rule", "action": "ALLOW",
				"conditions": []any{map[string]any{"operator": "OR", "operands": []any{operand}}},
			})
			msg := fx.runErr(resourcePolicyAccessCreate, d)
			if tc.wantErr == "" {
				if msg != "" {
					t.Fatalf("expected the rule to be created, got %q", msg)
				}
				return
			}
			if !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want an error containing %q", msg, tc.wantErr)
			}
			for _, r := range fx.api.Server.AllRequests() {
				if r.Method == "POST" {
					t.Errorf("no rule may be created for an invalid operand, got POST %s", r.Path)
				}
			}
		})
	}
}

// SCIM operands must use a value the IdP actually has (case-insensitive).
func TestResourcePolicyAccessRule_ScimOperandValues(t *testing.T) {
	for name, tc := range map[string]struct {
		rhs     string
		wantErr string
	}{
		"known value, other case": {"engineering", ""},
		"unknown value":           {"Marketing", "valid SCIM Attribute Value"},
	} {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			ps := fx.mockPolicySet("ACCESS_POLICY")
			fx.on("GET", fx.path("idp", "idp1", "scimattribute", "attr1"), common.SuccessResponse(map[string]any{"id": "attr1", "idpId": "idp1"}))
			fx.on("GET", common.ZPAUserConfigPath(common.TestCustomerID, "scimattribute", "idpId", "idp1", "attributeId", "attr1"),
				common.SuccessResponse(common.ZPAList([]string{"Engineering", "Sales"})))
			fx.on("POST", fx.path("policySet", ps, "rule"), common.CreatedResponse(map[string]any{"id": "r1"}))
			fx.on("GET", fx.path("policySet", ps, "rule", "r1"), common.SuccessResponse(map[string]any{"id": "r1"}))

			d := fx.resourceData(resourcePolicyAccessRule(), "", map[string]any{
				"name": "rule", "action": "ALLOW",
				"conditions": []any{map[string]any{"operator": "OR", "operands": []any{
					map[string]any{"object_type": "SCIM", "lhs": "attr1", "rhs": tc.rhs, "idp_id": "idp1"},
				}}},
			})
			msg := fx.runErr(resourcePolicyAccessCreate, d)
			if tc.wantErr == "" && msg != "" || tc.wantErr != "" && !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want %q", msg, tc.wantErr)
			}
		})
	}
}
