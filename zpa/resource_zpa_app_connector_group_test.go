package zpa

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
)

func TestAccResourceAppConnectorGroup_Basic(t *testing.T) {
	var groups appconnectorgroup.AppConnectorGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConnectorGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppConnectorGroupConfigure(resourceTypeAndName, initialName, variable.AppConnectorDescription, variable.AppConnectorEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AppConnectorDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.AppConnectorEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tcp_quick_ack_app", strconv.FormatBool(variable.TCPQuickAckApp)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tcp_quick_ack_assistant", strconv.FormatBool(variable.TCPQuickAckAssistant)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tcp_quick_ack_read_assistant", strconv.FormatBool(variable.TCPQuickAckReadAssistant)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "use_in_dr_mode", strconv.FormatBool(variable.UseInDrMode)),
				),
			},

			// Update test
			{
				Config: testAccCheckAppConnectorGroupConfigure(resourceTypeAndName, updatedName, variable.AppConnectorDescriptionUpdate, variable.AppConnectorEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AppConnectorDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.AppConnectorEnabledUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tcp_quick_ack_app", strconv.FormatBool(variable.TCPQuickAckAppUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tcp_quick_ack_assistant", strconv.FormatBool(variable.TCPQuickAckAssistantUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tcp_quick_ack_read_assistant", strconv.FormatBool(variable.TCPQuickAckReadAssistantUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "use_in_dr_mode", strconv.FormatBool(variable.UseInDrModeUpdate)),
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

func testAccCheckAppConnectorGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAAppConnectorGroup {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.AppConnectorGroup
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		rule, _, err := appconnectorgroup.Get(service, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("app connector group with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAppConnectorGroupExists(resource string, rule *appconnectorgroup.AppConnectorGroup) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.AppConnectorGroup
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		receivedRule, _, err := appconnectorgroup.Get(service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckAppConnectorGroupConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name                          = "%s"
	description                   = "%s"
	enabled                       = "%s"
	country_code                  = "US"
	city_country                  = "San Jose, US"
	latitude                      = "37.33874"
	longitude                     = "-121.8852525"
	location                      = "San Jose, CA, USA"
	upgrade_day                   = "SUNDAY"
	upgrade_time_in_secs          = "66600"
	override_version_profile      = true
	version_profile_id            = 0
	dns_query_type                = "IPV4_IPV6"
	tcp_quick_ack_app 			  = true
	tcp_quick_ack_assistant 	  = true
	tcp_quick_ack_read_assistant  = true
	use_in_dr_mode 				  = false
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the app connector group
		resourcetype.ZPAAppConnectorGroup,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),

		// Data source type and name
		resourcetype.ZPAAppConnectorGroup,
		resourceName,

		// Reference to the resource
		resourcetype.ZPAAppConnectorGroup, resourceName,
	)
}
