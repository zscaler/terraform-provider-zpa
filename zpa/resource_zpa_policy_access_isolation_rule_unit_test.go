package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestResourcePolicyIsolationRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("ISOLATION_POLICY", false)
	spec.resource = resourcePolicyIsolationRule()
	spec.config = map[string]any{
		"name":                     "Isolate exporter",
		"action":                   "ISOLATE",
		"zpn_isolation_profile_id": "iso1",
		"conditions":               clientTypeConditionV1(),
	}
	spec.response = policysetcontroller.PolicyRule{
		ID: "r1", Name: "Isolate exporter", Action: "ISOLATE", PolicySetID: "ps-ISOLATION_POLICY",
		ZpnIsolationProfileID: "iso1", Conditions: clientTypeConditionsResponseV1(),
	}
	spec.wantBody = map[string]any{
		"action": "ISOLATE", "zpnIsolationProfileId": "iso1", "conditions": clientTypeConditionBody,
	}
	spec.wantState = map[string]any{"zpn_isolation_profile_id": "iso1"}
	testCRUD(t, spec)
}
