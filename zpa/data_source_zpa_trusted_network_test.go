package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTrustedNetwork_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceTrustedNetworkConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_trusted_network.foobar", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceTrustedNetworkConfig_basic = `
data "zpa_trusted_network" "foobar" {
	name = "Corp-Trusted-Networks (zscalerthree.net)"
}`
