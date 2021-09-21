package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceSegmentGroup_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zpa_segment_group.by_id"
	resourceName2 := "data.zpa_segment_group.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSegmentGroupById(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceSegmentGroup(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
				PreventPostDestroyRefresh: true,
			},
			{

				Config: testAccDataSourceSegmentGroupByName(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceSegmentGroup(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName2, "enabled", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSegmentGroupById(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "test_segment_group" {
	name = "%s"
	description = "%s"
	enabled = true
	config_space = "DEFAULT"
}
data "zpa_segment_group" "by_id" {
	id = zpa_segment_group.test_segment_group.id
}
	`, rName, rDesc)
}

func testAccDataSourceSegmentGroupByName(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "zpa_segment_group" "test_segment_group" {
		name = "%s"
		description = "%s"
		enabled = true
		config_space = "DEFAULT"
	}
data "zpa_segment_group" "by_name" {
	name = zpa_segment_group.test_segment_group.name
}
	`, rName, rDesc)
}

func testAccDataSourceSegmentGroup(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
