package resources_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccProvisioningKey_basic(t *testing.T) {
	var key provisioningkey.ProvisioningKey
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_provisioning_key.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckProvisioningKeyDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisioningKeyConfig(rName, "CONNECTOR_GRP", "2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisioningKeyExists(zClient, resourceName, &key),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "association_type", "CONNECTOR_GRP"),
					resource.TestCheckResourceAttr(resourceName, "max_usage", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccProvisioningKeyConfig(rName, "CONNECTOR_GRP", "10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisioningKeyExists(zClient, resourceName, &key),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "association_type", "CONNECTOR_GRP"),
					resource.TestCheckResourceAttr(resourceName, "max_usage", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "zcomponent_id"),
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

func testAccCheckProvisioningKeyExists(zClient *client.Client, resourceName string, key *provisioningkey.ProvisioningKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("provisioning key not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("provisioning key ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		associationType := rs.Primary.Attributes["association_type"]
		received, _, err := provisioningkey.GetByName(ctx, service, associationType, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("failed to fetch provisioning key %s: %w", rs.Primary.ID, err)
		}

		*key = *received
		return nil
	}
}

func testAccCheckProvisioningKeyDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_provisioning_key" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			associationType := rs.Primary.Attributes["association_type"]

			// Use Get by ID instead of GetByName to avoid cached search query results
			// Get by ID uses a different cache key that should be invalidated on DELETE
			key, _, err := provisioningkey.Get(ctx, service, associationType, rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("id %s already exists", rs.Primary.ID)
			}

			if key != nil {
				return fmt.Errorf("provisioning key with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccProvisioningKeyConfig(name, associationType, maxUsage string) string {
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

data "zpa_enrollment_cert" "connector" {
  name = "Connector"
}

resource "zpa_provisioning_key" "test" {
  name                = "tf-acc-test-%s"
  association_type    = "%s"
  enabled             = "%s"
  max_usage           = "%s"
  zcomponent_id       = zpa_app_connector_group.test.id
  enrollment_cert_id  = data.zpa_enrollment_cert.connector.id
  depends_on          = [data.zpa_enrollment_cert.connector, zpa_app_connector_group.test]
}
`, name, name, associationType, strconv.FormatBool(true), maxUsage)
}
