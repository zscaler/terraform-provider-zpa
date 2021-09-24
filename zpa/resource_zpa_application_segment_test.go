package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceZPAApplicationSegment(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	sgName := acctest.RandString(10)
	sgDesc := acctest.RandString(20)
	port := acctest.RandIntRange(1000, 9999)
	resourceName := "zpa_application_segment.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationSegmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZPAApplicationSegmentConfigBasic(port, rName, rDesc, sgName, sgDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckApplicationSegmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttr(resourceName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(resourceName, "is_cname_enabled", "true"),
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

func testAccResourceZPAApplicationSegmentConfigBasic(port int, rName, rDesc, sgName, sgDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_application_segment" "test" {
		name = "%s"
		description = "%s"
		enabled = true
		health_reporting = "ON_ACCESS"
		bypass_type = "NEVER"
		is_cname_enabled = true
		tcp_port_ranges = ["%d", "%d"]
		domain_names = ["test.example.com"]
		segment_group_id = zpa_segment_group.test_app_group.id
		server_groups {
			id = []
		}
	}
	resource "zpa_segment_group" "test_app_group" {
		name = "%s"
		description = "%s"
		enabled = true
	}
	`, rName, rDesc, port, port, sgName, sgDesc)
}

func tesAccCheckApplicationSegmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Application Segment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Application Segment ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.applicationsegment.GetByName(rs.Primary.Attributes["name"])
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

func testAccCheckApplicationSegmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_application_segment" {
			continue
		}

		_, _, err := client.applicationsegment.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Application Segment still exists")
		}

		return nil
	}
	return nil
}
