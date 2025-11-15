package resources_test

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

const (
	privateCloudGroupDescription             = "Example private cloud group"
	privateCloudGroupDescriptionUpdate       = "Updated example private cloud group"
	privateCloudGroupCityCountry             = "San Jose, US"
	privateCloudGroupLatitude                = "37.33874"
	privateCloudGroupLongitude               = "-121.8852525"
	privateCloudGroupLocation                = "San Jose, CA, USA"
	privateCloudGroupUpgradeDay              = "SUNDAY"
	privateCloudGroupUpgradeTimeInSecs       = "66600"
	privateCloudGroupUpgradeDayUpdate        = "MONDAY"
	privateCloudGroupUpgradeTimeInSecsUpdate = "72000"
	privateCloudGroupVersionProfileID        = "0"
	privateCloudGroupIsPublic                = "TRUE"
)

func TestAccPrivateCloudGroup_basic(t *testing.T) {
	var group private_cloud_group.PrivateCloudGroup

	nameSuffix := sdkacctest.RandString(6)
	initialName := fmt.Sprintf("tf-acc-test-%s", nameSuffix)
	updatedName := fmt.Sprintf("tf-updated-%s", nameSuffix)
	resourceName := "zpa_private_cloud_group.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPrivateCloudGroupDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateCloudGroupConfig(initialName, privateCloudGroupDescription, privateCloudGroupUpgradeDay, privateCloudGroupUpgradeTimeInSecs),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPrivateCloudGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", initialName),
					resource.TestCheckResourceAttr(resourceName, "description", privateCloudGroupDescription),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "city_country", privateCloudGroupCityCountry),
					resource.TestCheckResourceAttr(resourceName, "latitude", privateCloudGroupLatitude),
					resource.TestCheckResourceAttr(resourceName, "longitude", privateCloudGroupLongitude),
					resource.TestCheckResourceAttr(resourceName, "location", privateCloudGroupLocation),
					resource.TestCheckResourceAttr(resourceName, "upgrade_day", privateCloudGroupUpgradeDay),
					resource.TestCheckResourceAttr(resourceName, "upgrade_time_in_secs", privateCloudGroupUpgradeTimeInSecs),
					resource.TestCheckResourceAttr(resourceName, "version_profile_id", privateCloudGroupVersionProfileID),
					resource.TestCheckResourceAttr(resourceName, "override_version_profile", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_public", privateCloudGroupIsPublic),
				),
			},
			{
				Config: testAccPrivateCloudGroupConfig(updatedName, privateCloudGroupDescriptionUpdate, privateCloudGroupUpgradeDayUpdate, privateCloudGroupUpgradeTimeInSecsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPrivateCloudGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", privateCloudGroupDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "upgrade_day", privateCloudGroupUpgradeDayUpdate),
					resource.TestCheckResourceAttr(resourceName, "upgrade_time_in_secs", privateCloudGroupUpgradeTimeInSecsUpdate),
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

func testAccPrivateCloudGroupConfig(name, description, upgradeDay, upgradeTime string) string {
	return fmt.Sprintf(`
resource "zpa_private_cloud_group" "test" {
  name                     = "%[1]s"
  description              = "%[2]s"
  enabled                  = true
  city_country             = "%[3]s"
  latitude                 = "%[4]s"
  longitude                = "%[5]s"
  location                 = "%[6]s"
  upgrade_day              = "%[7]s"
  upgrade_time_in_secs     = "%[8]s"
  version_profile_id       = "%[9]s"
  override_version_profile = true
  is_public                = "%[10]s"
}
`, name, description, privateCloudGroupCityCountry, privateCloudGroupLatitude, privateCloudGroupLongitude, privateCloudGroupLocation, upgradeDay, upgradeTime, privateCloudGroupVersionProfileID, privateCloudGroupIsPublic)
}

func testAccCheckPrivateCloudGroupExists(zClient *client.Client, resourceName string, group *private_cloud_group.PrivateCloudGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("private cloud group not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("private cloud group ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		resp, _, err := private_cloud_group.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch private cloud group %s: %w", rs.Primary.ID, err)
		}

		*group = *resp
		return nil
	}
}

func testAccCheckPrivateCloudGroupDestroy(zClient *client.Client) func(*terraform.State) error {
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
