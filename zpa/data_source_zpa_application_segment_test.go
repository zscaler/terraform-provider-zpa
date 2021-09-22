package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceApplicationSegment_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	port := acctest.RandIntRange(1000, 9999)
	rDesc := acctest.RandString(15)
	sgrName := acctest.RandString(15)
	sgrDesc := acctest.RandString(15)
	resourceName := "data.zpa_application_segment.by_id"
	resourceName2 := "data.zpa_application_segment.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApplicationSegmentByID(port, rName, rDesc, sgrName, sgrDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceApplicationSegment(resourceName),
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

func testAccDataSourceApplicationSegmentByID(port int, rName, rDesc, sgName, sgDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_application_segment" "test_application" {
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
	data "zpa_application_segment" "by_name" {
		id = zpa_application_segment.test_application.id
	}
	data "zpa_application_segment" "by_id" {
		id = zpa_application_segment.test_application.id
	}
	`, rName, rDesc, port, port, sgName, sgName)
}

func testAccDataSourceApplicationSegment(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
