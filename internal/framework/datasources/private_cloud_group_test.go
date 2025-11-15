package datasources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

func TestAccDataSourcePrivateCloudGroup_basic(t *testing.T) {
	name := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(6))
	resourceName := "zpa_private_cloud_group.test"
	dataSourceName := "data.zpa_private_cloud_group.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPrivateCloudGroupDestroyDataSource(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateCloudGroupDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled", resourceName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "city_country", resourceName, "city_country"),
					resource.TestCheckResourceAttrPair(dataSourceName, "latitude", resourceName, "latitude"),
					resource.TestCheckResourceAttrPair(dataSourceName, "longitude", resourceName, "longitude"),
					resource.TestCheckResourceAttrPair(dataSourceName, "location", resourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceName, "upgrade_day", resourceName, "upgrade_day"),
					resource.TestCheckResourceAttrPair(dataSourceName, "upgrade_time_in_secs", resourceName, "upgrade_time_in_secs"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version_profile_id", resourceName, "version_profile_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "override_version_profile", resourceName, "override_version_profile"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_public", resourceName, "is_public"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccPrivateCloudGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zpa_private_cloud_group" "test" {
  name                     = "%[1]s"
  description              = "Example private cloud group"
  enabled                  = true
  city_country             = "San Jose, US"
  latitude                 = "37.33874"
  longitude                = "-121.8852525"
  location                 = "San Jose, CA, USA"
  upgrade_day              = "SUNDAY"
  upgrade_time_in_secs     = "66600"
  version_profile_id       = "0"
  override_version_profile = true
  is_public                = "TRUE"
}

data "zpa_private_cloud_group" "test" {
  id = zpa_private_cloud_group.test.id
}
`, name)
}

func testAccCheckPrivateCloudGroupDestroyDataSource(zClient *client.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_private_cloud_group" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			_, _, err := private_cloud_group.Get(ctx, service, rs.Primary.ID)
			if err == nil {
				if _, delErr := private_cloud_group.Delete(ctx, service, rs.Primary.ID); delErr != nil && !helpers.IsObjectNotFoundError(delErr) {
					return fmt.Errorf("private cloud group %s still exists and couldn't be deleted: %w", rs.Primary.ID, delErr)
				}
				continue
			}

			if helpers.IsObjectNotFoundError(err) {
				continue
			}

			return fmt.Errorf("error checking private cloud group %s destruction: %w", rs.Primary.ID, err)
		}

		return nil
	}
}
