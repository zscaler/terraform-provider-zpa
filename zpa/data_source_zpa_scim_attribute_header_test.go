package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScimAttributeHeader_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceScimAttributeHeaderConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_scim_attribute_header.email_value", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_scim_attribute_header.email_value", "idp_name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceScimAttributeHeaderConfig_basic = `
data "zpa_scim_attribute_header" "email_value" {
    name = "emails.value"
    idp_name = "SGIO-User-Okta"
}
`
