package zpa

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func reorderConfig(orders ...string) map[string]any {
	rules := []any{}
	for i := 0; i+1 < len(orders); i += 2 {
		rules = append(rules, map[string]any{"id": orders[i], "order": orders[i+1]})
	}
	return map[string]any{"policy_type": "ACCESS_POLICY", "rules": rules}
}

// reorderRequest mocks the access policy rules and returns the rule IDs sent
// in the reorder request, in order.
func reorderRequest(t *testing.T, existing []policysetcontroller.PolicyRule, config map[string]any) []string {
	t.Helper()
	fx := newUnitFixture(t)
	ps := fx.mockPolicySet("ACCESS_POLICY")
	fx.on("GET", fx.path("policySet", "rules", "policyType", "ACCESS_POLICY"), common.SuccessResponse(common.ZPAList(existing)))
	fx.on("PUT", fx.path("policySet", ps, "reorder"), common.NoContentResponse())

	d := fx.resourceData(resourcePolicyAccessRuleReorder(), "", config)
	fx.run(resourcePolicyAccessReorderUpdate, d)
	if d.Id() != "ACCESS_POLICY-reorder" {
		t.Errorf("id = %q, want ACCESS_POLICY-reorder", d.Id())
	}
	for _, r := range fx.api.Server.AllRequests() {
		if r.Method == "PUT" && r.Path == fx.path("policySet", ps, "reorder") {
			var ids []string
			if err := json.Unmarshal(r.Body, &ids); err != nil {
				t.Fatalf("reorder body is not a list of rule IDs: %s", r.Body)
			}
			return ids
		}
	}
	t.Fatal("no reorder request was sent")
	return nil
}

var existingAccessRules = []policysetcontroller.PolicyRule{
	{ID: "r1", Name: "one", RuleOrder: "1"},
	{ID: "r2", Name: "two", RuleOrder: "2"},
	{ID: "r3", Name: "three", RuleOrder: "3"},
	{ID: "rd", Name: "Default_Rule", RuleOrder: "4"},
}

// Configured rules go first in their order; the others keep their relative
// order and the default rule stays last.
func TestResourcePolicyAccessRuleReorder_Order(t *testing.T) {
	got := reorderRequest(t, existingAccessRules, reorderConfig("r3", "1", "r1", "2"))
	if want := []string{"r3", "r1", "r2", "rd"}; !reflect.DeepEqual(got, want) {
		t.Errorf("reorder = %v, want %v", got, want)
	}
}

// The Zscaler Deception rule at order 1 is not managed by Terraform and stays first.
func TestResourcePolicyAccessRuleReorder_DeceptionRuleStaysFirst(t *testing.T) {
	existing := append([]policysetcontroller.PolicyRule{{ID: "dec", Name: "Zscaler Deception", RuleOrder: "1"}}, existingAccessRules...)
	got := reorderRequest(t, existing, reorderConfig("r3", "1", "r1", "2"))
	if want := []string{"dec", "r3", "r1", "r2", "rd"}; !reflect.DeepEqual(got, want) {
		t.Errorf("reorder = %v, want %v", got, want)
	}
}

func TestResourcePolicyAccessRuleReorder_DuplicateOrderRejected(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("policySet", "rules", "policyType", "ACCESS_POLICY"), common.SuccessResponse(common.ZPAList(existingAccessRules)))

	d := fx.resourceData(resourcePolicyAccessRuleReorder(), "", reorderConfig("r1", "1", "r2", "1"))
	if msg := fx.runErr(resourcePolicyAccessReorderUpdate, d); !strings.Contains(msg, "duplicate order '1'") {
		t.Fatalf("expected duplicate orders to be rejected, got %q", msg)
	}
	for _, r := range fx.api.Server.AllRequests() {
		if r.Method == "PUT" {
			t.Errorf("no reorder request may be sent, got %s %s", r.Method, r.Path)
		}
	}
}

// Rule order is not reset on destroy: delete only removes the resource from state.
func TestResourcePolicyAccessRuleReorder_DeleteMakesNoAPICall(t *testing.T) {
	fx := newUnitFixture(t)
	d := fx.resourceData(resourcePolicyAccessRuleReorder(), "ACCESS_POLICY-reorder", reorderConfig("r1", "1"))
	fx.run(resourcePolicyAccessReorderDelete, d)
	if n := len(fx.api.Server.AllRequests()); n != 0 {
		t.Errorf("expected no API calls, got %d", n)
	}
}
