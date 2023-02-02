package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicyType_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourcePolicyTypeConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.access_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.global_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.timeout_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.reauth_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.client_forwarding_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.inspection_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.isolation_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.bypass_policy"),
					testAccDataSourcePolicyTypeCheck("data.zpa_policy_type.siem_policy"),
				),
			},
		},
	})
}

func testAccDataSourcePolicyTypeCheck(policy_type string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(policy_type, "id"),
		resource.TestCheckResourceAttrSet(policy_type, "policy_type"),
	)
}

var testAccCheckDataSourcePolicyTypeConfig_basic = `
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

data "zpa_policy_type" "global_policy" {
    policy_type = "GLOBAL_POLICY"
}

data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

data "zpa_policy_type" "reauth_policy" {
    policy_type = "REAUTH_POLICY"
}

data "zpa_policy_type" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

data "zpa_policy_type" "inspection_policy" {
    policy_type = "INSPECTION_POLICY"
}

data "zpa_policy_type" "isolation_policy" {
    policy_type = "ISOLATION_POLICY"
}


data "zpa_policy_type" "bypass_policy" {
    policy_type = "BYPASS_POLICY"
}

data "zpa_policy_type" "siem_policy" {
    policy_type = "SIEM_POLICY"
}
`
