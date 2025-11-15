package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"

	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/helpers"
)

func TestAccC2CIPRanges_basic(t *testing.T) {
	var ipRange c2c_ip_ranges.IPRanges

	name := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(6))
	resourceName := "zpa_c2c_ip_ranges.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckC2CIPRangesDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccC2CIPRangesConfig(name, "testAcc_c2c_ip_ranges", true, "192.168.1.1", "192.168.1.254"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckC2CIPRangesExists(zClient, resourceName, &ipRange),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "testAcc_c2c_ip_ranges"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "location_hint", "Created_via_Terraform"),
					resource.TestCheckResourceAttr(resourceName, "ip_range_begin", "192.168.1.1"),
					resource.TestCheckResourceAttr(resourceName, "ip_range_end", "192.168.1.254"),
					resource.TestCheckResourceAttr(resourceName, "location", "San Jose, CA, USA"),
					resource.TestCheckResourceAttr(resourceName, "sccm_flag", "true"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "latitude_in_db", "37.33874"),
					resource.TestCheckResourceAttr(resourceName, "longitude_in_db", "-121.8852525"),
				),
			},
			{
				Config: testAccC2CIPRangesConfig(name, "this is update c2c ip ranges test", true, "10.0.1.1", "10.0.1.254"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckC2CIPRangesExists(zClient, resourceName, &ipRange),
					resource.TestCheckResourceAttr(resourceName, "description", "this is update c2c ip ranges test"),
					resource.TestCheckResourceAttr(resourceName, "ip_range_begin", "10.0.1.1"),
					resource.TestCheckResourceAttr(resourceName, "ip_range_end", "10.0.1.254"),
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

func testAccC2CIPRangesConfig(name, description string, enabled bool, start, end string) string {
	return fmt.Sprintf(`
resource "zpa_c2c_ip_ranges" "test" {
  name           = "%[1]s"
  description    = "%[2]s"
  enabled        = %[3]t
  location_hint  = "Created_via_Terraform"
  ip_range_begin = "%[4]s"
  ip_range_end   = "%[5]s"
  location       = "San Jose, CA, USA"
  sccm_flag      = true
  country_code   = "US"
  latitude_in_db = "37.33874"
  longitude_in_db = "-121.8852525"
}
`, name, description, enabled, start, end)
}

func testAccCheckC2CIPRangesExists(zClient *client.Client, resourceName string, ipRange *c2c_ip_ranges.IPRanges) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s ID is not set", resourceName)
		}

		service := zClient.Service
		if micro := rs.Primary.Attributes["microtenant_id"]; micro != "" {
			service = service.WithMicroTenant(micro)
		}

		resp, _, err := c2c_ip_ranges.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch c2c ip range %s: %w", rs.Primary.ID, err)
		}
		*ipRange = *resp
		return nil
	}
}

func testAccCheckC2CIPRangesDestroy(zClient *client.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_c2c_ip_ranges" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if micro := rs.Primary.Attributes["microtenant_id"]; micro != "" {
				service = service.WithMicroTenant(micro)
			}

			_, _, err := c2c_ip_ranges.Get(ctx, service, rs.Primary.ID)
			if err == nil {
				_, delErr := c2c_ip_ranges.Delete(ctx, service, rs.Primary.ID)
				if delErr == nil || helpers.IsObjectNotFoundError(delErr) {
					continue
				}
				return fmt.Errorf("c2c ip range %s still exists: %v", rs.Primary.ID, delErr)
			}

			if helpers.IsObjectNotFoundError(err) {
				continue
			}

			return fmt.Errorf("error checking c2c ip range %s destruction: %w", rs.Primary.ID, err)
		}

		return nil
	}
}
