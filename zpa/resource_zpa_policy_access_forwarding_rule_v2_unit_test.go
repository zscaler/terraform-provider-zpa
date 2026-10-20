package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func TestResourcePolicyForwardingRuleV2_CRUD(t *testing.T) {
	spec := policyRuleSpec("CLIENT_FORWARDING_POLICY", true)
	spec.resource = resourcePolicyForwardingRuleV2()
	spec.config = map[string]any{
		"name":       "Rule CRM",
		"action":     "INTERCEPT",
		"conditions": valuesConditionV2("APP", "app1"),
	}
	spec.response = policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Rule CRM", Action: "INTERCEPT", PolicySetID: "ps-CLIENT_FORWARDING_POLICY",
		Conditions: valuesConditionResponseV2("APP", "app1"),
	}
	spec.wantBody = map[string]any{
		"action": "INTERCEPT", "policySetId": "ps-CLIENT_FORWARDING_POLICY", "conditions": valuesConditionBodyV2("APP", "app1"),
	}
	spec.wantState = map[string]any{"conditions": valuesConditionStateV2("APP", "app1")}
	testCRUD(t, spec)
}
