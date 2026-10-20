package zpa

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

// mockAccessPolicy registers the access policy set lookup and returns its ID.
func mockAccessPolicy(fx *unitFixture) string {
	fx.on("GET", fx.path("policySet", "policyType", "ACCESS_POLICY"),
		common.SuccessResponse(policysetcontroller.PolicySet{ID: "ps1", PolicyType: "1"}))
	return "ps1"
}

func accessRuleV2Config() map[string]any {
	return map[string]any{
		"name":   "Allow CRM",
		"action": "ALLOW",
		"conditions": []any{
			map[string]any{
				"operator": "OR",
				"operands": []any{
					map[string]any{"object_type": "APP", "values": []any{"2001"}},
				},
			},
			map[string]any{
				"operator": "OR",
				"operands": []any{
					map[string]any{
						"object_type":  "SAML",
						"entry_values": []any{map[string]any{"lhs": "attr1", "rhs": "user@acme.com"}},
					},
				},
			},
		},
	}
}

// The API returns conditions in the v1 shape (one operand per value, lhs/rhs).
func accessRuleV2Response() policysetcontrollerv2.PolicyRuleResource {
	return policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Allow CRM", Action: "ALLOW", PolicySetID: "ps1",
		Conditions: []policysetcontrollerv2.PolicyRuleResourceConditions{
			{Operator: "OR", Operands: []policysetcontrollerv2.PolicyRuleResourceOperands{
				{ObjectType: "APP", LHS: "id", RHS: "2001"},
			}},
			{Operator: "OR", Operands: []policysetcontrollerv2.PolicyRuleResourceOperands{
				{ObjectType: "SAML", LHS: "attr1", RHS: "user@acme.com"},
			}},
		},
	}
}

func TestResourcePolicyAccessRuleV2_CreateRequest(t *testing.T) {
	fx := newUnitFixture(t)
	ps := mockAccessPolicy(fx)
	fx.on("POST", fx.pathV2("policySet", ps, "rule"), common.CreatedResponse(policysetcontrollerv2.PolicyRule{ID: "r1"}))
	fx.on("GET", fx.path("policySet", ps, "rule", "r1"), common.SuccessResponse(accessRuleV2Response()))

	d := fx.resourceData(resourcePolicyAccessRuleV2(), "", accessRuleV2Config())
	if diags := resourcePolicyAccessV2Create(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("create failed: %v", diags)
	}
	if d.Id() != "r1" || d.Get("policy_set_id") != ps {
		t.Errorf("id=%q policy_set_id=%v, want r1 / %s", d.Id(), d.Get("policy_set_id"), ps)
	}

	bodies := fx.requestBodies("POST", fx.pathV2("policySet", ps, "rule"))
	if len(bodies) != 1 {
		t.Fatalf("expected 1 v2 create request, got %d", len(bodies))
	}
	body := bodies[0]
	if body["action"] != "ALLOW" || body["policySetId"] != ps {
		t.Errorf("unexpected rule fields: action=%v policySetId=%v", body["action"], body["policySetId"])
	}
	// Each condition is sent in the v2 shape: APP as values, SAML as lhs/rhs entry values.
	operands := map[string]map[string]any{}
	for _, c := range body["conditions"].([]any) {
		for _, o := range c.(map[string]any)["operands"].([]any) {
			op := o.(map[string]any)
			operands[op["objectType"].(string)] = op
		}
	}
	if app := operands["APP"]; app == nil || !strings.Contains(strings.Join(toStrings(app["values"]), ","), "2001") {
		t.Errorf("APP operand not sent as values: %v", operands["APP"])
	}
	saml := operands["SAML"]
	if saml == nil {
		t.Fatalf("SAML operand missing: %v", body["conditions"])
	}
	entries, _ := saml["entryValues"].([]any)
	if len(entries) != 1 || entries[0].(map[string]any)["lhs"] != "attr1" || entries[0].(map[string]any)["rhs"] != "user@acme.com" {
		t.Errorf("SAML operand not sent as entry values: %v", saml)
	}

	// The v1-shaped response is mapped back to the v2 condition blocks in state.
	if n := d.Get("conditions").(*schema.Set).Len(); n != 2 {
		t.Errorf("state has %d conditions, want 2", n)
	}
}

