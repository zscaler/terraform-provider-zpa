package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccDataSourceLSSConfigController_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPALSSController)
	rPort := acctest.RandIntRange(1000, 9999)
	rIP, _ := acctest.RandIpAddress("192.168.100.0/25")

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLSSConfigControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLSSConfigControllerConfigure(resourceTypeAndName, generatedName, generatedName, generatedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, rIP, rPort, variable.LSSControllerEnabled, variable.LSSControllerTLSEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "config.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "config.0.name", resourceTypeAndName, "config.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "config.0.description", resourceTypeAndName, "config.0.description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "config.0.lss_host", resourceTypeAndName, "config.0.lss_host"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "config.0.lss_port", resourceTypeAndName, "config.0.lss_port"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "config.0.source_log_type", resourceTypeAndName, "config.0.source_log_type"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "config.0.enabled", strconv.FormatBool(variable.LSSControllerEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "config.0.use_tls", strconv.FormatBool(variable.LSSControllerTLSEnabled)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "connector_groups.#", "1"),
				),
			},
		},
	})
}
