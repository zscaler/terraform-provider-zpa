package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func TestResourcePolicyIsolationRuleV2_CRUD(t *testing.T) {
	spec := policyRuleSpec("ISOLATION_POLICY", true)
	spec.resource = resourcePolicyIsolationRuleV2()
	spec.config = map[string]any{
		"name":                     "Rule CRM",
		"action":                   "ISOLATE",
		"conditions":               valuesConditionV2("APP", "app1"),
		"zpn_isolation_profile_id": "iso1",
	}
	spec.response = policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Rule CRM", Action: "ISOLATE", PolicySetID: "ps-ISOLATION_POLICY",
		Conditions: valuesConditionResponseV2("APP", "app1"), ZpnIsolationProfileID: "iso1",
	}
	spec.wantBody = map[string]any{
		"action": "ISOLATE", "policySetId": "ps-ISOLATION_POLICY", "conditions": valuesConditionBodyV2("APP", "app1"),
		"zpnIsolationProfileId": "iso1",
	}
	spec.wantState = map[string]any{"conditions": valuesConditionStateV2("APP", "app1"), "zpn_isolation_profile_id": "iso1"}
	testCRUD(t, spec)
}
