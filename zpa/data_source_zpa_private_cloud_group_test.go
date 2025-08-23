package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourcePrivateCloudGroup_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPrivateCloudGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateCloudGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPrivateCloudGroupConfigure(resourceTypeAndName, generatedName, variable.PrivateCloudGroupDescription, variable.PrivateCloudGroupEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "enabled", resourceTypeAndName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "city_country", resourceTypeAndName, "city_country"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "latitude", resourceTypeAndName, "latitude"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "longitude", resourceTypeAndName, "longitude"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "location", resourceTypeAndName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "upgrade_day", resourceTypeAndName, "upgrade_day"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "upgrade_time_in_secs", resourceTypeAndName, "upgrade_time_in_secs"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "site_id", resourceTypeAndName, "site_id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "version_profile_id", resourceTypeAndName, "version_profile_id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "override_version_profile", resourceTypeAndName, "override_version_profile"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "is_public", resourceTypeAndName, "is_public"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PrivateCloudGroupEnabled)),
				),
			},
		},
	})
}
