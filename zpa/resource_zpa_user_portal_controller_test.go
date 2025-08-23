package zpa

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

func TestAccResourceUserPortalController_Basic(t *testing.T) {
	var userPortalController portal_controller.UserPortalController
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAUserPortalController)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserPortalControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserPortalControllerConfigure(resourceTypeAndName, initialName, variable.UserPortalControllerDescription, variable.UserPortalControllerEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserPortalControllerExists(resourceTypeAndName, &userPortalController),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.UserPortalControllerDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.UserPortalControllerEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification", variable.UserPortalControllerUserNotification),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification_enabled", strconv.FormatBool(variable.UserPortalControllerUserNotificationEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "domain", variable.UserPortalControllerDomain),
				),
			},

			// Update test
			{
				Config: testAccCheckUserPortalControllerConfigureUpdate(resourceTypeAndName, updatedName, variable.UserPortalControllerDescriptionUpdate, variable.UserPortalControllerEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserPortalControllerExists(resourceTypeAndName, &userPortalController),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.UserPortalControllerDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.UserPortalControllerEnabledUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification", variable.UserPortalControllerUserNotificationUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "domain", variable.UserPortalControllerDomainUpdate),
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

func testAccCheckUserPortalControllerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAUserPortalController {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		controller, _, err := portal_controller.Get(context.Background(), service, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if controller != nil {
			return fmt.Errorf("user portal controller with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckUserPortalControllerExists(resource string, controller *portal_controller.UserPortalController) resource.TestCheckFunc {
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
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		receivedController, _, err := portal_controller.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*controller = *receivedController

		return nil
	}
}

func testAccCheckUserPortalControllerConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
data "zpa_ba_certificate" "this" {
  name = "%s"
}

resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = %s
	user_notification = "%s"
	user_notification_enabled = %s
	certificate_id = data.zpa_ba_certificate.this.id
	domain = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Certificate data source
		variable.BACertificateName,

		// Resource type and name for the user portal controller
		resourcetype.ZPAUserPortalController,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		variable.UserPortalControllerUserNotification,
		strconv.FormatBool(variable.UserPortalControllerUserNotificationEnabled),
		variable.UserPortalControllerDomain,

		// Data source type and name
		resourcetype.ZPAUserPortalController, resourceName,

		// Reference to the resource
		resourcetype.ZPAUserPortalController, resourceName,
	)
}

func testAccCheckUserPortalControllerConfigureUpdate(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
data "zpa_ba_certificate" "this" {
  name = "%s"
}

resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = %s
	user_notification = "%s"
	user_notification_enabled = %s
	certificate_id = data.zpa_ba_certificate.this.id
	domain = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Certificate data source
		variable.BACertificateName,

		// Resource type and name for the user portal controller
		resourcetype.ZPAUserPortalController,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		variable.UserPortalControllerUserNotificationUpdate,
		strconv.FormatBool(variable.UserPortalControllerUserNotificationEnabled),
		variable.UserPortalControllerDomainUpdate,

		// Data source type and name
		resourcetype.ZPAUserPortalController, resourceName,

		// Reference to the resource
		resourcetype.ZPAUserPortalController, resourceName,
	)
}
