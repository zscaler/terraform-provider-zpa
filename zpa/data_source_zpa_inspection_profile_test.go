package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
)

func TestAccDataSourceInspectionProfile_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionCustomControl)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInspectionProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config: testAccCheckInspectionProfileConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "action", resourceTypeAndName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "default_action", resourceTypeAndName, "default_action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "paranoia_level", resourceTypeAndName, "paranoia_level"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "severity", resourceTypeAndName, "severity"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "rules.#", "2"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
*/
