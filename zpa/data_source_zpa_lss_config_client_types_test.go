package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLSSClientTypes_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: (testAccCheckDataSourceLSSClientTypesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"data.zpa_lss_config_client_types.all_client_types", ""),
				),
			},
		},
	})
}

var testAccCheckDataSourceLSSClientTypesConfig_basic = `
data "zpa_lss_config_client_types" "all_client_types" {
}`
