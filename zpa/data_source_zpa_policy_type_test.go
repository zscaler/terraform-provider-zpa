package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicyType_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourcePolicyTypeConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_policy_type.access_policy", "policy_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_policy_type.timeout_policy", "policy_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_policy_type.reauth_policy", "policy_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_policy_type.siem_policy", "policy_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_policy_type.client_forwarding_policy", "policy_type"),
				),
			},
		},
	})
}

const testAccCheckDataSourcePolicyTypeConfig_basic = `
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

data "zpa_policy_type" "reauth_policy" {
    policy_type = "REAUTH_POLICY"
}

data "zpa_policy_type" "siem_policy" {
    policy_type = "SIEM_POLICY"
}

data "zpa_policy_type" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}
`
