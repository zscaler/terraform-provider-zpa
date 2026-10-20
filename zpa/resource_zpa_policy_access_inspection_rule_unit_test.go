package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestResourcePolicyInspectionRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("INSPECTION_POLICY", false)
	spec.resource = resourcePolicyInspectionRule()
	spec.config = map[string]any{
		"name":                      "Inspect exporter",
		"action":                    "INSPECT",
		"zpn_inspection_profile_id": "ip1",
		"conditions":                clientTypeConditionV1(),
	}
	spec.response = policysetcontroller.PolicyRule{
		ID: "r1", Name: "Inspect exporter", Action: "INSPECT", PolicySetID: "ps-INSPECTION_POLICY",
		ZpnInspectionProfileID: "ip1", Conditions: clientTypeConditionsResponseV1(),
	}
	spec.wantBody = map[string]any{
		"action": "INSPECT", "zpnInspectionProfileId": "ip1", "conditions": clientTypeConditionBody,
	}
	spec.wantState = map[string]any{"zpn_inspection_profile_id": "ip1"}
	testCRUD(t, spec)
}
