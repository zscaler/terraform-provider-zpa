package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccAppConnectorGroup_basic(t *testing.T) {
	var group appconnectorgroup.AppConnectorGroup
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_app_connector_group.test"
	initialName := fmt.Sprintf("tf-acc-test-%s", rName)
	updatedName := fmt.Sprintf("tf-updated-%s", rName)
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAppConnectorGroupDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConnectorGroupConfig(initialName, "testAcc_app_connector_group", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAppConnectorGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", initialName),
					resource.TestCheckResourceAttr(resourceName, "description", "testAcc_app_connector_group"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tcp_quick_ack_app", "true"),
					resource.TestCheckResourceAttr(resourceName, "tcp_quick_ack_assistant", "true"),
					resource.TestCheckResourceAttr(resourceName, "tcp_quick_ack_read_assistant", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_in_dr_mode", "false"),
				),
			},
			{
				Config: testAccAppConnectorGroupConfig(updatedName, "this is update app connector group test", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAppConnectorGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "this is update app connector group test"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tcp_quick_ack_app", "true"),
					resource.TestCheckResourceAttr(resourceName, "tcp_quick_ack_assistant", "true"),
					resource.TestCheckResourceAttr(resourceName, "tcp_quick_ack_read_assistant", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_in_dr_mode", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAppConnectorGroupExists(zClient *client.Client, resourceName string, group *appconnectorgroup.AppConnectorGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("app connector group not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("app connector group ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		received, _, err := appconnectorgroup.Get(ctx, service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch app connector group %s: %w", rs.Primary.ID, err)
		}

		*group = *received
		return nil
	}
}

func testAccCheckAppConnectorGroupDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_app_connector_group" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			_, _, err := appconnectorgroup.Get(ctx, service, rs.Primary.ID)
			if err == nil {
				if _, delErr := appconnectorgroup.Delete(ctx, service, rs.Primary.ID); delErr != nil {
					if respErr, ok := delErr.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
						continue
					}
					return fmt.Errorf("app connector group %s still exists and failed to delete: %w", rs.Primary.ID, delErr)
				}
				continue
			}

			if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
				continue
			}

			return fmt.Errorf("error checking app connector group %s destruction: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccAppConnectorGroupConfig(name, description string, enabled bool) string {
	return fmt.Sprintf(`
resource "zpa_app_connector_group" "test" {
  name                          = "%s"
  description                   = "%s"
  enabled                       = "%t"
  country_code                  = "US"
  city_country                  = "San Jose, US"
  latitude                      = "37.33874"
  longitude                     = "-121.8852525"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  version_profile_name          = "Default"
  dns_query_type                = "IPV4_IPV6"
  tcp_quick_ack_app             = true
  tcp_quick_ack_assistant       = true
  tcp_quick_ack_read_assistant  = true
  use_in_dr_mode                = false
}
`, name, description, enabled)
}
