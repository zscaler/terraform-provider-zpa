package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourceCBIBanners_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPACBIBannerController)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCBIBannerConfigure(resourceTypeAndName, generatedName, variable.PrimaryColor, variable.TextColor, variable.NotificationTitle, variable.NotificationText, variable.Banner, variable.Persist, variable.Logo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "primary_color", resourceTypeAndName, "primary_color"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "text_color", resourceTypeAndName, "text_color"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "notification_title", resourceTypeAndName, "notification_title"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "notification_text", resourceTypeAndName, "notification_text"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "logo", resourceTypeAndName, "logo"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "banner", strconv.FormatBool(variable.Banner)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "persist", strconv.FormatBool(variable.Persist)),
				),
			},
		},
	})
}
