package zpa

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// A Deception rule that is itself configured is placed at its configured order.
func TestResourcePolicyAccessRuleReorder_ConfiguredDeceptionRuleUsesItsOrder(t *testing.T) {
	existing := append([]policysetcontroller.PolicyRule{{ID: "dec", Name: "Zscaler Deception", RuleOrder: "1"}}, existingAccessRules...)
	got := reorderRequest(t, existing, reorderConfig("dec", "1", "r3", "2", "r1", "3"))
	if want := []string{"dec", "r3", "r1", "r2", "rd"}; !reflect.DeepEqual(got, want) {
		t.Errorf("reorder = %v, want %v", got, want)
	}
}

// readOrders mocks the access policy rules, reads the reorder resource with
// the given configured rules and returns the orders it stored by rule ID.
func readOrders(t *testing.T, existing []policysetcontroller.PolicyRule, config map[string]any) map[string]string {
	t.Helper()
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("policySet", "rules", "policyType", "ACCESS_POLICY"), common.SuccessResponse(common.ZPAList(existing)))

	d := fx.resourceData(resourcePolicyAccessRuleReorder(), "ACCESS_POLICY-reorder", config)
	fx.run(resourcePolicyAccessReorderRead, d)
	got := map[string]string{}
	for _, r := range d.Get("rules").(*schema.Set).List() {
		rule := r.(map[string]any)
		got[rule["id"].(string)] = rule["order"].(string)
	}
	return got
}

// Without a Deception rule the live orders are stored as they are.
func TestResourcePolicyAccessRuleReorder_ReadStoresLiveOrders(t *testing.T) {
	got := readOrders(t, existingAccessRules, reorderConfig("r1", "3", "r3", "1"))
	if want := map[string]string{"r1": "1", "r3": "3"}; !reflect.DeepEqual(got, want) {
		t.Errorf("orders = %v, want %v", got, want)
	}
}

// With the Deception rule at order 1, the configured orders are written one
// position lower than declared; reading them back removes that offset again, so
// a configuration declaring orders 1..N does not show a change in every plan.
func TestResourcePolicyAccessRuleReorder_ReadReversesDeceptionOffset(t *testing.T) {
	live := []policysetcontroller.PolicyRule{
		{ID: "dec", Name: "Zscaler Deception", RuleOrder: "1"},
		{ID: "r3", Name: "three", RuleOrder: "2"},
		{ID: "r1", Name: "one", RuleOrder: "3"},
		{ID: "r2", Name: "two", RuleOrder: "4"},
		{ID: "rd", Name: "Default_Rule", RuleOrder: "5"},
	}
	got := readOrders(t, live, reorderConfig("r3", "1", "r1", "2"))
	if want := map[string]string{"r3": "1", "r1": "2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("orders = %v, want %v", got, want)
	}
}

// A configured Deception rule is not an offset: the live orders are stored as they are.
func TestResourcePolicyAccessRuleReorder_ReadConfiguredDeceptionRuleKeepsLiveOrders(t *testing.T) {
	live := []policysetcontroller.PolicyRule{
		{ID: "dec", Name: "Zscaler Deception", RuleOrder: "1"},
		{ID: "r3", Name: "three", RuleOrder: "2"},
		{ID: "rd", Name: "Default_Rule", RuleOrder: "3"},
	}
	got := readOrders(t, live, reorderConfig("dec", "1", "r3", "2"))
	if want := map[string]string{"dec": "1", "r3": "2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("orders = %v, want %v", got, want)
	}
}

// Update writes the orders with the offset and reads them back without it: the
// state after the apply holds the configured orders.
func TestResourcePolicyAccessRuleReorder_UpdateWithDeceptionRoundTrips(t *testing.T) {
	fx := newUnitFixture(t)
	ps := fx.mockPolicySet("ACCESS_POLICY")
	before := append([]policysetcontroller.PolicyRule{{ID: "dec", Name: "Zscaler Deception", RuleOrder: "1"}}, existingAccessRules...)
	after := []policysetcontroller.PolicyRule{
		{ID: "dec", Name: "Zscaler Deception", RuleOrder: "1"},
		{ID: "r3", Name: "three", RuleOrder: "2"},
		{ID: "r1", Name: "one", RuleOrder: "3"},
		{ID: "r2", Name: "two", RuleOrder: "4"},
		{ID: "rd", Name: "Default_Rule", RuleOrder: "5"},
	}
	rulesPath := fx.path("policySet", "rules", "policyType", "ACCESS_POLICY")
	fx.on("GET", rulesPath, common.SuccessResponse(common.ZPAList(before)))
	fx.api.OnFunc("PUT", fx.path("policySet", ps, "reorder"), func(*http.Request, []byte) common.MockResponse {
		fx.on("GET", rulesPath, common.SuccessResponse(common.ZPAList(after)))
		return common.NoContentResponse()
	})

	d := fx.resourceData(resourcePolicyAccessRuleReorder(), "", reorderConfig("r3", "1", "r1", "2"))
	fx.run(resourcePolicyAccessReorderUpdate, d)

	got := map[string]string{}
	for _, r := range d.Get("rules").(*schema.Set).List() {
		rule := r.(map[string]any)
		got[rule["id"].(string)] = rule["order"].(string)
	}
	if want := map[string]string{"r3": "1", "r1": "2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("state orders = %v, want %v", got, want)
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
