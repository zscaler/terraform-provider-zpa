package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAppConnectorGroup_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zpa_app_connector_group.by_id"
	resourceName2 := "data.zpa_app_connector_group.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppConnectorGroupByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAppConnectorGroup(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "37.3382082"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "-121.8863286"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "37.3382082"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "-121.8863286"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceAppConnectorGroupByID(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_app_connector_group" "test_res" {
		name                          = "%s"
		description                   = "%s"
		enabled                       = true
		country_code                  = "US"
		latitude                      = "37.3382082"
		longitude                     = "-121.8863286"
		location                      = "San Jose, CA, USA"
		upgrade_day                   = "SUNDAY"
		upgrade_time_in_secs          = "66600"
		override_version_profile      = true
		version_profile_id            = 0
		dns_query_type                = "IPV4"
	}
	data "zpa_app_connector_group" "by_name" {
		name = zpa_app_connector_group.test_res.name
	}
	data "zpa_app_connector_group" "by_id" {
		id = zpa_app_connector_group.test_res.id
	}
	`, rName, rDesc)
}

func testAccDataSourceAppConnectorGroup(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
