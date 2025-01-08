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
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praconsole"
)

func TestAccResourcePRAConsoleController_Basic(t *testing.T) {
	var praConsole praconsole.PRAConsole
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAConsoleController)
	domainName := "pra_" + generatedName

	praPortalTypeAndName, _, praPortalGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAPortalController)
	praPortalHCL := testAccCheckPRAPortalControllerConfigure(praPortalTypeAndName, "tf-acc-test-"+praPortalGeneratedName, variable.PraPortalDescription, variable.PraPortalEnabled, variable.PraUserNotificationEnabled, domainName, variable.PraUserNotification)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAConsoleControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAConsoleControllerConfigure(resourceTypeAndName, generatedName, initialName, variable.PraConsoleDescription, variable.PraConsoleEnabled, praPortalHCL, praPortalTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRAConsoleControllerExists(resourceTypeAndName, &praConsole),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.PraConsoleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PraConsoleEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pra_application.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pra_portals.#", "1"),
				),
			},
			// Update test
			{
				Config: testAccCheckPRAConsoleControllerConfigure(resourceTypeAndName, generatedName, updatedName, variable.PraConsoleDescriptionUpdate, variable.PraConsoleEnabled, praPortalHCL, praPortalTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRAConsoleControllerExists(resourceTypeAndName, &praConsole),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.PraConsoleDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PraConsoleEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pra_application.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pra_portals.#", "1"),
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

func testAccCheckPRAConsoleControllerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPRAConsoleController {
			continue
		}

		console, _, err := praconsole.Get(apiClient.PRAConsole, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if console != nil {
			return fmt.Errorf("pra console with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPRAConsoleControllerExists(resource string, console *praconsole.PRAConsole) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedConsole, _, err := praconsole.Get(apiClient.PRAConsole, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*console = *receivedConsole

		return nil
	}
}

func testAccCheckPRAConsoleControllerConfigure(resourceTypeAndName, generatedName, name, description string, enabled bool, praPortalHCL, praPortalTypeAndName string) string {
	return fmt.Sprintf(`

// pra portal controller resource
%s

// pra console controller resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		praPortalHCL,
		getPRAConsoleControllerResourceHCL(generatedName, name, description, praPortalTypeAndName, enabled),

		// data source variables
		resourcetype.ZPAPRAConsoleController,
		generatedName,
		resourceTypeAndName,
	)
}

func getPRAConsoleControllerResourceHCL(generatedName, name, description, praPortalTypeAndName string, enabled bool) string {
	return fmt.Sprintf(`

resource "zpa_application_segment_pra" "this" {
	name             = "tf-acc-test-%s"
	description      = "tf-acc-test-%s"
	enabled          = true
	health_reporting = "ON_ACCESS"
	bypass_type      = "NEVER"
	is_cname_enabled = true
	tcp_port_ranges  = ["3223", "3223", "3392", "3392"]
	domain_names     = ["ssh_pra3223.example.com", "rdp_pra3392.example.com"]
	segment_group_id = zpa_segment_group.this.id
	common_apps_dto {
		apps_config {
		domain               = "rdp_pra3392.example.com"
		application_protocol = "RDP"
		connection_security  = "ANY"
		application_port     = "3392"
		enabled              = true
		app_types            = ["SECURE_REMOTE_ACCESS"]
		}
		apps_config {
		domain               = "ssh_pra3223.example.com"
		application_protocol = "SSH"
		application_port     = "3223"
		enabled              = true
		app_types            = ["SECURE_REMOTE_ACCESS"]
		}
	}
}
  
resource "zpa_segment_group" "this" {
	name        = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled     = true
}

data "zpa_application_segment_by_type" "rdp_pra3392" {
    application_type = "SECURE_REMOTE_ACCESS"
    name = "rdp_pra3392"
	depends_on = [zpa_application_segment_pra.this]
}

resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = "%s"
	pra_application {
		id = data.zpa_application_segment_by_type.rdp_pra3392.id
	}
	 pra_portals {
		id = ["${%s.id}"]
	}
	depends_on = [%s]
}
`,
		// PRA Application Segment and Segment Group name generation
		generatedName,
		generatedName,
		generatedName,
		generatedName,

		// Resource type and name for the PRA Console
		resourcetype.ZPAPRAConsoleController,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		praPortalTypeAndName,
		praPortalTypeAndName,
	)
}
