package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

// privileged_capabilities booleans are sent as a list of capability names.
func TestResourcePolicyCapabilitiesAccessRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("CAPABILITIES_POLICY", true)
	spec.resource = resourcePolicyCapabilitiesAccessRule()
	spec.config = map[string]any{
		"name":       "PRA capabilities",
		"action":     "CHECK_CAPABILITIES",
		"conditions": valuesConditionV2("APP", "app1"),
		"privileged_capabilities": []any{map[string]any{
			"clipboard_copy": true, "file_upload": true, "record_session": true,
		}},
	}
	spec.response = policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "PRA capabilities", Action: "CHECK_CAPABILITIES", PolicySetID: "ps-CAPABILITIES_POLICY",
		Conditions: valuesConditionResponseV2("APP", "app1"),
		PrivilegedCapabilities: policysetcontrollerv2.PrivilegedCapabilities{
			Capabilities: []string{"CLIPBOARD_COPY", "FILE_UPLOAD", "RECORD_SESSION"},
		},
	}
	spec.wantBody = map[string]any{
		"action": "CHECK_CAPABILITIES", "conditions": valuesConditionBodyV2("APP", "app1"),
		"privilegedCapabilities": map[string]any{
			"capabilities": []any{"CLIPBOARD_COPY", "FILE_UPLOAD", "RECORD_SESSION"},
		},
	}
	spec.wantState = map[string]any{
		"privileged_capabilities": []any{map[string]any{
			"clipboard_copy": true, "file_upload": true, "record_session": true, "file_download": false,
		}},
	}
	testCRUD(t, spec)
}
