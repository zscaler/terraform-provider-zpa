package datasources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

var machineGroupNames = []string{
	"BD-MGR01", "BD-MGR02", "BD MGR 03",
}

func TestAccDataSourceMachineGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceMachineGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateMachineGroupChecks()...,
				),
			},
		},
	})
}

func generateMachineGroupChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range machineGroupNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_machine_group.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourceMachineGroup_basic() string {
	var configs string
	for _, name := range machineGroupNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_machine_group" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
