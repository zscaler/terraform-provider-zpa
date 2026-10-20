package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

// privileged_portal_capabilities booleans are sent as a list of capability names.
func TestResourcePolicyPortalAccessRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("PRIVILEGED_PORTAL_POLICY", true)
	spec.resource = resourcePolicyPortalAccessRule()
	spec.config = map[string]any{
		"name":       "Portal capabilities",
		"action":     "CHECK_PRIVILEGED_PORTAL_CAPABILITIES",
		"conditions": valuesConditionV2("PRIVILEGE_PORTAL", "pp1"),
		"privileged_portal_capabilities": []any{map[string]any{
			"delete_file": true, "request_approvals": true, "upload_inspected_scan": true,
		}},
	}
	spec.response = policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Portal capabilities", Action: "CHECK_PRIVILEGED_PORTAL_CAPABILITIES", PolicySetID: "ps-PRIVILEGED_PORTAL_POLICY",
		Conditions: valuesConditionResponseV2("PRIVILEGE_PORTAL", "pp1"),
		PrivilegedPortalCapabilities: policysetcontrollerv2.PrivilegedPortalCapabilities{
			Capabilities: []string{"DELETE_FILE", "REQUEST_APPROVALS", "UPLOAD_INSPECTED_SCAN"},
		},
	}
	spec.wantBody = map[string]any{
		"action": "CHECK_PRIVILEGED_PORTAL_CAPABILITIES", "conditions": valuesConditionBodyV2("PRIVILEGE_PORTAL", "pp1"),
		"privilegedPortalCapabilities": map[string]any{
			"capabilities": []any{"DELETE_FILE", "REQUEST_APPROVALS", "UPLOAD_INSPECTED_SCAN"},
		},
	}
	spec.wantState = map[string]any{
		"privileged_portal_capabilities": []any{map[string]any{
			"delete_file": true, "request_approvals": true, "upload_inspected_scan": true, "review_approvals": false,
		}},
	}
	testCRUD(t, spec)
}
