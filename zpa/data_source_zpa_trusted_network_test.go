package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTrustedNetwork_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceTrustedNetworkConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrustedNetworkCheck("data.zpa_trusted_network.testAcc01"),
					testAccDataSourceTrustedNetworkCheck("data.zpa_trusted_network.testAcc02"),
				),
			},
		},
	})
}

func testAccDataSourceTrustedNetworkCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceTrustedNetworkConfig_basic = `
data "zpa_trusted_network" "testAcc01" {
    name = "BD-TrustedNetwork01"
}

data "zpa_trusted_network" "testAcc02" {
    name = "BD-TrustedNetwork02"
}
`
