package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCloudConnectorGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCloudConnectorGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCloudConnectorGroupCheck("data.zpa_cloud_connector_group.zs-cc-vpc"),
				),
			},
		},
	})
}

func testAccDataSourceCloudConnectorGroupCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceCloudConnectorGroupConfig_basic = `
data "zpa_cloud_connector_group" "zs-cc-vpc" {
    name = "zs-cc-vpc-096108eb5d9e68d71-ca-central-1a"
}`
*/
