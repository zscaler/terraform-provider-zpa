package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

func lssConfig(sourceLogType string) map[string]any {
	return map[string]any{
		"config": []any{map[string]any{
			"name":            "SIEM user activity",
			"enabled":         true,
			"format":          "json",
			"lss_host":        "splunk.acme.com",
			"lss_port":        "5001",
			"source_log_type": sourceLogType,
			"use_tls":         true,
		}},
		"connector_groups": []any{map[string]any{"id": []any{"acg1"}}},
		"policy_rule_resource": []any{map[string]any{
			"name":   "policy_rule_resource-lss",
			"action": "ALLOW",
			"conditions": []any{map[string]any{
				"operator": "OR",
				"operands": []any{map[string]any{
					"object_type": "CLIENT_TYPE",
					"values":      []any{"zpn_client_type_exporter"},
				}},
			}},
		}},
	}
}

func TestResourceLSSConfigController_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceLSSConfigController(),
		config:   lssConfig("zpn_trans_log"),
		id:       "lss1",
		response: lssconfigcontroller.LSSResource{
			ID: "lss1",
			LSSConfig: &lssconfigcontroller.LSSConfig{
				ID: "lss1", Name: "SIEM user activity", Enabled: true, Format: "json",
				LSSHost: "splunk.acme.com", LSSPort: "5001", SourceLogType: "zpn_trans_log", UseTLS: true,
			},
			ConnectorGroups: []lssconfigcontroller.ConnectorGroups{{ID: "acg1"}},
		},
		createPath: v2Path("lssConfig"),
		itemPath:   v2Path("lssConfig", "lss1"),
		wantBody: map[string]any{
			"config": map[string]any{
				"lssHost": "splunk.acme.com", "lssPort": "5001", "sourceLogType": "zpn_trans_log", "useTls": true,
			},
			"connectorGroups": []any{map[string]any{"id": "acg1"}},
			"policyRuleResource": map[string]any{
				"conditions": []any{map[string]any{"operands": []any{map[string]any{
					"objectType": "CLIENT_TYPE", "values": []any{"zpn_client_type_exporter"},
				}}}},
			},
		},
		wantState: map[string]any{
			"config":           []any{map[string]any{"lss_host": "splunk.acme.com", "source_log_type": "zpn_trans_log"}},
			"connector_groups": []any{map[string]any{"id": []any{"acg1"}}},
		},
	})
}

func TestResourceLSSConfigController_InvalidConfigRejected(t *testing.T) {
	cases := map[string]struct {
		mutate  func(c map[string]any)
		wantErr string
	}{
		"filter on a log type without filters": {
			func(c map[string]any) {
				cfg := c["config"].([]any)[0].(map[string]any)
				cfg["source_log_type"] = "zpn_audit_log"
				cfg["filter"] = []any{"ZPN_STATUS_AUTHENTICATED"}
				delete(c, "policy_rule_resource")
			},
			"filter is not supported for source log type zpn_audit_log",
		},
		"unknown client type": {
			func(c map[string]any) {
				operand := c["policy_rule_resource"].([]any)[0].(map[string]any)["conditions"].([]any)[0].(map[string]any)["operands"].([]any)[0].(map[string]any)
				operand["values"] = []any{"zpn_client_type_unknown"}
			},
			"invalid value: zpn_client_type_unknown for object_type CLIENT_TYPE",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			config := lssConfig("zpn_trans_log")
			tc.mutate(config)
			d := fx.resourceData(resourceLSSConfigController(), "", config)
			if msg := fx.runErr(resourceLSSConfigControllerCreate, d); !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want an error containing %q", msg, tc.wantErr)
			}
			if n := len(fx.api.Server.AllRequests()); n != 0 {
				t.Errorf("no API call may be made, got %d", n)
			}
		})
	}
}

// Filters and policy operands are checked against the source log type before the
// LSS configuration is created.
func TestResourceLSSConfigController_FilterAndOperandValidation(t *testing.T) {
	type operand = map[string]any
	cases := []struct {
		name          string
		sourceLogType string
		filters       []any
		operand       operand // nil: no policy rule
		wantErr       string
	}{
		{"user status filter", "zpn_auth_log", []any{"ZPN_STATUS_AUTHENTICATED"}, nil, ""},
		{"user status filter invalid", "zpn_auth_log", []any{"CLT_INVALID_DOMAIN"}, nil, "invalid filter: CLT_INVALID_DOMAIN"},
		{"user activity filter", "zpn_trans_log", []any{"CLT_INVALID_DOMAIN"}, nil, ""},
		{"app connector status filter", "zpn_ast_auth_log", []any{"ZPN_STATUS_DISCONNECTED"}, nil, ""},
		{"service edge status filter", "zpn_sys_auth_log", []any{"ZPN_STATUS_AUTH_FAILED"}, nil, ""},
		{"filter on unknown log type", "zpn_custom_log", []any{"ZPN_STATUS_AUTHENTICATED"}, nil, "invalid source_log_type"},
		{"saml operand", "zpn_auth_log", nil, operand{"object_type": "SAML", "entry_values": []any{map[string]any{"lhs": "attr1", "rhs": "jdoe@acme.com"}}}, ""},
		{"saml operand without entry values", "zpn_auth_log", nil, operand{"object_type": "SAML"}, "entry_values must be provided for object_type SAML"},
		{"saml operand without rhs", "zpn_trans_log", nil, operand{"object_type": "SCIM", "entry_values": []any{map[string]any{"lhs": "attr1", "rhs": ""}}}, "both lhs and rhs must be provided"},
		{"client type with entry values", "zpn_trans_log", nil, operand{"object_type": "CLIENT_TYPE", "values": []any{"zpn_client_type_exporter"}, "entry_values": []any{map[string]any{"lhs": "a", "rhs": "b"}}}, "entry_values should not be provided for object_type CLIENT_TYPE"},
		{"operand on log type without operands", "zpn_ast_auth_log", nil, operand{"object_type": "CLIENT_TYPE", "values": []any{"zpn_client_type_exporter"}}, "invalid source_log_type: zpn_ast_auth_log"},
		{"object type not allowed for log type", "zpn_auth_log", nil, operand{"object_type": "APP", "values": []any{"2001"}}, "invalid object_type: APP for source_log_type: zpn_auth_log"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fx := newUnitFixture(t)
			fx.catchAll("lss1", lssconfigcontroller.LSSResource{ID: "lss1", LSSConfig: &lssconfigcontroller.LSSConfig{ID: "lss1", Name: "lss"}})
			config := lssConfig(tc.sourceLogType)
			delete(config, "policy_rule_resource")
			if tc.filters != nil {
				config["config"].([]any)[0].(map[string]any)["filter"] = tc.filters
			}
			if tc.operand != nil {
				config["policy_rule_resource"] = []any{map[string]any{
					"name": "lss rule", "action": "ALLOW",
					"conditions": []any{map[string]any{"operator": "OR", "operands": []any{tc.operand}}},
				}}
			}
			d := fx.resourceData(resourceLSSConfigController(), "", config)
			msg := fx.runErr(resourceLSSConfigControllerCreate, d)
			if tc.wantErr == "" && msg != "" || tc.wantErr != "" && !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want %q", msg, tc.wantErr)
			}
		})
	}
}
