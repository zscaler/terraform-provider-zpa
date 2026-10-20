package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestResourcePolicyForwardingRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("CLIENT_FORWARDING_POLICY", false)
	spec.resource = resourcePolicyForwardingRule()
	spec.config = map[string]any{
		"name":       "Bypass exporter",
		"action":     "BYPASS",
		"conditions": clientTypeConditionV1(),
	}
	spec.response = policysetcontroller.PolicyRule{
		ID: "r1", Name: "Bypass exporter", Action: "BYPASS", PolicySetID: "ps-CLIENT_FORWARDING_POLICY",
		Conditions: clientTypeConditionsResponseV1(),
	}
	spec.wantBody = map[string]any{
		"action": "BYPASS", "policySetId": "ps-CLIENT_FORWARDING_POLICY", "conditions": clientTypeConditionBody,
	}
	spec.wantState = map[string]any{"action": "BYPASS"}
	testCRUD(t, spec)
}
