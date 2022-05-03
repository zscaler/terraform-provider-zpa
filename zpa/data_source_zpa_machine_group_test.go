package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMachineGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceMachineGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceMachineGroupCheck("data.zpa_machine_group.bd_mgr01"),
					testAccDataSourceMachineGroupCheck("data.zpa_machine_group.bd_mgr02"),
				),
			},
		},
	})
}

func testAccDataSourceMachineGroupCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

const testAccCheckDataSourceMachineGroup_basic = `
data "zpa_machine_group" "bd_mgr01" {
    name = "BD-MGR01"
}
data "zpa_machine_group" "bd_mgr02" {
    name = "BD-MGR02"
}
`
