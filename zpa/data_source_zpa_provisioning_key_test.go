package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
)

func TestAccDataSourceProvisioningKey_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAProvisioningKey)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConnectorGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config: testAccCheckProvisioningKeyAppConnectorGroupConfigure(resourceTypeAndName, generatedName, generatedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ConnectorGroupType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "association_type", resourceTypeAndName, "association_type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "max_usage", resourceTypeAndName, "max_usage"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "enrollment_cert_id", resourceTypeAndName, "enrollment_cert_id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "zcomponent_id", resourceTypeAndName, "zcomponent_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ProvisioningKeyEnabled)),
				),
			},
		},
	})
}
