package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceProvisioningKey(t *testing.T) {

	rName1 := acctest.RandString(10)
	rName2 := acctest.RandString(10)
	resourceName1 := "zpa_provisioning_key.test_connector_grp"
	resourceName2 := "zpa_provisioning_key.test_edge_connector_grp"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProvisioningKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProvisioningKeyConnectorGroupConfigBasic(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyExists(resourceName1),
					resource.TestCheckResourceAttr(resourceName1, "name", rName1),
					resource.TestCheckResourceAttr(resourceName1, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName1, "max_usage", "2"),
					resource.TestCheckResourceAttr(resourceName1, "association_type", "CONNECTOR_GRP"),
				),
			},
			{
				ResourceName:      resourceName1,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceProvisioningKeyServiceEdgeConfigBasic(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "name", rName2),
					resource.TestCheckResourceAttr(resourceName2, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "max_usage", "2"),
					resource.TestCheckResourceAttr(resourceName2, "association_type", "SERVICE_EDGE_GRP"),
				),
			},
			{
				ResourceName:      resourceName2,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceProvisioningKeyConnectorGroupConfigBasic(rName1 string) string {
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
        version_profile_id            = "2"
        dns_query_type                = "IPV4"
    }
    resource "zpa_provisioning_key" "test_connector_grp" {
        name                     = "%s"
        association_type         = "CONNECTOR_GRP"
        enabled                  = true
        enrollment_cert_id       = data.zpa_enrollment_cert.connector.id
        max_usage                = "2"
        zcomponent_id            = zpa_app_connector_group.app_connector_test.id
    }
    `, rName1)
}

func testAccResourceProvisioningKeyServiceEdgeConfigBasic(rName2 string) string {
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
    resource "zpa_provisioning_key" "test_edge_connector_grp" {
        name                     = "%s"
        association_type         = "SERVICE_EDGE_GRP"
        enabled                  = true
        enrollment_cert_id       = data.zpa_enrollment_cert.service_edge.id
        max_usage                = "2"
        zcomponent_id            = zpa_service_edge_group.service_edge_group_test.id
    }
    `, rName2)
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
		resp, _, err := client.provisioningkey.GetByName(rs.Primary.Attributes["association_type"], rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}
		if resp.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("name Not found in created attributes")
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

		_, _, err := client.provisioningkey.GetByName(rs.Primary.Attributes["association_type"], rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Provisioning Key still exists")
		}

		return nil
	}
	return nil
}
