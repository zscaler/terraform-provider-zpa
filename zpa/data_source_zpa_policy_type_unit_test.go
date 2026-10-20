package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestDataSourcePolicyType(t *testing.T) {
	cases := map[string]struct {
		config     map[string]any
		policyType string
	}{
		"given type":                  {map[string]any{"policy_type": "ACCESS_POLICY"}, "ACCESS_POLICY"},
		"defaults to the global type": {map[string]any{}, "GLOBAL_POLICY"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			fx.on("GET", fx.path("policySet", "policyType", tc.policyType), common.SuccessResponse(policysetcontroller.PolicySet{
				ID: "ps1", Name: "Policy", PolicyType: "1", Enabled: true,
				Rules: []policysetcontroller.PolicyRule{{ID: "r1", Name: "Allow", Action: "ALLOW"}},
			}))

			d := fx.resourceData(dataSourcePolicyType(), "", tc.config)
			fx.run(dataSourcePolicyTypeRead, d)
			if d.Id() != "ps1" {
				t.Errorf("id = %q, want ps1", d.Id())
			}
			assertState(t, d, map[string]any{"rules": []any{map[string]any{"id": "r1", "action": "ALLOW"}}})
		})
	}
}
