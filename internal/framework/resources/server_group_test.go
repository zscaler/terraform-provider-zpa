package resources_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccServerGroup_basic(t *testing.T) {
	var group servergroup.ServerGroup
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_server_group.test"
	dataSourceName := "data.zpa_server_group.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckServerGroupDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccServerGroupConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "dynamic_discovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "app_connector_groups.#", "1"),
					// Data source checks
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttr(dataSourceName, "app_connector_groups.#", "1"),
				),
			},
			// Update test
			{
				Config: testAccServerGroupConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerGroupExists(zClient, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "dynamic_discovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "app_connector_groups.#", "1"),
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

func testAccCheckServerGroupExists(zClient *client.Client, resourceName string, group *servergroup.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("server group not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server group ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		received, _, err := servergroup.Get(ctx, service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch server group %s: %w", rs.Primary.ID, err)
		}

		*group = *received
		return nil
	}
}

func testAccCheckServerGroupDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_server_group" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			// Use Get by ID to avoid cached search query results
			// Get by ID uses a different cache key that should be invalidated on DELETE
			group, _, err := servergroup.Get(ctx, service, rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("id %s already exists", rs.Primary.ID)
			}

			if group != nil {
				return fmt.Errorf("server group with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccServerGroupConfig(rName string) string {
	return fmt.Sprintf(`
resource "zpa_app_connector_group" "test" {
  name                          = "tf-acc-test-%s"
  description                   = "testAcc_app_connector_group"
  enabled                       = "true"
  country_code                  = "US"
  city_country                  = "San Jose, US"
  latitude                      = "37.33874"
  longitude                     = "-121.8852525"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  dns_query_type                = "IPV4_IPV6"
  tcp_quick_ack_app             = true
  tcp_quick_ack_assistant       = true
  tcp_quick_ack_read_assistant  = true
  use_in_dr_mode                = false
}

resource "zpa_server_group" "test" {
  name             = "tf-acc-test-%s"
  description      = "tf-acc-test-%s"
  enabled          = "%s"
  dynamic_discovery = "%s"
  app_connector_groups {
    id = [zpa_app_connector_group.test.id]
  }
  depends_on = [zpa_app_connector_group.test]
}

data "zpa_server_group" "test" {
  id = zpa_server_group.test.id
}
`, rName, rName, rName, strconv.FormatBool(true), strconv.FormatBool(true))
}
