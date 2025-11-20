package datasources_test

import (
	"fmt"
	"strconv"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccProvisioningKeyDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_provisioning_key.test"
	dataSourceName := "data.zpa_provisioning_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisioningKeyDataSourceConfig(rName, "CONNECTOR_GRP", "2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "association_type", resourceName, "association_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_usage", resourceName, "max_usage"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enrollment_cert_id", resourceName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "zcomponent_id", resourceName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", strconv.FormatBool(true)),
				),
			},
		},
	})
}

func testAccProvisioningKeyDataSourceConfig(name, associationType, maxUsage string) string {
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

data "zpa_enrollment_cert" "connector" {
  name = "Connector"
}

resource "zpa_provisioning_key" "test" {
  name                = "tf-acc-test-%s"
  association_type    = "%s"
  enabled             = "%s"
  max_usage           = "%s"
  zcomponent_id       = zpa_app_connector_group.test.id
  enrollment_cert_id  = data.zpa_enrollment_cert.connector.id
  depends_on          = [data.zpa_enrollment_cert.connector, zpa_app_connector_group.test]
}

data "zpa_provisioning_key" "test" {
  id               = zpa_provisioning_key.test.id
  association_type = "CONNECTOR_GRP"
}
`, name, name, associationType, strconv.FormatBool(true), maxUsage)
}
