package datasources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

func TestAccDataSourceApplicationServer_basic(t *testing.T) {
	name := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(6))
	resourceName := "zpa_application_server.test"
	dataSourceName := "data.zpa_application_server.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckApplicationServerDestroyDataSource(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationServerDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "address", resourceName, "address"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccApplicationServerDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zpa_application_server" "test" {
  name        = "%[1]s"
  description = "test application server"
  address     = "apps-server.example.com"
  enabled     = true
}

data "zpa_application_server" "test" {
  id = zpa_application_server.test.id
}
`, name)
}

func testAccCheckApplicationServerDestroyDataSource(zClient *client.Client) resource.TestCheckFunc {
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
