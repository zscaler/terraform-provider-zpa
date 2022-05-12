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
  association_type = "CONNECTOR_GRP"
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
	zcomponent_id			 = "${%s.id}"
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
func testAccCheckProvisioningKeyServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, name, serviceEdgeGroupHCL, serviceEdgeGroupTypeAndName, provisioningKeyType string) string {
	return fmt.Sprintf(`
// app connector group resource
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
		serviceEdgeGroupProvisioningKeyResourceHCL(generatedName, name, serviceEdgeGroupTypeAndName, provisioningKeyType),

		// data source variables
		resourcetype.ZPAProvisioningKey,
		generatedName,
		resourceTypeAndName,
	)
}

func serviceEdgeGroupProvisioningKeyResourceHCL(generatedName, name, serviceEdgeGroupTypeAndName, provisioningKeyType string) string {
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
		generatedName,
		provisioningKeyType,
		strconv.FormatBool(variable.ProvisioningKeyEnabled),
		variable.ProvisioningKeyUsage,
		serviceEdgeGroupTypeAndName,
		serviceEdgeGroupTypeAndName,
	)
}
*/
