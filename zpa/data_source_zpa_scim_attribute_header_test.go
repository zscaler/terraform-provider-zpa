package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScimAttributeHeader_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceScimAttributeHeaderConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceScimAttributeHeaderCheck("data.zpa_scim_attribute_header.email_value"),
					testAccDataSourceScimAttributeHeaderCheck("data.zpa_scim_attribute_header.cost_center"),
					testAccDataSourceScimAttributeHeaderCheck("data.zpa_scim_attribute_header.department"),
					testAccDataSourceScimAttributeHeaderCheck("data.zpa_scim_attribute_header.name_family_name"),
				),
			},
		},
	})
}

func testAccDataSourceScimAttributeHeaderCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
		resource.TestCheckResourceAttrSet(name, "idp_name"),
	)
}

var testAccCheckDataSourceScimAttributeHeaderConfig_basic = `
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
`*/
