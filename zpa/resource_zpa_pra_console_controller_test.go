package zpa

/*
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

func TestAccResourcePRAConsoleControllerBasic(t *testing.T) {
	var praConsole praconsole.PRAConsole
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAConsoleController)
	domainName := "pra_" + generatedName

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	appSegmentTypeAndName, _, appSegmentGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationSegmentPRA)
	applicationSegmentHCL := testAccCheckApplicationSegmentPRAConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, appSegmentGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, variable.AppSegmentEnabled, variable.AppSegmentCnameEnabled)

	praPortalTypeAndName, _, praPortalGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAPortalController)
	praPortalHCL := testAccCheckPRAPortalControllerConfigure(praPortalTypeAndName, praPortalGeneratedName, variable.PraPortalDescription, variable.PraPortalEnabled, variable.PraUserNotificationEnabled, domainName, variable.PraUserNotification)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAConsoleControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAConsoleControllerConfigure(resourceTypeAndName, generatedName, initialName, variable.PraConsoleDescription, variable.PraConsoleEnabled, applicationSegmentHCL, appSegmentTypeAndName, praPortalHCL, praPortalTypeAndName),
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
				Config: testAccCheckPRAConsoleControllerConfigure(resourceTypeAndName, generatedName, updatedName, variable.PraConsoleDescription, variable.PraConsoleEnabled, applicationSegmentHCL, appSegmentTypeAndName, praPortalHCL, praPortalTypeAndName),
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

		console, _, err := apiClient.praconsole.Get(rs.Primary.ID)

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
		receivedConsole, _, err := apiClient.praconsole.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*console = *receivedConsole

		return nil
	}
}

func testAccCheckPRAConsoleControllerConfigure(resourceTypeAndName, generatedName, name, description string, enabled bool, applicationSegmentHCL, appSegmentTypeAndName, praPortalHCL, praPortalTypeAndName string) string {
	return fmt.Sprintf(`

// pra application segment resource
%s

// pra portal controller resource
%s

// pra console controller resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		applicationSegmentHCL,
		praPortalHCL,
		getPRAConsoleControllerResourceHCL(generatedName, name, description, appSegmentTypeAndName, praPortalTypeAndName, enabled),

		// data source variables
		resourcetype.ZPAPRAConsoleController,
		generatedName,
		resourceTypeAndName,
	)
}

func getPRAConsoleControllerResourceHCL(generatedName, name, description, appSegmentTypeAndName, praPortalTypeAndName string, enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = "%s"
	pra_application {
		id = "${%s.id}"
	}
	 pra_portals {
		id = ["${%s.id}"]
	}
	depends_on = [%s, %s]
}
`,
		// Resource type and name for the certificate
		resourcetype.ZPAPRAConsoleController,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		appSegmentTypeAndName,
		praPortalTypeAndName,
		praPortalTypeAndName,
		appSegmentTypeAndName,
	)
}
*/
