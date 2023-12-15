package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var policyTypeNames = []string{
	"ACCESS_POLICY", "GLOBAL_POLICY", "TIMEOUT_POLICY", "REAUTH_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY", "BYPASS_POLICY", "SIEM_POLICY",
}

func TestAccDataSourcePolicyType_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourcePolicyType_basic(),
				Check: resource.ComposeTestCheckFunc(
					generatePolicyTypeChecks()...,
				),
			},
		},
	})
}

func generatePolicyTypeChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range policyTypeNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_policy_type.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "policy_type"),
		))
	}
	return checks
}

func testAccCheckDataSourcePolicyType_basic() string {
	var configs string
	for _, name := range policyTypeNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_policy_type" "%s" {
    policy_type = "%s"
}
`, resourceName, name)
	}
	return configs
}
