package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccSCIMAttributeHeaderDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSCIMAttributeHeaderDataSourceConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccSCIMAttributeHeaderDataSourceCheck("data.zpa_scim_attribute_header.email_value"),
					testAccSCIMAttributeHeaderDataSourceCheck("data.zpa_scim_attribute_header.cost_center"),
					testAccSCIMAttributeHeaderDataSourceCheck("data.zpa_scim_attribute_header.department"),
					testAccSCIMAttributeHeaderDataSourceCheck("data.zpa_scim_attribute_header.name_family_name"),
				),
			},
		},
	})
}

func testAccSCIMAttributeHeaderDataSourceCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
		resource.TestCheckResourceAttrSet(name, "idp_name"),
	)
}

var testAccSCIMAttributeHeaderDataSourceConfig_basic = `
data "zpa_scim_attribute_header" "email_value" {
    name = "emails.value"
    idp_name = "BD_Okta_Users"
}

data "zpa_scim_attribute_header" "cost_center" {
    name = "costCenter"
    idp_name = "BD_Okta_Users"
}

data "zpa_scim_attribute_header" "department" {
    name = "department"
    idp_name = "BD_Okta_Users"
}

data "zpa_scim_attribute_header" "name_family_name" {
    name = "name.familyName"
    idp_name = "BD_Okta_Users"
}
`
