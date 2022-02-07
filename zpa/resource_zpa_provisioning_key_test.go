package zpa

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceProvisioningKey(t *testing.T) {

	rNameConnectorGrp := acctest.RandString(10)
	rNameServiceEdgeGrp := acctest.RandString(20)
	resourceName := "zpa_provisioning_key.test_connector_grp"
	resourceName2 := "zpa_provisioning_key.test_service_edge_grp"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProvisioningKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProvisioningKeyConnectorGroupConfigBasic(rNameConnectorGrp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameConnectorGrp),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_usage", "2"),
					resource.TestCheckResourceAttr(resourceName, "association_type", "CONNECTOR_GRP"),
				),
			},
			{
				Config: testAccResourceProvisioningKeyServiceEdgeConfigBasic(rNameServiceEdgeGrp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "name", rNameServiceEdgeGrp),
					resource.TestCheckResourceAttr(resourceName2, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "max_usage", "2"),
					resource.TestCheckResourceAttr(resourceName2, "association_type", "SERVICE_EDGE_GRP"),
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

func testAccResourceProvisioningKeyConnectorGroupConfigBasic(rNameConnectorGrp string) string {
	return fmt.Sprintf(`
	data "zpa_enrollment_cert" "connector" {
		name = "Connector"
	}
	resource "zpa_app_connector_group" "app_connector_test" {
		name                          = "Test"
		description                   = "Test"
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
	resource "zpa_provisioning_key" "test_connector_grp" {
		name                     = "%s"
		association_type         = "CONNECTOR_GRP"
		enabled                  = true
		enrollment_cert_id       = data.zpa_enrollment_cert.connector.id
		max_usage                = "2"
		zcomponent_id            = zpa_app_connector_group.app_connector_test.id
		depends_on = 			 = [zpa_app_connector_group.app_connector_test]
	}
	`, rNameConnectorGrp)
}

func testAccResourceProvisioningKeyServiceEdgeConfigBasic(rNameServiceEdgeGrp string) string {
	return fmt.Sprintf(`
	data "zpa_enrollment_cert" "service_edge" {
		name = "Service Edge"
	}
	resource "zpa_service_edge_group" "service_edge_group_test" {
		name                 = "Test"
		description          = "Test"
		upgrade_day          = "SUNDAY"
		upgrade_time_in_secs = "66600"
		latitude             = "37.3382082"
		longitude            = "-121.8863286"
		location             = "San Jose, CA, USA"
		override_version_profile = true
		version_profile_id   = "2"
	}
	resource "zpa_provisioning_key" "test_service_edge_grp" {
		name                     = "%s"
		association_type         = "SERVICE_EDGE_GRP"
		enabled                  = true
		enrollment_cert_id       = data.zpa_enrollment_cert.service_edge.id
		max_usage                = "2"
		zcomponent_id            = zpa_service_edge_group.service_edge_group_test.id
		depends_on = 			 = [zpa_app_connector_group.service_edge_group_test]
	}
	`, rNameServiceEdgeGrp)
}

func testAccCheckProvisioningKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Provisioning Key Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Provisioning Key ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.provisioningkey.GetBy(rs.Primary.Attributes["name"])
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

func testAccCheckProvisioningKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_provisioning_key" {
			continue
		}

		_, _, err := client.provisioningkey.GetByName(rs.Primary.Attributes["name"], rs.Primary.Attributes["association_type"])
		if err == nil {
			return fmt.Errorf("Provisioning Key still exists")
		}

		return nil
	}
	return nil
}
*/
