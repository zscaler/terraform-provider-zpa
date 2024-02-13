package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScimGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceScimGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.a000"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.b000"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.c000"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.d000"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.e000"),
				),
			},
		},
	})
}

func testAccDataSourceScimGroupCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceScimGroupConfig_basic = `
data "zpa_scim_groups" "a000" {
    name = "A000"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "b000" {
    name = "B000"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "c000" {
    name = "C000"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "d000" {
    name = "D000"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "e000" {
    name = "E000"
	idp_name = "BD_Okta_Users"
}
`
