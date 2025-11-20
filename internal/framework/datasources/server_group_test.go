package datasources_test

import (
	"fmt"
	"strconv"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccServerGroupDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_server_group.test"
	dataSourceName := "data.zpa_server_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServerGroupDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "dynamic_discovery", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "app_connector_groups.#", "1"),
				),
			},
		},
	})
}

func testAccServerGroupDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "zpa_app_connector_group" "test" {
  name                          = "tf-acc-test-%s"
  description                   = "testAcc_app_connector_group"
  enabled                       = "true"
  country_code                  = "US"
  city_country                  = "San Jose, US"
  latitude                      = "37.33874"
  longitude                     = "-121.8852525"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  dns_query_type                = "IPV4_IPV6"
  tcp_quick_ack_app             = true
  tcp_quick_ack_assistant       = true
  tcp_quick_ack_read_assistant  = true
  use_in_dr_mode                = false
}

resource "zpa_server_group" "test" {
  name             = "tf-acc-test-%s"
  description      = "tf-acc-test-%s"
  enabled          = "%s"
  dynamic_discovery = "%s"
  app_connector_groups {
    id = [zpa_app_connector_group.test.id]
  }
  depends_on = [zpa_app_connector_group.test]
}

data "zpa_server_group" "test" {
  id = zpa_server_group.test.id
}
`, rName, rName, rName, strconv.FormatBool(true), strconv.FormatBool(true))
}
