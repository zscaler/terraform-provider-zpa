package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/provisioningkey"
)

func TestAccResourceProvisioningKeyConnector_Basic(t *testing.T) {
	var groups provisioningkey.ProvisioningKey
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAProvisioningKey)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProvisioningKeyDestroyAppConnector,
		Steps: []resource.TestStep{
			// Test App Connector Group Provisioning Key
			{
				Config: testAccCheckProvisioningKeyAppConnectorGroupConfigure(resourceTypeAndName, generatedName, generatedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ConnectorGroupType, variable.ProvisioningKeyUsage),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyAppConnectorExists(resourceTypeAndName, &groups),
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
				Config: testAccCheckProvisioningKeyAppConnectorGroupConfigure(resourceTypeAndName, generatedName, generatedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ConnectorGroupType, variable.ProvisioningKeyUsageUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyAppConnectorExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "association_type", variable.ConnectorGroupType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "max_usage", variable.ProvisioningKeyUsageUpdate),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckProvisioningKeyDestroyAppConnector(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAProvisioningKey {
			continue
		}

		rule, _, err := provisioningkey.GetByName(apiClient.ProvisioningKey, rs.Primary.Attributes["association_type"], rs.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("provisioning key with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckProvisioningKeyAppConnectorExists(resource string, provisioningKey *provisioningkey.ProvisioningKey) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedKey, _, err := provisioningkey.GetByName(apiClient.ProvisioningKey, rs.Primary.Attributes["association_type"], rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*provisioningKey = *receivedKey

		return nil
	}
}

func testAccCheckProvisioningKeyAppConnectorGroupConfigure(resourceTypeAndName, generatedName, name, appConnectorGroupHCL, appConnectorGroupTypeAndName, provisioningKeyType, usage string) string {
	return fmt.Sprintf(`
// app connector group resource
%s

// provisioning key resource
%s

data "%s" "%s" {
  id = "${%s.id}"
  association_type = "CONNECTOR_GRP"
}
`,
		// resource variables
		appConnectorGroupHCL,
		appConnectorGroupProvisioningKeyResourceHCL(generatedName, name, appConnectorGroupTypeAndName, provisioningKeyType, usage),

		// data source variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		resourceTypeAndName,
	)
}

func appConnectorGroupProvisioningKeyResourceHCL(generatedName, name, appConnectorGroupTypeAndName, provisioningKeyType, usage string) string {
	return fmt.Sprintf(`

data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}

resource "%s" "%s" {
	name                     = "tf-acc-test-%s"
	association_type         = "%s"
	enabled                  = "%s"
	max_usage                = "%s"
	zcomponent_id			 = "${%s.id}"
	enrollment_cert_id       = data.zpa_enrollment_cert.connector.id
	depends_on = [ data.zpa_enrollment_cert.connector, %s ]

}
`,
		// resource variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		name,
		provisioningKeyType,
		strconv.FormatBool(variable.ProvisioningKeyEnabled),
		usage,
		appConnectorGroupTypeAndName,
		appConnectorGroupTypeAndName,
	)
}

/*
// Testing Provisioning Key for Service Edge Group
func TestAccResourceProvisioningKeyServiceEdgeGroup_Basic(t *testing.T) {
	var groups provisioningkey.ProvisioningKey
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAProvisioningKey)

	serviceEdgeGroupTypeAndName, _, serviceEdgeGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServiceEdgeGroup)
	serviceEdgeGroupHCL := testAccCheckServiceEdgeGroupConfigure(serviceEdgeGroupTypeAndName, serviceEdgeGroupGeneratedName, variable.ServiceEdgeDescription, variable.ServiceEdgeEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProvisioningKeyDestroyServiceEdgeGroup,
		Steps: []resource.TestStep{
			// Test Service Edge Group Provisioning Key
			{
				Config: testAccCheckProvisioningKeyServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, generatedName, serviceEdgeGroupHCL, serviceEdgeGroupTypeAndName, variable.ServiceEdgeGroupType, variable.ProvisioningKeyUsage),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyServiceEdgeGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "association_type", variable.ServiceEdgeGroupType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "max_usage", variable.ProvisioningKeyUsage),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},

			// Update Service Edge Group Provisioning Key
			{
				Config: testAccCheckProvisioningKeyServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, generatedName, serviceEdgeGroupHCL, serviceEdgeGroupTypeAndName, variable.ServiceEdgeGroupType, variable.ProvisioningKeyUsageUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisioningKeyServiceEdgeGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "association_type", variable.ServiceEdgeGroupType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "max_usage", variable.ProvisioningKeyUsageUpdate),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckProvisioningKeyDestroyServiceEdgeGroup(s *terraform.State) error {
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

func testAccCheckProvisioningKeyServiceEdgeGroupExists(resource string, provisioningkey *provisioningkey.ProvisioningKey) resource.TestCheckFunc {
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

func testAccCheckProvisioningKeyServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, name, serviceEdgeGroupHCL, serviceEdgeGroupTypeAndName, provisioningKeyType, usage string) string {
	return fmt.Sprintf(`
// service edge group resource
%s

// provisioning key resource
%s

data "%s" "%s" {
  id = "${%s.id}"
  association_type =  "SERVICE_EDGE_GRP"
}
`,
		// resource variables
		serviceEdgeGroupHCL,
		serviceEdgeGroupProvisioningKeyResourceHCL(generatedName, name, serviceEdgeGroupTypeAndName, provisioningKeyType, usage),

		// data source variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		resourceTypeAndName,
	)
}

func serviceEdgeGroupProvisioningKeyResourceHCL(generatedName, name, serviceEdgeGroupTypeAndName, provisioningKeyType, usage string) string {
	return fmt.Sprintf(`

	data "zpa_enrollment_cert" "service_edge" {
		name = "Service Edge"
	}

	resource "%s" "%s" {
		name                     = "tf-acc-test-%s"
		association_type         = "%s"
		enabled                  = "%s"
		max_usage                = "%s"
		zcomponent_id			 = "${%s.id}"
		enrollment_cert_id       = data.zpa_enrollment_cert.service_edge.id
		depends_on = [ data.zpa_enrollment_cert.service_edge, %s ]

	}
	`,
		// resource variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		name,
		provisioningKeyType,
		strconv.FormatBool(variable.ProvisioningKeyEnabled),
		usage,
		serviceEdgeGroupTypeAndName,
		serviceEdgeGroupTypeAndName,
	)
}
*/
