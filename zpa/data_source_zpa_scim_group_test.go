package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScimGroup_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceScimGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.contractors"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.engineering"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.executives"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.finance"),
					testAccDataSourceScimGroupCheck("data.zpa_scim_groups.marketing"),
				),
			},
		},
	})
}

func testAccDataSourceScimGroupCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
		resource.TestCheckResourceAttrSet(name, "idp_name"),
	)
}

var testAccCheckDataSourceScimGroupConfig_basic = `
data "zpa_scim_groups" "contractors" {
    name = "Contractors"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "engineering" {
    name = "Engineering"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "executives" {
    name = "Executives"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "finance" {
    name = "Finance"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "marketing" {
    name = "Marketing"
	idp_name = "BD_Okta_Users"
}
`
