package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourceUserPortalLink_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAUserPortalLink)

	userPortalControllerTypeAndName, _, userPortalControllerGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAUserPortalController)
	userPortalControllerHCL := testAccCheckUserPortalControllerConfigure(userPortalControllerTypeAndName, userPortalControllerGeneratedName, variable.UserPortalLinkDescription, variable.UserPortalLinkEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserPortalLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserPortalLinkConfigure(resourceTypeAndName, generatedName, generatedName, generatedName, userPortalControllerHCL, userPortalControllerTypeAndName, variable.UserPortalLinkEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "enabled", resourceTypeAndName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "link", resourceTypeAndName, "link"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "icon_text", resourceTypeAndName, "icon_text"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "protocol", resourceTypeAndName, "protocol"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.UserPortalLinkEnabled)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "user_portals.#", "1"),
				),
			},
		},
	})
}
