package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLSSClientTypes_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceLSSClientTypesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"data.zpa_lss_config_client_types.foobar", ""),
				),
			},
		},
	})
}

const testAccCheckDataSourceLSSClientTypesConfig_basic = `
data "zpa_lss_config_client_types" "foobar" {
}`
