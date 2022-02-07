package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCustomerVersionProfile_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceCustomerVersionProfileConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_customer_version_profile.default", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_customer_version_profile.previous_default", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_customer_version_profile.new_release", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceCustomerVersionProfileConfig_basic = `
data "zpa_customer_version_profile" "default" {
    name = "Default"
}

data "zpa_customer_version_profile" "previous_default" {
    name = "Default"
}

data "zpa_customer_version_profile" "new_release" {
    name = "Default"
}
`
