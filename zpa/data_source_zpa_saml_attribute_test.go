package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSamlAttribute_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceSamlAttributeConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_saml_attribute.email_user_sso", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_saml_attribute.department", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_saml_attribute.first_name", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_saml_attribute.last_name", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_saml_attribute.group", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceSamlAttributeConfig_basic = `
data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_SGIO-User-Okta"
}
data "zpa_saml_attribute" "department" {
    name = "DepartmentName_SGIO-User-Okta"
}
data "zpa_saml_attribute" "first_name" {
    name = "FirstName_SGIO-User-Okta"
}
data "zpa_saml_attribute" "last_name" {
    name = "LastName_SGIO-User-Okta"
}
data "zpa_saml_attribute" "group" {
    name = "GroupName_SGIO-User-Okta"
}
`
