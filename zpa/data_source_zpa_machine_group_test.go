package zpa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var machineGroupNames = []string{
	"BD-MGR01", "BD-MGR02", "BD MGR 03", "BD  MGR  04", "BD   MGR   05",
	"BD    MGR06", "BD  MGR 07", "BD  M GR   08", "BD   M   GR 09",
}

func TestAccDataSourceMachineGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

// createValidResourceName converts the given name to a valid Terraform resource name
func createValidResourceName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}
