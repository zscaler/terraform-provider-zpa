package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
)

func TestAccDataSourceCBIIsolationProfile_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPACBIExternalIsolationProfile)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBIIsolationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCBIIsolationProfileConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "user_experience.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "security_controls.#", "1"),
				),
			},
		},
	})
}
*/
