package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

func TestAccApplicationServer_basic(t *testing.T) {
	var server appservercontroller.ApplicationServer

	suffix := sdkacctest.RandString(6)
	initialName := fmt.Sprintf("tf-acc-test-%s", suffix)
	updatedName := fmt.Sprintf("tf-updated-%s", suffix)
	resourceName := "zpa_application_server.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckApplicationServerDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationServerConfig(initialName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckApplicationServerExists(zClient, resourceName, &server),
					resource.TestCheckResourceAttr(resourceName, "name", initialName),
					resource.TestCheckResourceAttr(resourceName, "description", "test application server"),
					resource.TestCheckResourceAttr(resourceName, "address", "apps-server.example.com"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccApplicationServerConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckApplicationServerExists(zClient, resourceName, &server),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "test application server"),
					resource.TestCheckResourceAttr(resourceName, "address", "apps-server.example.com"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
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

func testAccApplicationServerConfig(name string) string {
	return fmt.Sprintf(`
resource "zpa_application_server" "test" {
  name        = "%s"
  description = "test application server"
  address     = "apps-server.example.com"
  enabled     = true
}
`, name)
}

func testAccCheckApplicationServerExists(zClient *client.Client, resourceName string, server *appservercontroller.ApplicationServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("application server not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("application server ID not set")
		}

		service := zClient.Service
		if micro := rs.Primary.Attributes["microtenant_id"]; micro != "" {
			service = service.WithMicroTenant(micro)
		}

		resp, _, err := appservercontroller.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve application server %s: %w", rs.Primary.ID, err)
		}
		*server = *resp
		return nil
	}
}

func testAccCheckApplicationServerDestroy(zClient *client.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_application_server" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if micro := rs.Primary.Attributes["microtenant_id"]; micro != "" {
				service = service.WithMicroTenant(micro)
			}

			_, _, err := appservercontroller.Get(ctx, service, rs.Primary.ID)
			if err == nil {
				if _, delErr := appservercontroller.Delete(ctx, service, rs.Primary.ID); delErr != nil && !helpers.IsObjectNotFoundError(delErr) {
					return fmt.Errorf("application server %s still exists and could not be deleted: %w", rs.Primary.ID, delErr)
				}
				continue
			}

			if helpers.IsObjectNotFoundError(err) {
				continue
			}

			return fmt.Errorf("error checking application server %s destruction: %w", rs.Primary.ID, err)
		}

		return nil
	}
}
