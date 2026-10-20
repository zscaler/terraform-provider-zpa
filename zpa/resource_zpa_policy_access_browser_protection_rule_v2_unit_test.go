package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func TestResourcePolicyBrowserProtectionRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("CLIENTLESS_SESSION_PROTECTION_POLICY", true)
	spec.resource = resourcePolicyBrowserProtectionRule()
	spec.config = map[string]any{
		"name":       "Monitor portal sessions",
		"action":     "MONITOR",
		"conditions": valuesConditionV2("USER_PORTAL", "up1"),
	}
	spec.response = policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Monitor portal sessions", Action: "MONITOR", PolicySetID: "ps-CLIENTLESS_SESSION_PROTECTION_POLICY",
		Conditions: valuesConditionResponseV2("USER_PORTAL", "up1"),
	}
	spec.updateConfig = func(c map[string]any) { c["action"] = "DO_NOT_MONITOR" }
	spec.wantBody = map[string]any{"conditions": valuesConditionBodyV2("USER_PORTAL", "up1")}
	spec.wantState = map[string]any{"action": "MONITOR", "conditions": valuesConditionStateV2("USER_PORTAL", "up1")}
	testCRUD(t, spec)
}
