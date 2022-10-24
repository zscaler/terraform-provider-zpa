package zpa

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccDataSourceServiceEdgeGroup_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServiceEdgeGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceEdgeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, variable.ServiceEdgeDescription, variable.ServiceEdgeEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServiceEdgeEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_public", strconv.FormatBool(variable.ServiceEdgeIsPublic)),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "latitude", resourceTypeAndName, "latitude"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "longitude", resourceTypeAndName, "longitude"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "location", resourceTypeAndName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "version_profile_name", resourceTypeAndName, "version_profile_name"),
				),
			},
		},
	})
}
