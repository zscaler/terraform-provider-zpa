package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
)

func TestAccDataSourcePRAPortalController_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAPortalController)
	domainName := "pra_" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAPortalControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAPortalControllerConfigure(resourceTypeAndName, generatedName, variable.PraPortalDescription, variable.PraPortalEnabled, variable.PraUserNotificationEnabled, domainName, variable.PraUserNotification),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PraPortalEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_notification_enabled", strconv.FormatBool(variable.PraUserNotificationEnabled)),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "certificate_id", resourceTypeAndName, "certificate_id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "domain", resourceTypeAndName, "domain"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "user_notification", resourceTypeAndName, "user_notification"),
				),
			},
		},
	})
}
