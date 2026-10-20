package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestResourcePolicyTimeoutRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("TIMEOUT_POLICY", false)
	spec.resource = resourcePolicyTimeoutRule()
	spec.config = map[string]any{
		"name":                "Re-auth exporter",
		"action":              "RE_AUTH",
		"reauth_idle_timeout": "600",
		"reauth_timeout":      "172800",
		"conditions":          clientTypeConditionV1(),
	}
	spec.response = policysetcontroller.PolicyRule{
		ID: "r1", Name: "Re-auth exporter", Action: "RE_AUTH", PolicySetID: "ps-TIMEOUT_POLICY",
		ReauthIdleTimeout: "600", ReauthTimeout: "172800", Conditions: clientTypeConditionsResponseV1(),
	}
	spec.wantBody = map[string]any{
		"action": "RE_AUTH", "reauthIdleTimeout": "600", "reauthTimeout": "172800",
		"policySetId": "ps-TIMEOUT_POLICY", "conditions": clientTypeConditionBody,
	}
	spec.wantState = map[string]any{"reauth_timeout": "172800"}
	testCRUD(t, spec)
}
