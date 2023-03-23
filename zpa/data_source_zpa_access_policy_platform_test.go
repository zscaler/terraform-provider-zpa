package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAccessPolicyPlatforms_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: (testAccCheckDataSourceAccessPolicyPlatforms_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"data.zpa_access_policy_platforms.this", ""),
				),
			},
		},
	})
}

var testAccCheckDataSourceAccessPolicyPlatforms_basic = `
data "zpa_access_policy_platforms" "this" {}
`
