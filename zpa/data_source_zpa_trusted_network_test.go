package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var networkNames = []string{
	"BD-TrustedNetwork01", "BD-TrustedNetwork02", "BD-TrustedNetwork03", "BD Trusted Network 01",
}

func TestAccDataSourceTrustedNetwork_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceTrustedNetwork_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateTrustedNetworkChecks()...,
				),
			},
		},
	})
}

func generateTrustedNetworkChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range networkNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_trusted_network.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourceTrustedNetwork_basic() string {
	var configs string
	for _, name := range networkNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_trusted_network" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
