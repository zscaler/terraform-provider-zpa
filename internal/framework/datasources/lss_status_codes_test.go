package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccDataSourceLSSStatusCodes_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceLSSStatusCodesConfig_basic,
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
