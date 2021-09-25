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
		applications {
			id = []
		}
		app_connector_groups {
			id = [data.zpa_app_connector_group.example.id]
		}
	}
	data "zpa_app_connector_group" "example" {
		name = "SGIO-Vancouver"
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
