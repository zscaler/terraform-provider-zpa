package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func credentialRuleConfig() map[string]any {
	return map[string]any{
		"name":       "Inject credentials",
		"action":     "INJECT_CREDENTIALS",
		"conditions": valuesConditionV2("CONSOLE", "con1"),
		"credential": []any{map[string]any{"id": "cr1"}},
	}
}

func credentialRuleResponse() policysetcontrollerv2.PolicyRuleResource {
	return policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Inject credentials", Action: "INJECT_CREDENTIALS", PolicySetID: "ps-CREDENTIAL_POLICY",
		Conditions: valuesConditionResponseV2("CONSOLE", "con1"),
		Credential: &policysetcontrollerv2.Credential{ID: "cr1"},
	}
}

func TestResourcePolicyCredentialAccessRule_CRUD(t *testing.T) {
	spec := policyRuleSpec("CREDENTIAL_POLICY", true)
	spec.resource = resourcePolicyCredentialAccessRule()
	spec.config = credentialRuleConfig()
	spec.response = credentialRuleResponse()
	spec.wantBody = map[string]any{
		"action": "INJECT_CREDENTIALS", "conditions": valuesConditionBodyV2("CONSOLE", "con1"),
		"credential": map[string]any{"id": "cr1"},
	}
	spec.wantState = map[string]any{"credential": []any{map[string]any{"id": "cr1"}}}
	testCRUD(t, spec)
}

// Without a credential block in the plan, an update keeps the credential the rule already has.
func TestResourcePolicyCredentialAccessRule_UpdateKeepsExistingCredential(t *testing.T) {
	spec := policyRuleSpec("CREDENTIAL_POLICY", true)
	fx := newUnitFixture(t)
	spec.setup(fx)
	fx.on("GET", spec.itemPath, common.SuccessResponse(credentialRuleResponse()))
	fx.on("PUT", spec.updatePath, common.NoContentResponse())

	state := credentialRuleConfig()
	delete(state, "credential")
	config := credentialRuleConfig()
	delete(config, "credential")
	config["description"] = "updated"
	d := fx.updateData(resourcePolicyCredentialAccessRule(), "r1", state, config)
	fx.run(resourcePolicyCredentialAccessRuleUpdate, d)

	assertSubset(t, "update request", fx.body("PUT", spec.updatePath), map[string]any{"credential": map[string]any{"id": "cr1"}})
}
