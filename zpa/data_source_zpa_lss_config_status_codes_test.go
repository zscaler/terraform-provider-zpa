package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLSSStatusCodes_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: (testAccCheckDataSourceLSSStatusCodesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"data.zpa_lss_config_status_codes.status_codes", ""),
				),
			},
		},
	})
}

var testAccCheckDataSourceLSSStatusCodesConfig_basic = `
data "zpa_lss_config_status_codes" "status_codes" {
}`