func TestResourcePolicyAccessRuleV2_DuplicateObjectTypeRejectedBeforeAPI(t *testing.T) {
	fx := newUnitFixture(t)
	mockAccessPolicy(fx)
	config := accessRuleV2Config()
	config["conditions"] = []any{map[string]any{
		"operator": "OR",
		"operands": []any{
			map[string]any{"object_type": "APP", "values": []any{"2001"}},
			map[string]any{"object_type": "APP", "values": []any{"2002"}},
		},
	}}

	d := fx.resourceData(resourcePolicyAccessRuleV2(), "", config)
	diags := resourcePolicyAccessV2Create(context.Background(), d, fx.client)
	if !diags.HasError() || !strings.Contains(diags[0].Summary, "can only be specified once") {
		t.Fatalf("expected the duplicate object_type to be rejected, got %v", diags)
	}
	if len(fx.api.Server.AllRequests()) != 0 {
		t.Errorf("no API call may be made for an invalid rule, got %d", len(fx.api.Server.AllRequests()))
	}
}

func TestResourcePolicyAccessRuleV2_AppOperandWithoutValuesRejected(t *testing.T) {
	fx := newUnitFixture(t)
	ps := mockAccessPolicy(fx)
	config := accessRuleV2Config()
	config["conditions"] = []any{map[string]any{
		"operator": "OR",
		"operands": []any{map[string]any{"object_type": "APP"}},
	}}

	d := fx.resourceData(resourcePolicyAccessRuleV2(), "", config)
	diags := resourcePolicyAccessV2Create(context.Background(), d, fx.client)
	if !diags.HasError() || !strings.Contains(diags[0].Summary, "Application Segment ID must be provided") {
		t.Fatalf("expected an APP operand without values to be rejected, got %v", diags)
	}
	if n := fx.calls("POST", fx.pathV2("policySet", ps, "rule")); n != 0 {
		t.Errorf("no create request may be sent for an invalid rule, got %d", n)
	}
}

// The access policy set ID is looked up once per client and reused.
func TestResourcePolicyAccessRuleV2_PolicySetIDCached(t *testing.T) {
	fx := newUnitFixture(t)
	ps := mockAccessPolicy(fx)
	fx.on("GET", fx.path("policySet", ps, "rule", "r1"), common.SuccessResponse(accessRuleV2Response()))

	for i := 0; i < 3; i++ {
		d := fx.resourceData(resourcePolicyAccessRuleV2(), "r1", accessRuleV2Config())
		if diags := resourcePolicyAccessV2Read(skipListIndex(context.Background()), d, fx.client); diags.HasError() {
			t.Fatalf("read %d failed: %v", i, diags)
		}
	}
	if n := fx.calls("GET", fx.path("policySet", "policyType", "ACCESS_POLICY")); n != 1 {
		t.Errorf("policy set looked up %d times, want 1", n)
	}
}

func TestResourcePolicyAccessRuleV2_Delete(t *testing.T) {
	fx := newUnitFixture(t)
	ps := mockAccessPolicy(fx)
	fx.on("DELETE", fx.path("policySet", ps, "rule", "r1"), common.NoContentResponse())

	d := fx.resourceData(resourcePolicyAccessRuleV2(), "r1", accessRuleV2Config())
	if diags := resourcePolicyAccessV2Delete(context.Background(), d, fx.client); diags.HasError() {
		t.Fatalf("delete failed: %v", diags)
	}
	if n := fx.calls("DELETE", fx.path("policySet", ps, "rule", "r1")); n != 1 {
		t.Errorf("expected 1 delete request, got %d", n)
	}
}

