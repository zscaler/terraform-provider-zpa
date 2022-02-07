package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceServerGroup(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	connGroupName := acctest.RandString(10)
	connGroupDesc := acctest.RandString(20)
	resourceName := "zpa_server_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServerGroupConfigBasic(rName, rDesc, connGroupName, connGroupDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "dynamic_discovery", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceServerGroupConfigBasic(rName, rDesc, connGroupName, connGroupDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_server_group" "test" {
		name              = "%s"
		description       = "%s"
		enabled           = true
		dynamic_discovery = true
		app_connector_groups {
		  id = [zpa_app_connector_group.test_app_connector_group.id]
		}
		depends_on = [zpa_app_connector_group.test_app_connector_group]
	}
	resource "zpa_app_connector_group" "test_app_connector_group" {
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
	`, rName, rDesc, connGroupName, connGroupDesc)
}

func tesAccCheckServerGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Server Group Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Server Group Group ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.servergroup.GetByName(rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}
		if resp.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("name Not found in created attributes")
		}
		if resp.Description != rs.Primary.Attributes["description"] {
			return fmt.Errorf("description Not found in created attributes")
		}
		return nil
	}
}

func testAccCheckServerGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_server_group" {
			continue
		}

		_, _, err := client.servergroup.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Server Group still exists")
		}

		return nil
	}
	return nil
}
