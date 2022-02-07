package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceServiceEdgeGroup_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zpa_service_edge_group.by_id"
	resourceName2 := "data.zpa_service_edge_group.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServiceEdgeGroupByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceServiceEdgeGroup(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "37.3382082"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "-121.8863286"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "37.3382082"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "-121.8863286"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceServiceEdgeGroupByID(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "zpa_service_edge_group" "testAcc" {
	name                 = "%s"
	description          = "%s"
	upgrade_day          = "SUNDAY"
	upgrade_time_in_secs = "66600"
	latitude             = "37.3382082"
	longitude            = "-121.8863286"
	location             = "San Jose, CA, USA"
	version_profile_id   = "0"
}

data "zpa_service_edge_group" "by_name" {
	name = zpa_service_edge_group.testAcc.name
}
data "zpa_service_edge_group" "by_id" {
	id = zpa_service_edge_group.testAcc.id
}
	`, rName, rDesc)
}

func testAccDataSourceServiceEdgeGroup(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
