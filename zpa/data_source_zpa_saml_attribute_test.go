package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSamlAttribute_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceSamlAttributeConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceSamlAttributeCheck("data.zpa_saml_attribute.email_user_sso"),
					testAccDataSourceSamlAttributeCheck("data.zpa_saml_attribute.department"),
					testAccDataSourceSamlAttributeCheck("data.zpa_saml_attribute.first_name"),
					testAccDataSourceSamlAttributeCheck("data.zpa_saml_attribute.last_name"),
					testAccDataSourceSamlAttributeCheck("data.zpa_saml_attribute.group"),
				),
			},
		},
	})
}

func testAccDataSourceSamlAttributeCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceSamlAttributeConfig_basic = `
data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_BD_Okta_Users"
	idp_name = "BD_Okta_Users"
}
data "zpa_saml_attribute" "department" {
    name = "DepartmentName_BD_Okta_Users"
	idp_name = "BD_Okta_Users"
}
data "zpa_saml_attribute" "first_name" {
    name = "FirstName_BD_Okta_Users"
	idp_name = "BD_Okta_Users"
}
data "zpa_saml_attribute" "last_name" {
    name = "LastName_BD_Okta_Users"
	idp_name = "BD_Okta_Users"
}
data "zpa_saml_attribute" "group" {
    name = "GroupName_BD_Okta_Users"
	idp_name = "BD_Okta_Users"
}
`
*/
