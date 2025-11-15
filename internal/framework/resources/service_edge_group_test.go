package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccServiceEdgeGroup_basic(t *testing.T) {
	var group serviceedgegroup.ServiceEdgeGroup
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_service_edge_group.test"
	initialName := fmt.Sprintf("tf-acc-test-%s", rName)
	updatedName := fmt.Sprintf("tf-updated-%s", rName)
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckServiceEdgeGroupDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceEdgeGroupConfig(initialName, "testAcc_service_edge_group", true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceEdgeGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", initialName),
					resource.TestCheckResourceAttr(resourceName, "description", "testAcc_service_edge_group"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "37.33874"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "-121.8852525"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
					resource.TestCheckResourceAttr(resourceName, "version_profile_name", "Default"),
				),
			},
			{
				Config: testAccServiceEdgeGroupConfig(updatedName, "testAcc_service_edge_group", true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceEdgeGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "testAcc_service_edge_group"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "37.33874"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "-121.8852525"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
					resource.TestCheckResourceAttr(resourceName, "version_profile_name", "Default"),
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

func testAccCheckServiceEdgeGroupExists(zClient *client.Client, resourceName string, group *serviceedgegroup.ServiceEdgeGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("service edge group not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("service edge group ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		received, _, err := serviceedgegroup.Get(ctx, service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch service edge group %s: %w", rs.Primary.ID, err)
		}

		*group = *received
		return nil
	}
}

func testAccCheckServiceEdgeGroupDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_service_edge_group" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			_, _, err := serviceedgegroup.Get(ctx, service, rs.Primary.ID)
			if err == nil {
				if _, delErr := serviceedgegroup.Delete(ctx, service, rs.Primary.ID); delErr != nil {
					if respErr, ok := delErr.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
						continue
					}
					return fmt.Errorf("service edge group %s still exists and failed to delete: %w", rs.Primary.ID, delErr)
				}
				continue
			}

			if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
				continue
			}

			return fmt.Errorf("error checking service edge group %s destruction: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccServiceEdgeGroupConfig(name, description string, enabled, isPublic bool) string {
	return fmt.Sprintf(`
resource "zpa_service_edge_group" "test" {
  name                      = "%s"
  description               = "%s"
  enabled                   = "%t"
  is_public                 = "%t"
  upgrade_day               = "SUNDAY"
  upgrade_time_in_secs      = "66600"
  country_code              = "US"
  city_country              = "San Jose, US"
  latitude                  = "37.33874"
  longitude                 = "-121.8852525"
  location                  = "San Jose, CA, USA"
  override_version_profile  = true
  version_profile_id        = 0
  version_profile_name      = "Default"
  grace_distance_enabled    = true
  grace_distance_value      = "10"
  grace_distance_value_unit = "KMS"
}
`, name, description, enabled, isPublic)
}
