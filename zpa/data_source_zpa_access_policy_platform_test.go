package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAccessPolicyPlatforms_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAccessPolicyPlatforms_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_platforms.this", "linux"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_platforms.this", "android"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_platforms.this", "windows"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_platforms.this", "ios"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_platforms.this", "mac"),
				),
			},
		},
	})
}

var testAccCheckDataSourceAccessPolicyPlatforms_basic = `
data "zpa_access_policy_platforms" "this" {}
`
