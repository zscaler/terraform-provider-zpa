package resources_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccPRAPortalController_basic(t *testing.T) {
	var portal praportal.PRAPortal
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_pra_portal_controller.test"
	zClient := acctest.TestClient(t)

	initialName := fmt.Sprintf("tf-acc-test-%s", rName)
	updatedName := fmt.Sprintf("tf-updated-%s", rName)
	domainName := fmt.Sprintf("pra_%s", rName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPRAPortalControllerDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPRAPortalControllerConfig(initialName, "Portal Controller Test", true, true, domainName, "Created with Terraform"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPRAPortalControllerExists(zClient, resourceName, &portal),
					resource.TestCheckResourceAttr(resourceName, "name", initialName),
					resource.TestCheckResourceAttr(resourceName, "description", "Portal Controller Test"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "user_notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("%s.bd-hashicorp.com", domainName)),
					resource.TestCheckResourceAttr(resourceName, "user_notification", "Created with Terraform"),
				),
			},
			// Update test
			{
				Config: testAccPRAPortalControllerConfig(updatedName, "Portal Controller Test Update", true, true, domainName, "Created with Terraform"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPRAPortalControllerExists(zClient, resourceName, &portal),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "Portal Controller Test Update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "user_notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("%s.bd-hashicorp.com", domainName)),
					resource.TestCheckResourceAttr(resourceName, "user_notification", "Created with Terraform"),
				),
			},
			// Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPRAPortalControllerExists(zClient *client.Client, resourceName string, portal *praportal.PRAPortal) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("pra portal controller not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("pra portal controller ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		received, _, err := praportal.Get(ctx, service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch pra portal controller %s: %w", rs.Primary.ID, err)
		}

		*portal = *received
		return nil
	}
}

func testAccCheckPRAPortalControllerDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_pra_portal_controller" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			portal, _, err := praportal.Get(ctx, service, rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("id %s already exists", rs.Primary.ID)
			}

			if portal != nil {
				return fmt.Errorf("pra portal controller with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccPRAPortalControllerConfig(name, description string, enabled, notificationEnabled bool, domainName, userNotification string) string {
	return fmt.Sprintf(`
data "zpa_ba_certificate" "this" {
	name = "pra01.bd-hashicorp.com"
}

resource "zpa_pra_portal_controller" "test" {
	name = "%s"
	description = "%s"
	domain = "%s.bd-hashicorp.com"
	user_notification = "%s"
	enabled = "%s"
	user_notification_enabled = "%s"
	certificate_id = data.zpa_ba_certificate.this.id
}

data "zpa_pra_portal_controller" "test" {
  id = zpa_pra_portal_controller.test.id
}
`, name, description, domainName, userNotification, strconv.FormatBool(enabled), strconv.FormatBool(notificationEnabled))
}
