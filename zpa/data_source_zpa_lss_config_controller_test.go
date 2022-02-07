package zpa

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceLSSConfigController_ByIdAndName(t *testing.T) {
	// rName := acctest.RandString(15)
	// port := acctest.RandIntRange(1000, 9999)
	// rDesc := acctest.RandString(15)
	connectorGroupName := acctest.RandString(15)
	connectorGroupDesc := acctest.RandString(15)
	resourceName := "data.zpa_lss_config_controller.by_id"
	// resourceName2 := "data.zpa_lss_config_controller.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLSSConfigControllerByID(connectorGroupName, connectorGroupDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceLSSConfigController(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", connectorGroupName),
					resource.TestCheckResourceAttr(resourceName, "description", connectorGroupDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					// resource.TestCheckResourceAttr(resourceName2, "name", rName),
					// resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
					// resource.TestCheckResourceAttr(resourceName2, "enabled", "true"),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceLSSConfigControllerByID(connectorGroupName, connectorGroupDesc string) string {
	return fmt.Sprintf(`

	data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
		log_type="zpn_trans_log"
	}

	resource "zpa_app_connector_group" "test_app_connector_lss" {
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

	resource "zpa_lss_config_controller" "test_lss_config_controller" {
		config {
		  name            = "LSS_Test"
		  description     = "LSS_Test"
		  enabled         = true
		  format          = data.zpa_lss_config_log_type_formats.zpn_trans_log.json
		  lss_host        = "192.168.1.1"
		  lss_port        = "20000"
		  source_log_type = "zpn_trans_log"
		  use_tls         = true
		}
		connector_groups {
			id = [ zpa_app_connector_group.test_app_connector_lss.id ]
		  }
	}

		data "zpa_lss_config_controller" "by_id" {
			id = zpa_lss_config_controller.test_lss_config_controller.id
		}
	`, connectorGroupName, connectorGroupDesc)
}

func testAccDataSourceLSSConfigController(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
*/
