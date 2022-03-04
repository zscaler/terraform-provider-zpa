package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/provisioningkey"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccResourceProvisioningKeyBasic(t *testing.T) {
	var provisioning_key provisioningkey.ProvisioningKey
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAProvisioningKey)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConnectorProvisioningKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppConnectorProvisioningKeyConfigure(resourceTypeAndName, generatedName, variable.ConnectorGroupType, variable.ProvisioningKeyUsage, variable.ProvisioningKeyEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorProvisioningKeyExists(resourceTypeAndName, &provisioning_key),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "association_type", variable.ConnectorGroupType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "max_usage", variable.ProvisioningKeyUsage),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},

			// Update test
			{
				Config: testAccCheckAppConnectorProvisioningKeyConfigure(resourceTypeAndName, generatedName, variable.ConnectorGroupType, variable.ProvisioningKeyUsage, variable.ProvisioningKeyEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorProvisioningKeyExists(resourceTypeAndName, &provisioning_key),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "association_type", variable.ConnectorGroupType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "max_usage", variable.ProvisioningKeyUsage),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},
		},
	})
}

func testAccCheckAppConnectorProvisioningKeyDestroy(s *terraform.State) error {
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

func testAccCheckAppConnectorProvisioningKeyExists(resource string, key *provisioningkey.ProvisioningKey) resource.TestCheckFunc {
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
		*key = *receivedKey

		return nil
	}
}

func testAccCheckAppConnectorProvisioningKeyConfigure(resourceTypeAndName, generatedName, association_type, max_usage string, enabled bool) string {
	return fmt.Sprintf(`
data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}

resource "zpa_app_connector_group" "testAcc_connector_group" {
	name                          = "TestAcc_Connector_Group"
	description                   = "TestAcc_Connector_Group"
	enabled                       = true
	country_code                  = "USA"
	latitude                      = "49.1041779"
	longitude                     = "-122.6603519"
	location                      = "New York, NY, USA"
	upgrade_day                   = "SUNDAY"
	upgrade_time_in_secs          = "66600"
	override_version_profile      = true
	version_profile_id            = "2"
	dns_query_type                = "IPV4"
  }

resource "%s" "%s" {
	name             = "%s"
	association_type = "%s"
	max_usage        = "%s"
	enabled 		 = "%s"
	enrollment_cert_id = data.zpa_enrollment_cert.connector.id
	zcomponent_id = zpa_app_connector_group.testAcc_connector_group.id
}

`,
		// resource variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		generatedName,
		variable.ConnectorGroupType,
		variable.ProvisioningKeyUsage,
		strconv.FormatBool(enabled),
	)
}
