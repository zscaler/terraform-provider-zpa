package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceBrowserAccess(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	port := acctest.RandIntRange(1000, 65635)
	sgName := acctest.RandString(10)
	sgDesc := acctest.RandString(20)
	srvName := acctest.RandString(10)
	srvDesc := acctest.RandString(20)
	appConnectorName := acctest.RandString(10)
	appConnectorDesc := acctest.RandString(10)
	resourceName := "zpa_browser_access.testAcc_browser_access"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZPABrowserAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZPABrowserAccessConfigBasic(port, rName, rDesc, sgName, sgDesc, srvName, srvDesc, appConnectorName, appConnectorDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZPABrowserAccessExists(resourceName),
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

func testAccResourceZPABrowserAccessConfigBasic(port int, rName, rDesc, sgName, sgDesc, srvName, srvDesc, appConnectorName, appConnectorDesc string) string {
	return fmt.Sprintf(`
data "zpa_ba_certificate" "jenkins" {
	name = "jenkins.securitygeek.io"
}
resource "zpa_browser_access" "testAcc_browser_access" {
	name             = "%s"
	description      = "%s"
	enabled          = true
	health_reporting = "ON_ACCESS"
	bypass_type      = "NEVER"
	is_cname_enabled = true
	tcp_port_range {
        from = "%d"
        to = "%d"
    }
	domain_names     = ["jenkins.securitygeek.io"]
	segment_group_id = zpa_segment_group.testAcc_segment_group.id

	clientless_apps {
		name                 = "jenkins.securitygeek.io"
		application_protocol = "HTTP"
		application_port     = "%d"
		certificate_id       = data.zpa_ba_certificate.jenkins.id
		trust_untrusted_cert = true
		enabled              = true
		domain               = "jenkins.securitygeek.io"
	}
	server_groups {
		id = [
			zpa_server_group.testAcc_server_group.id
		]
	}
	depends_on = [zpa_server_group.testAcc_server_group, zpa_segment_group.testAcc_segment_group]
}
resource "zpa_segment_group" "testAcc_segment_group" {
	name = "%s"
	description = "%s"
	enabled = true
}

resource "zpa_server_group" "testAcc_server_group" {
	name              = "%s"
	description       = "%s"
	enabled           = true
	dynamic_discovery = true
	app_connector_groups {
		id = [zpa_app_connector_group.testAcc_app_connector_group.id]
	}
	depends_on = [zpa_app_connector_group.testAcc_app_connector_group]
}

resource "zpa_app_connector_group" "testAcc_app_connector_group" {
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

	`, rName, rDesc, port, port, port, sgName, sgDesc, srvName, srvDesc, appConnectorName, appConnectorDesc)
}

func testAccCheckZPABrowserAccessExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Browser Access App Segment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Browser Access App Segment ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.browseraccess.GetByName(rs.Primary.Attributes["name"])
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

func testAccCheckZPABrowserAccessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_browser_access" {
			continue
		}

		_, _, err := client.applicationsegment.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Browser Access App Segment still exists")
		}

		return nil
	}
	return nil
}
