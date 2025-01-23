package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourcePRAConsoleController_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAConsoleController)
	domainName := "pra_" + generatedName

	praPortalTypeAndName, _, praPortalGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAPortalController)
	praPortalHCL := testAccCheckPRAPortalControllerConfigure(praPortalTypeAndName, praPortalGeneratedName, variable.PraPortalDescription, variable.PraPortalEnabled, variable.PraUserNotificationEnabled, domainName, variable.PraUserNotification)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAConsoleControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAConsoleControllerConfigure(resourceTypeAndName, generatedName, generatedName, variable.PraConsoleDescription, variable.PraConsoleEnabled, praPortalHCL, praPortalTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PraConsoleEnabled)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "pra_application.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "pra_portals.#", "1"),
				),
			},
		},
	})
}
