package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccSCIMGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSCIMGroupDataSourceConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccSCIMGroupDataSourceCheck("data.zpa_scim_groups.engineering"),
					testAccSCIMGroupDataSourceCheck("data.zpa_scim_groups.contractors"),
					testAccSCIMGroupDataSourceCheck("data.zpa_scim_groups.finance"),
					testAccSCIMGroupDataSourceCheck("data.zpa_scim_groups.executives"),
				),
			},
		},
	})
}

func testAccSCIMGroupDataSourceCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccSCIMGroupDataSourceConfig_basic = `
data "zpa_scim_groups" "engineering" {
    name = "Engineering"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "contractors" {
    name = "Contractors"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "finance" {
    name = "Finance"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "executives" {
    name = "Executives"
	idp_name = "BD_Okta_Users"
}
`
