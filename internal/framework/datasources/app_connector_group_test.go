package datasources_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccAppConnectorGroupDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_app_connector_group.test"
	dataSourceName := "data.zpa_app_connector_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConnectorGroupDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled", resourceName, "enabled"),
				),
			},
		},
	})
}

func testAccAppConnectorGroupDataSourceConfig(rName string) string {
	name := fmt.Sprintf("tf-acc-test-%s", rName)
	return fmt.Sprintf(`
resource "zpa_app_connector_group" "test" {
  name                          = "%s"
  description                   = "testAcc_app_connector_group"
  enabled                       = "true"
  country_code                  = "US"
  city_country                  = "San Jose, US"
  latitude                      = "37.33874"
  longitude                     = "-121.8852525"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  version_profile_name          = "Default"
  dns_query_type                = "IPV4_IPV6"
  tcp_quick_ack_app             = true
  tcp_quick_ack_assistant       = true
  tcp_quick_ack_read_assistant  = true
  use_in_dr_mode                = false
}

data "zpa_app_connector_group" "test" {
  id = zpa_app_connector_group.test.id
}
`, name)
}
