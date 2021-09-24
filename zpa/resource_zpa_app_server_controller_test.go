package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceApplicationServer(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	resourceName := "zpa_application_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationServerConfigBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckApplicationServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "address", "test.example.com"),
					resource.TestCheckResourceAttr(resourceName, "config_space", "DEFAULT"),
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

func testAccResourceApplicationServerConfigBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_application_server" "test" {
		name 		 = "%s"
		description  = "%s"
		address      = "test.example.com"
		enabled      = true
		config_space = "DEFAULT"
	}
	`, rName, rDesc)
}

func tesAccCheckApplicationServerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Application Server Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Application Server ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.appservercontroller.GetByName(rs.Primary.Attributes["name"])
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

func testAccCheckApplicationServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_application_server" {
			continue
		}

		_, _, err := client.appservercontroller.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Application Server still exists")
		}

		return nil
	}
	return nil
}
