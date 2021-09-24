package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceZPAServerGroup_basic(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	resourceName := "zpa_server_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZPAServerGroupConfigBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
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

func testAccResourceZPAServerGroupConfigBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_server_group" "test" {
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
	`, rName, rDesc)
}

func tesAccCheckServerGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Server Group Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Server Group ID is set")
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
			return fmt.Errorf("Server group still exists")
		}

		return nil
	}
	return nil
}