func toStrings(v any) []string {
	items, _ := v.([]any)
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// v2 policy rules take one operand per object type with a list of values; the
// API reads them back in the v1 shape (one operand per value).
func valuesConditionV2(objectType, value string) []any {
	return []any{map[string]any{
		"operator": "OR",
		"operands": []any{map[string]any{"object_type": objectType, "values": []any{value}}},
	}}
}

func valuesConditionResponseV2(objectType, value string) []policysetcontrollerv2.PolicyRuleResourceConditions {
	return []policysetcontrollerv2.PolicyRuleResourceConditions{{
		Operator: "OR",
		Operands: []policysetcontrollerv2.PolicyRuleResourceOperands{{ObjectType: objectType, LHS: "id", RHS: value}},
	}}
}

func valuesConditionBodyV2(objectType, value string) []any {
	return []any{map[string]any{
		"operator": "OR",
		"operands": []any{map[string]any{"objectType": objectType, "values": []any{value}}},
	}}
}

func valuesConditionStateV2(objectType, value string) []any {
	return []any{map[string]any{
		"operator": "OR",
		"operands": []any{map[string]any{"object_type": objectType, "values": []any{value}}},
	}}
}

func TestResourcePolicyAccessRuleV2_CRUD(t *testing.T) {
	spec := policyRuleSpec("ACCESS_POLICY", true)
	spec.resource = resourcePolicyAccessRuleV2()
	spec.config = accessRuleV2Config()
	response := accessRuleV2Response()
	response.PolicySetID = "ps-ACCESS_POLICY"
	spec.response = response
	spec.updateConfig = func(c map[string]any) { c["description"] = "Updated" }
	spec.wantBody = map[string]any{"action": "ALLOW", "conditions": []any{
		map[string]any{"operands": []any{map[string]any{"objectType": "APP", "values": []any{"2001"}}}},
		map[string]any{"operands": []any{map[string]any{"objectType": "SAML"}}},
	}}
	spec.wantState = map[string]any{"policy_set_id": "ps-ACCESS_POLICY"}
	testCRUD(t, spec)
}

// v2 conditions are checked before the rule is created: ID-based types need
// values, attribute-based types need well-formed lhs/rhs entry values.
func TestResourcePolicyAccessRuleV2_ConditionValidation(t *testing.T) {
	values := func(v ...string) map[string]any {
		out := []any{}
		for _, s := range v {
			out = append(out, s)
		}
		return map[string]any{"values": out}
	}
	entry := func(lhs, rhs string) map[string]any {
		return map[string]any{"entry_values": []any{map[string]any{"lhs": lhs, "rhs": rhs}}}
	}
	none := map[string]any{}
	cases := []struct {
		name, objectType string
		operand          map[string]any
		wantErr          string
	}{
		{"app", "APP", values("2001"), ""},
		{"app without values", "APP", none, "Application Segment ID must be provided"},
		{"segment group without values", "APP_GROUP", none, "Segment Group ID must be provided"},
		{"machine group without values", "MACHINE_GRP", none, "Machine Group ID must be provided"},
		{"location without values", "LOCATION", none, "Location ID must be provided"},
		{"edge connector group without values", "EDGE_CONNECTOR_GROUP", none, "Edge Connector Group ID must be provided"},
		{"branch connector group without values", "BRANCH_CONNECTOR_GROUP", none, "Branch Connector Group ID must be provided"},
		{"user portal without values", "USER_PORTAL", none, "User Portal ID must be provided"},
		{"chrome posture profile without values", "CHROME_POSTURE_PROFILE", none, "Chrome Posture Profile ID must be provided"},
		{"client type", "CLIENT_TYPE", values("zpn_client_type_zapp"), ""},
		{"client type without values", "CLIENT_TYPE", none, "valid Client Types"},
		{"client type unknown", "CLIENT_TYPE", values("zpn_client_type_x"), "invalid Client Type 'zpn_client_type_x'"},
		{"platform", "PLATFORM", entry("windows", "true"), ""},
		{"platform without entries", "PLATFORM", none, "valid platform types"},
		{"platform unknown", "PLATFORM", entry("solaris", "true"), "valid platform types"},
		{"risk factor", "RISK_FACTOR_TYPE", entry("ZIA", "HIGH"), ""},
		{"risk factor without entries", "RISK_FACTOR_TYPE", none, "valid risk factor values"},
		{"risk factor bad lhs", "RISK_FACTOR_TYPE", entry("ZPA", "HIGH"), "LHS must be 'ZIA'"},
		{"risk factor bad rhs", "RISK_FACTOR_TYPE", entry("ZIA", "SEVERE"), "RHS must be one of"},
		{"posture", "POSTURE", entry("udid1", "true"), ""},
		{"posture without entries", "POSTURE", none, "valid Posture UDID"},
		{"posture without udid", "POSTURE", entry("", "true"), "LHS must be a valid Posture UDID"},
		{"posture bad rhs", "POSTURE", entry("udid1", "yes"), "'true' or 'false' for POSTURE"},
		{"trusted network", "TRUSTED_NETWORK", entry("net1", "true"), ""},
		{"trusted network without entries", "TRUSTED_NETWORK", none, "valid Network ID"},
		{"trusted network without id", "TRUSTED_NETWORK", entry("", "true"), "LHS must be a valid Network ID"},
		{"trusted network bad rhs", "TRUSTED_NETWORK", entry("net1", "yes"), "'true' or 'false' for TRUSTED_NETWORK"},
		{"country code", "COUNTRY_CODE", entry("US", "true"), ""},
		{"country code without entries", "COUNTRY_CODE", none, "valid country code"},
		{"country code invalid", "COUNTRY_CODE", entry("USA", "true"), "'USA' is not a valid ISO-3166"},
		{"country code bad rhs", "COUNTRY_CODE", entry("US", "false"), "'true' for COUNTRY_CODE"},
		{"saml", "SAML", entry("attr1", "jdoe@acme.com"), ""},
		{"saml without entries", "SAML", none, "entry_values must be provided for SAML"},
		{"saml without attribute", "SAML", entry("", "jdoe@acme.com"), "valid SAML attribute ID"},
		{"saml without value", "SAML", entry("attr1", ""), "cannot be empty for SAML"},
		{"scim", "SCIM", entry("idp1", "Engineering"), ""},
		{"scim without entries", "SCIM", none, "entry_values must be provided for SCIM"},
		{"scim without idp", "SCIM", entry("", "Engineering"), "valid IdP ID and cannot be empty for SCIM object_type"},
		{"scim without value", "SCIM", entry("idp1", ""), "cannot be empty for SCIM object_type"},
		{"scim group", "SCIM_GROUP", entry("idp1", "101"), ""},
		{"scim group without entries", "SCIM_GROUP", none, "entry_values must be provided for SCIM_GROUP"},
		{"scim group without idp", "SCIM_GROUP", entry("", "101"), "valid IdP ID and cannot be empty for SCIM_GROUP"},
		{"scim group without group", "SCIM_GROUP", entry("idp1", ""), "valid scim group ID"},
		{"chrome enterprise", "CHROME_ENTERPRISE", entry("managed", "true"), ""},
		{"chrome enterprise without entries", "CHROME_ENTERPRISE", none, "entry_values must be provided for CHROME_ENTERPRISE"},
		{"chrome enterprise bad lhs", "CHROME_ENTERPRISE", entry("unmanaged", "true"), "LHS must be 'managed'"},
		{"chrome enterprise bad rhs", "CHROME_ENTERPRISE", entry("managed", "maybe"), "'true' or 'false' for CHROME_ENTERPRISE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fx := newUnitFixture(t)
			fx.catchAll("r1", map[string]any{"id": "r1", "name": "rule"})
			operand := map[string]any{"object_type": tc.objectType}
			for k, v := range tc.operand {
				operand[k] = v
			}
			d := fx.resourceData(resourcePolicyAccessRuleV2(), "", map[string]any{
				"name": "rule", "action": "ALLOW",
				"conditions": []any{map[string]any{"operator": "OR", "operands": []any{operand}}},
			})
			msg := fx.runErr(resourcePolicyAccessV2Create, d)
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
					t.Errorf("no rule may be created for an invalid condition, got POST %s", r.Path)
				}
			}
		})
	}
}
