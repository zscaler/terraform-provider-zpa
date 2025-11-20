package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccSAMLAttributeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSAMLAttributeDataSourceConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccSAMLAttributeDataSourceCheck("data.zpa_saml_attribute.email_user_sso"),
					testAccSAMLAttributeDataSourceCheck("data.zpa_saml_attribute.department"),
					testAccSAMLAttributeDataSourceCheck("data.zpa_saml_attribute.first_name"),
					testAccSAMLAttributeDataSourceCheck("data.zpa_saml_attribute.last_name"),
					testAccSAMLAttributeDataSourceCheck("data.zpa_saml_attribute.group"),
				),
			},
		},
	})
}

func testAccSAMLAttributeDataSourceCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccSAMLAttributeDataSourceConfig_basic = `
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
