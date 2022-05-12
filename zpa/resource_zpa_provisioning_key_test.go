package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/provisioningkey"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccResourceProvisioningKeyBasic(t *testing.T) {
	var groups provisioningkey.ProvisioningKey
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAProvisioningKey)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProvisioningKeyDestroy,
		Steps: []resource.TestStep{

			// Test App Connector Group Provisioning Key
			{
				Config: testAccCheckProvisioningKeyAppConnectorGroupConfigure(resourceTypeAndName, generatedName, generatedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ConnectorGroupType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "association_type", variable.ConnectorGroupType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "max_usage", variable.ProvisioningKeyUsage),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},

			// Update App Connector Group Provisioning Key
			{
				Config: testAccCheckProvisioningKeyAppConnectorGroupConfigure(resourceTypeAndName, generatedName, generatedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ConnectorGroupType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "association_type", variable.ConnectorGroupType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "max_usage", variable.ProvisioningKeyUsage),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},
		},
	})
}

func testAccCheckProvisioningKeyDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAProvisioningKey {
			continue
		}

		rule, _, err := apiClient.provisioningkey.GetByName(rs.Primary.Attributes["association_type"], rs.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("provisioning key with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckProvisioningKeyExists(resource string, provisioningkey *provisioningkey.ProvisioningKey) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedKey, _, err := apiClient.provisioningkey.GetByName(rs.Primary.Attributes["association_type"], rs.Primary.Attributes["name"])

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*provisioningkey = *receivedKey

		return nil
	}
}

func testAccCheckProvisioningKeyAppConnectorGroupConfigure(resourceTypeAndName, generatedName, name, appConnectorGroupHCL, appConnectorGroupTypeAndName, provisioningKeyType string) string {
	return fmt.Sprintf(`
// app connector group resource
%s

// provisioning key resource
%s

data "%s" "%s" {
  id = "${%s.id}"
  association_type =  "CONNECTOR_GRP"
}
`,
		// resource variables
		appConnectorGroupHCL,
		appConnectorGroupProvisioningKeyResourceHCL(generatedName, name, appConnectorGroupTypeAndName, provisioningKeyType),

		// data source variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		resourceTypeAndName,
	)
}

func appConnectorGroupProvisioningKeyResourceHCL(generatedName, name, appConnectorGroupTypeAndName, provisioningKeyType string) string {
	return fmt.Sprintf(`

data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}

resource "%s" "%s" {
	name                     = "tf-acc-test-%s"
	association_type         = "%s"
	enabled                  = "%s"
	max_usage                = "%s"
	zcomponent_id			= "%s"
	enrollment_cert_id       = data.zpa_enrollment_cert.connector.id
	depends_on = [ data.zpa_enrollment_cert.connector, %s ]

}
`,
		// resource variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		generatedName,
		provisioningKeyType,
		strconv.FormatBool(variable.ProvisioningKeyEnabled),
		variable.ProvisioningKeyUsage,
		appConnectorGroupTypeAndName,
		appConnectorGroupTypeAndName,
	)
}

/*
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
*/
