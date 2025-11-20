package datasources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

var policyTypeNames = []string{
	"ACCESS_POLICY", "GLOBAL_POLICY", "TIMEOUT_POLICY", "REAUTH_POLICY", "SIEM_POLICY",
	"CLIENT_FORWARDING_POLICY", "BYPASS_POLICY", "INSPECTION_POLICY", "CREDENTIAL_POLICY",
	"CAPABILITIES_POLICY", "ISOLATION_POLICY", "REDIRECTION_POLICY",
}

func TestAccDataSourcePolicyType_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
