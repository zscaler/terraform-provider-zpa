package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCustomerVersionProfile_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCustomerVersionProfileConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCustomerVersionProfileCheck("data.zpa_customer_version_profile.default"),
					testAccDataSourceCustomerVersionProfileCheck("data.zpa_customer_version_profile.previous_default"),
					testAccDataSourceCustomerVersionProfileCheck("data.zpa_customer_version_profile.new_release"),
				),
			},
		},
	})
}

func testAccDataSourceCustomerVersionProfileCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceCustomerVersionProfileConfig_basic = `
data "zpa_customer_version_profile" "default" {
    name = "Default"
}

data "zpa_customer_version_profile" "previous_default" {
    name = "Previous Default"
}

data "zpa_customer_version_profile" "new_release" {
    name = "New Release"
}
`
