package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourceC2CIPRanges_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAC2CIPRanges)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckC2CIPRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckC2CIPRangesConfigure(resourceTypeAndName, generatedName, variable.C2CIPRangesDescription, variable.C2CIPRangesEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "enabled", resourceTypeAndName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "location_hint", resourceTypeAndName, "location_hint"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ip_range_begin", resourceTypeAndName, "ip_range_begin"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ip_range_end", resourceTypeAndName, "ip_range_end"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "location", resourceTypeAndName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "sccm_flag", resourceTypeAndName, "sccm_flag"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "country_code", resourceTypeAndName, "country_code"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "latitude_in_db", resourceTypeAndName, "latitude_in_db"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "longitude_in_db", resourceTypeAndName, "longitude_in_db"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.C2CIPRangesEnabled)),
				),
			},
		},
	})
}
