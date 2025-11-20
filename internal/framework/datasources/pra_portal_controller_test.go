package datasources_test

import (
	"fmt"
	"strconv"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccPRAPortalControllerDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_pra_portal_controller.test"
	dataSourceName := "data.zpa_pra_portal_controller.test"
	domainName := fmt.Sprintf("pra_%s", rName)
	generatedName := fmt.Sprintf("tf-acc-test-%s", rName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPRAPortalControllerDataSourceConfig(generatedName, "Portal Controller Test", true, true, domainName, "Created with Terraform"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr(resourceName, "user_notification_enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttrPair(dataSourceName, "certificate_id", resourceName, "certificate_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "domain", resourceName, "domain"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_notification", resourceName, "user_notification"),
				),
			},
		},
	})
}

func testAccPRAPortalControllerDataSourceConfig(name, description string, enabled, notificationEnabled bool, domainName, userNotification string) string {
	return fmt.Sprintf(`
data "zpa_ba_certificate" "this" {
	name = "pra01.bd-hashicorp.com"
}

resource "zpa_pra_portal_controller" "test" {
	name = "%s"
	description = "%s"
	domain = "%s.bd-hashicorp.com"
	user_notification = "%s"
	enabled = "%s"
	user_notification_enabled = "%s"
	certificate_id = data.zpa_ba_certificate.this.id
}

data "zpa_pra_portal_controller" "test" {
  id = zpa_pra_portal_controller.test.id
}
`, name, description, domainName, userNotification, strconv.FormatBool(enabled), strconv.FormatBool(notificationEnabled))
}
