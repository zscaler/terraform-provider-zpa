package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIdpController_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceIdpController_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceIdpControllerCheck("data.zpa_idp_controller.bd_okta_users"),
					testAccDataSourceIdpControllerCheck("data.zpa_idp_controller.bd_okta_admin"),
				),
			},
		},
	})
}

func testAccDataSourceIdpControllerCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceIdpController_basic = `
data "zpa_idp_controller" "bd_okta_users" {
    name = "BD_Okta_Users"
}

data "zpa_idp_controller" "bd_okta_admin" {
    name = "BD_Okta_Admin"
}
`
