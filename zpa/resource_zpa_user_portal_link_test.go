package zpa

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

func TestAccResourceUserPortalLink_Basic(t *testing.T) {
	var userPortalLink portal_link.UserPortalLink
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAUserPortalLink)

	userPortalControllerTypeAndName, _, userPortalControllerGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAUserPortalController)
	userPortalControllerHCL := testAccCheckUserPortalControllerConfigure(userPortalControllerTypeAndName, userPortalControllerGeneratedName, variable.UserPortalLinkDescription, variable.UserPortalLinkEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserPortalLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserPortalLinkConfigure(resourceTypeAndName, generatedName, generatedName, generatedName, userPortalControllerHCL, userPortalControllerTypeAndName, variable.UserPortalLinkEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserPortalLinkExists(resourceTypeAndName, &userPortalLink),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.UserPortalLinkEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "link", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "icon_text", variable.UserPortalLinkIconText),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocol", variable.UserPortalLinkProtocol),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_portals.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckUserPortalLinkConfigureUpdate(resourceTypeAndName, generatedName, generatedName, generatedName, userPortalControllerHCL, userPortalControllerTypeAndName, variable.UserPortalLinkEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserPortalLinkExists(resourceTypeAndName, &userPortalLink),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.UserPortalLinkEnabledUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "link", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "icon_text", variable.UserPortalLinkIconTextUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocol", variable.UserPortalLinkProtocolUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_portals.#", "1"),
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

func testAccCheckUserPortalLinkDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAUserPortalLink {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		link, _, err := portal_link.Get(context.Background(), service, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if link != nil {
			return fmt.Errorf("user portal link with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckUserPortalLinkExists(resource string, link *portal_link.UserPortalLink) resource.TestCheckFunc {
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

		receivedLink, _, err := portal_link.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*link = *receivedLink

		return nil
	}
}

func testAccCheckUserPortalLinkConfigure(resourceTypeAndName, generatedName, name, description, userPortalControllerHCL, userPortalControllerTypeAndName string, enabled bool) string {
	return fmt.Sprintf(`

// user portal controller resource
%s

// user portal link resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		userPortalControllerHCL,
		getUserPortalLinkResourceHCL(generatedName, name, description, userPortalControllerTypeAndName, enabled),

		// data source variables
		resourcetype.ZPAUserPortalLink,
		generatedName,
		resourceTypeAndName,
	)
}

func getUserPortalLinkResourceHCL(generatedName, name, description, userPortalControllerTypeAndName string, enabled bool) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled = "%s"
	link = "tf-acc-test-%s"
	icon_text = "%s"
	protocol = "%s"
	user_portals {
		id = ["${%s.id}"]
	}
	depends_on = [ %s ]
}
`,
		// resource variables
		resourcetype.ZPAUserPortalLink,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		name,
		variable.UserPortalLinkIconText,
		variable.UserPortalLinkProtocol,
		userPortalControllerTypeAndName,
		userPortalControllerTypeAndName,
	)
}

func testAccCheckUserPortalLinkConfigureUpdate(resourceTypeAndName, generatedName, name, description, userPortalControllerHCL, userPortalControllerTypeAndName string, enabled bool) string {
	return fmt.Sprintf(`

// user portal controller resource
%s

// user portal link resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		userPortalControllerHCL,
		getUserPortalLinkResourceHCLUpdate(generatedName, name, description, userPortalControllerTypeAndName, enabled),

		// data source variables
		resourcetype.ZPAUserPortalLink,
		generatedName,
		resourceTypeAndName,
	)
}

func getUserPortalLinkResourceHCLUpdate(generatedName, name, description, userPortalControllerTypeAndName string, enabled bool) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled = "%s"
	link = "tf-acc-test-%s"
	icon_text = "%s"
	protocol = "%s"
	user_portals {
		id = ["${%s.id}"]
	}
	depends_on = [ %s ]
}
`,
		// resource variables
		resourcetype.ZPAUserPortalLink,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		name,
		variable.UserPortalLinkIconTextUpdate,
		variable.UserPortalLinkProtocolUpdate,
		userPortalControllerTypeAndName,
		userPortalControllerTypeAndName,
	)
}
