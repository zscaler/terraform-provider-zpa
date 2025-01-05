package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourceServerGroup_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerGroupConfigure(resourceTypeAndName, generatedName, generatedName, generatedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServerGroupEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dynamic_discovery", strconv.FormatBool(variable.ServerGroupDynamicDiscovery)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "app_connector_groups.#", "1"),
				),
			},
		},
	})
}
