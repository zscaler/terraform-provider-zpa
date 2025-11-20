package datasources_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccServiceEdgeGroupDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_service_edge_group.test"
	dataSourceName := "data.zpa_service_edge_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceEdgeGroupDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
					resource.TestCheckResourceAttrPair(dataSourceName, "latitude", resourceName, "latitude"),
					resource.TestCheckResourceAttrPair(dataSourceName, "longitude", resourceName, "longitude"),
					resource.TestCheckResourceAttrPair(dataSourceName, "location", resourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version_profile_name", resourceName, "version_profile_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "grace_distance_enabled", resourceName, "grace_distance_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "grace_distance_value", resourceName, "grace_distance_value"),
					resource.TestCheckResourceAttrPair(dataSourceName, "grace_distance_value_unit", resourceName, "grace_distance_value_unit"),
				),
			},
		},
	})
}

func testAccServiceEdgeGroupDataSourceConfig(rName string) string {
	name := fmt.Sprintf("tf-acc-test-%s", rName)
	return fmt.Sprintf(`
resource "zpa_service_edge_group" "test" {
  name                      = "%s"
  description               = "testAcc_service_edge_group"
  enabled                   = "true"
  is_public                 = "true"
  upgrade_day               = "SUNDAY"
  upgrade_time_in_secs      = "66600"
  country_code              = "US"
  city_country              = "San Jose, US"
  latitude                  = "37.33874"
  longitude                 = "-121.8852525"
  location                  = "San Jose, CA, USA"
  override_version_profile  = true
  version_profile_id        = 0
  version_profile_name      = "Default"
  grace_distance_enabled    = true
  grace_distance_value      = "10"
  grace_distance_value_unit = "KMS"
}

data "zpa_service_edge_group" "test" {
  id = zpa_service_edge_group.test.id
}
`, name)
}
