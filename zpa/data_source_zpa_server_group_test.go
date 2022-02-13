package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceServerGroup_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zpa_server_group.by_id"
	resourceName2 := "data.zpa_server_group.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerGroupByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceServerGroup(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName2, "enabled", "true"),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceServerGroupByID(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_server_group" "test_app_group" {
		name = "%s"
		description = "%s"
		enabled = true
		dynamic_discovery = true
		app_connector_groups {
			id = [zpa_app_connector_group.testAcc.id]
		}
		depends_on = [zpa_app_connector_group.testAcc]
	}

	resource "zpa_app_connector_group" "testAcc" {
		name                          = "testAcc"
		description                   = "testAcc"
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
	  
	data "zpa_server_group" "by_name" {
		name = zpa_server_group.test_app_group.name
	}
	data "zpa_server_group" "by_id" {
		id = zpa_server_group.test_app_group.id
	}
	`, rName, rDesc)
}

func testAccDataSourceServerGroup(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
