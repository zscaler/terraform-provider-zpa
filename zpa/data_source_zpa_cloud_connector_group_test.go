package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCloudConnectorGroup_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceCloudConnectorGroupConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_cloud_connector_group.foobar", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceCloudConnectorGroupConfig_basic = `
data "zpa_cloud_connector_group" "foobar" {
    name = "zs-cc-vpc-096108eb5d9e68d71-ca-central-1a"
}`
