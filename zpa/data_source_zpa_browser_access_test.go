package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccDataSourceBrowserAccess_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPABrowserAccess)
	rPort := acctest.RandIntRange(1000, 9999)

	serverGroupTypeAndName, _, serverGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)
	serverGroupHCL := testAccCheckServerGroupConfigure(serverGroupTypeAndName, serverGroupGeneratedName, serverGroupGeneratedName, serverGroupGeneratedName, serverGroupGeneratedName, "", variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBrowserAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBrowserAccessConfigure(resourceTypeAndName, generatedName, generatedName, generatedName, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName, rPort, variable.BrowserAccessEnabled, variable.BrowserAccessCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "bypass_type", resourceTypeAndName, "bypass_type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "health_reporting", resourceTypeAndName, "health_reporting"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "segment_group_id", resourceTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.BrowserAccessEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.BrowserAccessCnameEnabled)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "clientless_apps.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "udp_port_ranges.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "tcp_port_ranges.#", "2"),
					// resource.TestCheckResourceAttr(dataSourceTypeAndName, "server_groups.#", "1"),
				),
			},
		},
	})
}
