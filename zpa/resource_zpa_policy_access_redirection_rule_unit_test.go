package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func redirectionRuleConfig() map[string]any {
	return map[string]any{
		"name":                "Redirect exporter",
		"action":              "REDIRECT_PREFERRED",
		"service_edge_groups": []any{map[string]any{"id": []any{"seg1"}}},
		"conditions":          clientTypeConditionV1(),
	}
}

// policy_set_id is optional: without it, create resolves the redirection policy
// set and posts the rule to it (policySet/ps-REDIRECTION_POLICY/rule).
func TestResourcePolicyRedirectionRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("REDIRECTION_POLICY", false)
	spec.resource = resourcePolicyRedictionRule()
	spec.config = redirectionRuleConfig()
	spec.response = policysetcontroller.PolicyRule{
		ID: "r1", Name: "Redirect exporter", Action: "REDIRECT_PREFERRED", PolicySetID: "ps-REDIRECTION_POLICY",
		ServiceEdgeGroups: []serviceedgegroup.ServiceEdgeGroup{{ID: "seg1"}}, Conditions: clientTypeConditionsResponseV1(),
	}
	spec.wantBody = map[string]any{
		"action":            "REDIRECT_PREFERRED",
		"serviceEdgeGroups": []any{map[string]any{"id": "seg1"}},
	}
	spec.wantState = map[string]any{"service_edge_groups": []any{map[string]any{"id": []any{"seg1"}}}}
	testCRUD(t, spec)
}

func TestResourcePolicyRedirectionRule_ActionAndServiceEdgeGroups(t *testing.T) {
	cases := map[string]struct {
		action  string
		groups  bool
		wantErr string
	}{
		"preferred without groups": {"REDIRECT_PREFERRED", false, "must be selected"},
		"always without groups":    {"REDIRECT_ALWAYS", false, "must be selected"},
		"default with groups":      {"REDIRECT_DEFAULT", true, "must be empty"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			config := redirectionRuleConfig()
			config["action"] = tc.action
			if !tc.groups {
				delete(config, "service_edge_groups")
			}
			d := fx.resourceData(resourcePolicyRedictionRule(), "", config)
			if msg := fx.runErr(resourcePolicyRedictionRuleCreate, d); !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want an error containing %q", msg, tc.wantErr)
			}
			if n := len(fx.api.Server.AllRequests()); n != 0 {
				t.Errorf("no API call may be made, got %d", n)
			}
		})
	}
}
