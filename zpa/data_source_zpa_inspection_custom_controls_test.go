package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
)

func TestAccDataSourceInspectionCustomControls_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionCustomControl)

	// appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	// appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInspectionCustomControlsDestroy,
		Steps: []resource.TestStep{

			{
				Config: testAccCheckInspectionCustomControlsConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "action", resourceTypeAndName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "default_action", resourceTypeAndName, "default_action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "paranoia_level", resourceTypeAndName, "paranoia_level"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "severity", resourceTypeAndName, "severity"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "rules.#", "1"),
				),
			},
		},
	})
}
