package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceZPASegmentGroup_basic(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	resourceName := "zpa_segment_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZPASegmentGroupConfigBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckSegmentGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
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

func testAccResourceZPASegmentGroupConfigBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_segment_group" "test" {
		name= "%s"
		description= "%s"
	}
	`, rName, rDesc)
}

func tesAccCheckSegmentGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Segment Group Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Segment Group ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.segmentgroup.GetByName(rs.Primary.Attributes["name"])
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

func testAccCheckSegmentGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_segment_group" {
			continue
		}

		_, _, err := client.segmentgroup.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Segment group still exists")
		}

		return nil
	}
	return nil
}
