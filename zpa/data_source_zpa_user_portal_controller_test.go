package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourceUserPortalController_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAUserPortalController)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserPortalControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserPortalControllerConfigure(resourceTypeAndName, generatedName, variable.UserPortalControllerDescription, variable.UserPortalControllerEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "enabled", resourceTypeAndName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "user_notification", resourceTypeAndName, "user_notification"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "user_notification_enabled", resourceTypeAndName, "user_notification_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "certificate_id", resourceTypeAndName, "certificate_id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "domain", resourceTypeAndName, "domain"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.UserPortalControllerEnabled)),
				),
			},
		},
	})
}
