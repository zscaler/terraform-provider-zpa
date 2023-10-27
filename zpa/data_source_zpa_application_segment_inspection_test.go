package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
)

func TestAccDataSourceApplicationSegmentInspection_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationSegmentInspection)
	rDomain := acctest.RandomWithPrefix("tf-acc-test")

	serverGroupTypeAndName, _, serverGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)
	serverGroupHCL := testAccCheckServerGroupConfigure(serverGroupTypeAndName, serverGroupGeneratedName, serverGroupGeneratedName, serverGroupGeneratedName, serverGroupGeneratedName, "", variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationSegmentInspectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationSegmentInspectionConfigure(resourceTypeAndName, generatedName, generatedName, generatedName, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName, variable.AppSegmentEnabled, rDomain, variable.AppSegmentCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "bypass_type", resourceTypeAndName, "bypass_type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "health_reporting", resourceTypeAndName, "health_reporting"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "segment_group_id", resourceTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.AppSegmentEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.AppSegmentCnameEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "common_apps_dto.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "tcp_port_ranges.#", "2"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
