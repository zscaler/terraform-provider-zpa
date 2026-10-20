package zpa

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func timeoutRuleV2Config(idle, timeout string) map[string]any {
	return map[string]any{
		"name":                "Re-auth CRM",
		"action":              "RE_AUTH",
		"reauth_idle_timeout": idle,
		"reauth_timeout":      timeout,
		"conditions":          valuesConditionV2("APP", "app1"),
	}
}

// Human-readable timeouts are sent in seconds and read back human-readable.
func TestResourcePolicyTimeoutRuleV2_CRUD(t *testing.T) {
	spec := policyRuleSpec("TIMEOUT_POLICY", true)
	spec.resource = resourcePolicyTimeoutRuleV2()
	spec.config = timeoutRuleV2Config("10 Minutes", "2 Days")
	spec.response = policysetcontrollerv2.PolicyRuleResource{
		ID: "r1", Name: "Re-auth CRM", Action: "RE_AUTH", PolicySetID: "ps-TIMEOUT_POLICY",
		ReauthIdleTimeout: "600", ReauthTimeout: "172800", Conditions: valuesConditionResponseV2("APP", "app1"),
	}
	spec.wantBody = map[string]any{
		"action": "RE_AUTH", "reauthIdleTimeout": "600", "reauthTimeout": "172800",
		"conditions": valuesConditionBodyV2("APP", "app1"),
	}
	spec.wantState = map[string]any{
		"reauth_idle_timeout": "10 Minutes", "reauth_timeout": "2 Days",
		"conditions": valuesConditionStateV2("APP", "app1"),
	}
	testCRUD(t, spec)
}

func TestResourcePolicyTimeoutRuleV2_NeverSentAsMinusOne(t *testing.T) {
	spec := policyRuleSpec("TIMEOUT_POLICY", true)
	fx := newUnitFixture(t)
	spec.setup(fx)
	rule := policysetcontrollerv2.PolicyRuleResource{ID: "r1", Name: "Re-auth CRM", ReauthIdleTimeout: "-1", ReauthTimeout: "-1"}
	fx.on("POST", spec.createPath, common.CreatedResponse(rule))
	fx.on("GET", spec.itemPath, common.SuccessResponse(rule))

	d := fx.resourceData(resourcePolicyTimeoutRuleV2(), "", timeoutRuleV2Config("Never", "Never"))
	fx.run(resourcePolicyTimeoutRuleV2Create, d)
	assertSubset(t, "create request", fx.body("POST", spec.createPath),
		map[string]any{"reauthIdleTimeout": "-1", "reauthTimeout": "-1"})
}

func TestResourcePolicyTimeoutRuleV2_ShortTimeoutRejected(t *testing.T) {
	fx := newUnitFixture(t)
	fx.mockPolicySet("TIMEOUT_POLICY")
	d := fx.resourceData(resourcePolicyTimeoutRuleV2(), "", timeoutRuleV2Config("5 Minutes", "2 Days"))
	if msg := fx.runErr(resourcePolicyTimeoutRuleV2Create, d); !strings.Contains(msg, "at least 10 minutes") {
		t.Fatalf("expected a timeout under 10 minutes to be rejected, got %q", msg)
	}
	if n := fx.calls("POST", v2Path("policySet", "ps-TIMEOUT_POLICY", "rule")); n != 0 {
		t.Errorf("no create request may be sent, got %d", n)
	}
}

// The API reads "Never" back as "never" and "60 Minutes" as "1 Hour": equal
// durations must not show as a change in the plan.
func TestResourcePolicyTimeoutRuleV2_EquivalentTimeoutsNoDiff(t *testing.T) {
	fx := newUnitFixture(t)
	r := resourcePolicyTimeoutRuleV2()
	cases := []struct{ state, config string }{
		{"never", "Never"},
		{"1 Hour", "60 Minutes"},
		{"2 Days", "48 Hours"},
	}
	for _, c := range cases {
		state := fx.resourceData(r, "r1", timeoutRuleV2Config(c.state, c.state)).State()
		diff, err := r.Diff(context.Background(), state, terraformConfig(timeoutRuleV2Config(c.config, c.config)), fx.client)
		if err != nil {
			t.Fatalf("plan failed: %v", err)
		}
		if diff != nil && (diff.Attributes["reauth_timeout"] != nil || diff.Attributes["reauth_idle_timeout"] != nil) {
			t.Errorf("state %q and config %q are the same duration but planned a change", c.state, c.config)
		}
	}
	// a different duration is still a change
	state := fx.resourceData(r, "r1", timeoutRuleV2Config("1 Hour", "1 Hour")).State()
	diff, _ := r.Diff(context.Background(), state, terraformConfig(timeoutRuleV2Config("2 Hours", "1 Hour")), fx.client)
	if diff == nil || diff.Attributes["reauth_idle_timeout"] == nil {
		t.Error("a different timeout must still plan a change")
	}
}
