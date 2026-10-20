package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func TestResourcePolicyInspectionRuleV2_CRUD(t *testing.T) {
	spec := policyRuleSpec("INSPECTION_POLICY", true)
	spec.resource = resourcePolicyInspectionRuleV2()
	spec.config = map[string]any{
		"name":                      "Rule CRM",
		"action":                    "INSPECT",
		"conditions":                valuesConditionV2("APP", "app1"),
		"zpn_inspection_profile_id": "ip1",
	}
	spec.response = policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Rule CRM", Action: "INSPECT", PolicySetID: "ps-INSPECTION_POLICY",
		Conditions: valuesConditionResponseV2("APP", "app1"), ZpnInspectionProfileID: "ip1",
	}
	spec.wantBody = map[string]any{
		"action": "INSPECT", "policySetId": "ps-INSPECTION_POLICY", "conditions": valuesConditionBodyV2("APP", "app1"),
		"zpnInspectionProfileId": "ip1",
	}
	spec.wantState = map[string]any{"conditions": valuesConditionStateV2("APP", "app1"), "zpn_inspection_profile_id": "ip1"}
	testCRUD(t, spec)
}
