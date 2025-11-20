package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccDataSourceLSSClientTypes_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceLSSClientTypesConfig_basic,
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
