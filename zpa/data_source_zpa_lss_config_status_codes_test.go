package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLSSStatusCodes_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceLSSStatusCodesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"data.zpa_lss_config_status_codes.foobar", ""),
				),
			},
		},
	})
}

const testAccCheckDataSourceLSSStatusCodesConfig_basic = `
data "zpa_lss_config_status_codes" "foobar" {
}`
