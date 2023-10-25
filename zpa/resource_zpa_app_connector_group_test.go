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
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
)

func TestAccResourceAppConnectorGroupBasic(t *testing.T) {
	var groups appconnectorgroup.AppConnectorGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConnectorGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppConnectorGroupConfigure(resourceTypeAndName, generatedName, variable.AppConnectorDescription, variable.AppConnectorEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
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
				Config: testAccCheckAppConnectorGroupConfigure(resourceTypeAndName, generatedName, variable.AppConnectorDescriptionUpdate, variable.AppConnectorEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
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

		rule, _, err := apiClient.appconnectorgroup.Get(rs.Primary.ID)

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
		receivedRule, _, err := apiClient.appconnectorgroup.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckAppConnectorGroupConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	return fmt.Sprintf(`
// app connector group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		appConnectorGroupResourceHCL(generatedName, description, enabled),

		// data source variables
		resourcetype.ZPAAppConnectorGroup,
		generatedName,
		resourceTypeAndName,
	)
}

func appConnectorGroupResourceHCL(generatedName, description string, enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name                          = "tf-acc-test-%s"
	description                   = "%s"
	enabled                       = "%s"
	country_code                  = "US"
	city_country                  = "San Jose, US"
	latitude                      = "37.3382082"
	longitude                     = "-121.8863286"
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
`,
		// resource variables
		resourcetype.ZPAAppConnectorGroup,
		generatedName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
	)
}
