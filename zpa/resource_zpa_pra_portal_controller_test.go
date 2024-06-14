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
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praportal"
)

func TestAccResourcePRAPortalControllerBasic(t *testing.T) {
	var praPortal praportal.PRAPortal
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAPortalController)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName
	domainName := "pra_" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAPortalControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAPortalControllerConfigure(resourceTypeAndName, initialName, variable.PraPortalDescription, variable.PraPortalEnabled, variable.PraUserNotificationEnabled, domainName, variable.PraUserNotification),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRAPortalControllerExists(resourceTypeAndName, &praPortal),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.PraPortalDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PraPortalEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification_enabled", strconv.FormatBool(variable.PraUserNotificationEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "domain", domainName+".bd-hashicorp.com"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification", variable.PraUserNotification),
				),
			},
			// Update test
			{
				Config: testAccCheckPRAPortalControllerConfigure(resourceTypeAndName, updatedName, variable.PraPortalDescriptionUpdate, variable.PraPortalEnabled, variable.PraUserNotificationEnabled, domainName, variable.PraUserNotification),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRAPortalControllerExists(resourceTypeAndName, &praPortal),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.PraPortalDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PraPortalEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification_enabled", strconv.FormatBool(variable.PraUserNotificationEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "domain", domainName+".bd-hashicorp.com"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification", variable.PraUserNotification),
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

func testAccCheckPRAPortalControllerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPRAPortalController {
			continue
		}

		group, _, err := praportal.Get(apiClient.PRAPortal, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if group != nil {
			return fmt.Errorf("pra portal with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPRAPortalControllerExists(resource string, portal *praportal.PRAPortal) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedPortal, _, err := praportal.Get(apiClient.PRAPortal, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*portal = *receivedPortal

		return nil
	}
}

func testAccCheckPRAPortalControllerConfigure(resourceTypeAndName, generatedName, description string, enabled, notificationEnabled bool, domainName, userNotification string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

data "zpa_ba_certificate" "this" {
	name = "pra01.bd-hashicorp.com"
}

resource "%s" "%s" {
	name = "%s"
	description = "%s"
	domain = "%s.bd-hashicorp.com"
	user_notification = "%s"
	enabled = "%s"
	user_notification_enabled = "%s"
    certificate_id = data.zpa_ba_certificate.this.id
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the certificate
		resourcetype.ZPAPRAPortalController,
		resourceName,
		generatedName,
		description,
		domainName,
		userNotification,
		strconv.FormatBool(enabled),
		strconv.FormatBool(notificationEnabled),

		// Data source type and name
		resourcetype.ZPAPRAPortalController, resourceName,

		// Reference to the resource
		resourcetype.ZPAPRAPortalController, resourceName,
	)
}
