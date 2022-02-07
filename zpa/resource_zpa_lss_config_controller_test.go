package zpa

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceLSSConfigController(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	port := acctest.RandIntRange(1000, 9999)
	appConnectorName := acctest.RandString(10)
	appConnectorDesc := acctest.RandString(10)
	resourceName := "zpa_lss_config_controller.testAcc_lss_controller"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLSSConfigControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLSSConfigControllerConfigBasic(port, rName, rDesc, appConnectorName, appConnectorDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLSSConfigControllerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "config.name", rName),
					resource.TestCheckResourceAttr(resourceName, "config.description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "config.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "config.lss_host", "192.168.1.1"),
					// resource.TestCheckResourceAttr(resourceName, "config.lss_port", "11001"),
					resource.TestCheckResourceAttr(resourceName, "config.source_log_type", "zpn_trans_log"),
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

func testAccResourceLSSConfigControllerConfigBasic(port int, rName, rDesc, appConnectorName, appConnectorDesc string) string {
	return fmt.Sprintf(`

data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
	log_type = "zpn_trans_log"
  }

resource "zpa_lss_config_controller" "testAcc_lss_controller" {
	config {
		name            = "%s"
		description     = "%s"
		enabled         = true
		format          = data.zpa_lss_config_log_type_formats.zpn_trans_log.json
		lss_host        = "192.168.1.1"
		lss_port        = "%d"
		source_log_type = "zpn_trans_log"
		use_tls         = true
	}
	connector_groups {
		id = [zpa_app_connector_group.testAcc_lss_app_connector_group.id]
	  }
	depends_on = [zpa_app_connector_group.testAcc_lss_app_connector_group]
}

resource "zpa_app_connector_group" "testAcc_lss_app_connector_group" {
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
	`, rName, rDesc, port, appConnectorName, appConnectorDesc)
}

func testAccCheckLSSConfigControllerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("LSS Config Controller Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no LSS Config Controller ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.lssconfigcontroller.GetByName(rs.Primary.Attributes["config.name"])
		if err != nil {
			return err
		}
		if resp.LSSConfig.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("name Not found in created attributes")
		}
		if resp.LSSConfig.Description != rs.Primary.Attributes["description"] {
			return fmt.Errorf("description Not found in created attributes")
		}
		return nil
	}
}

func testAccCheckLSSConfigControllerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_lss_config_controller" {
			continue
		}

		_, _, err := client.lssconfigcontroller.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("LSS Config Controller still exists")
		}

		return nil
	}
	return nil
}
*/
