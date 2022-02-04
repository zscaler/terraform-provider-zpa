package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAppConnectorGroup(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	resourceName := "zpa_app_connector_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConnectorGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppConnectorGroupConfigBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAppConnectorGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "37.3382082"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "-121.8863286"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
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

func testAccResourceAppConnectorGroupConfigBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_app_connector_group" "test" {
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
	`, rName, rDesc)
}

func tesAccCheckAppConnectorGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("App Connector Group Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no App Connector Group ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.appconnectorgroup.GetByName(rs.Primary.Attributes["name"])
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

func testAccCheckAppConnectorGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_app_connector_group" {
			continue
		}

		_, _, err := client.appconnectorgroup.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("App Connector Group still exists")
		}

		return nil
	}
	return nil
}
