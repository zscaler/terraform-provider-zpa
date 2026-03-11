package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourceTagKey_Basic(t *testing.T) {
	nsResourceTypeAndName, _, nsGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPATagNamespace)
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPATagKey)

	nsName := "tf-acc-test-" + nsGeneratedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTagKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTagKeyConfigure(nsResourceTypeAndName, nsName, resourceTypeAndName, generatedName, variable.TagKeyDescription, variable.TagKeyEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "enabled", resourceTypeAndName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "namespace_id", resourceTypeAndName, "namespace_id"),
				),
			},
		},
	})
}
