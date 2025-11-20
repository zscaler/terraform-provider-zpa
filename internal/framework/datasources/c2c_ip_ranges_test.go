package datasources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/helpers"
)

func TestAccDataSourceC2CIPRanges_basic(t *testing.T) {
	name := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(6))
	resourceName := "zpa_c2c_ip_ranges.test"
	dataName := "data.zpa_c2c_ip_ranges.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckC2CIPRangesDestroyDataSource(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccC2CIPRangesDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "enabled", resourceName, "enabled"),
					resource.TestCheckResourceAttrPair(dataName, "location_hint", resourceName, "location_hint"),
					resource.TestCheckResourceAttrPair(dataName, "ip_range_begin", resourceName, "ip_range_begin"),
					resource.TestCheckResourceAttrPair(dataName, "ip_range_end", resourceName, "ip_range_end"),
					resource.TestCheckResourceAttrPair(dataName, "location", resourceName, "location"),
					resource.TestCheckResourceAttrPair(dataName, "sccm_flag", resourceName, "sccm_flag"),
					resource.TestCheckResourceAttrPair(dataName, "country_code", resourceName, "country_code"),
					resource.TestCheckResourceAttrPair(dataName, "latitude_in_db", resourceName, "latitude_in_db"),
					resource.TestCheckResourceAttrPair(dataName, "longitude_in_db", resourceName, "longitude_in_db"),
				),
			},
		},
	})
}

func testAccC2CIPRangesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zpa_c2c_ip_ranges" "test" {
  name           = "%[1]s"
  description    = "testAcc_c2c_ip_ranges"
  enabled        = true
  location_hint  = "Created_via_Terraform"
  ip_range_begin = "192.168.1.1"
  ip_range_end   = "192.168.1.254"
  location       = "San Jose, CA, USA"
  sccm_flag      = true
  country_code   = "US"
  latitude_in_db = "37.33874"
  longitude_in_db = "-121.8852525"
}

data "zpa_c2c_ip_ranges" "test" {
  id = zpa_c2c_ip_ranges.test.id
}
`, name)
}

func testAccCheckC2CIPRangesDestroyDataSource(zClient *client.Client) resource.TestCheckFunc {
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
