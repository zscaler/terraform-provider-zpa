package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceApplicationServer_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zpa_application_server.by_id"
	resourceName2 := "data.zpa_application_server.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApplicationServerByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceApplicationServer(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "address", "test.example.com"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName2, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "address", "test.example.com"),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceApplicationServerByID(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zpa_application_server" "test_res" {
		name 		 = "%s"
		description  = "%s"
		address      = "test.example.com"
		enabled      = true
		config_space = "DEFAULT"
	}
	data "zpa_application_server" "by_name" {
		name = zpa_application_server.test_res.name
	}
	data "zpa_application_server" "by_id" {
		id = zpa_application_server.test_res.id
	}
	`, rName, rDesc)
}

func testAccDataSourceApplicationServer(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
