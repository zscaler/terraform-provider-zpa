package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceProvisioningKey_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.zpa_provisioning_key.test-provisioning-key"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProvisioningKeyBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceProvisioningKey(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-provisioning-key-"+rName),
					resource.TestCheckResourceAttr(resourceName, "association_type", "CONNECTOR_GRP"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_usage", "2"),
				),
			},
		},
	})
}

func testAccDataSourceProvisioningKeyBasic(rName string) string {
	return fmt.Sprintf(`

resource "zpa_app_connector_group" "app_connector_test" {
	name                          = "Test"
	description                   = "Test"
	enabled                       = true
	country_code                  = "US"
	latitude                      = "37.3382082"
	longitude                     = "-121.8863286"
	location                      = "San Jose, CA, USA"
	upgrade_day                   = "SUNDAY"
	upgrade_time_in_secs          = "66600"
	override_version_profile      = true
	version_profile_id            = "2"
	dns_query_type                = "IPV4"
}

data "zpa_enrollment_cert" "connector" {
	name = "Connector"
}

resource "zpa_provisioning_key" "test-provisioning-key" {
	name                     = "test-provisioning-key-%s"
	association_type         = "CONNECTOR_GRP"
	enabled                  = true
	enrollment_cert_id       = data.zpa_enrollment_cert.connector.id
	max_usage                = "2"
	zcomponent_id            = zpa_app_connector_group.app_connector_test.id
	depends_on 				 = [ zpa_app_connector_group.app_connector_test ]
}

data "zpa_provisioning_key" "test-provisioning-key" {
	name 				= zpa_provisioning_key.test-provisioning-key.name
    association_type 	= "CONNECTOR_GRP"
}

	`, rName)
}

func testAccDataSourceProvisioningKey(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
