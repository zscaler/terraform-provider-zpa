package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var regionNames = []string{
	"Frankfurt",
}

func TestAccDataSourceCBIRegions_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCBIRegions_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateCBIRegionsChecks()...,
				),
			},
		},
	})
}

func generateCBIRegionsChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range regionNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_cloud_browser_isolation_region.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourceCBIRegions_basic() string {
	var configs string
	for _, name := range regionNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_cloud_browser_isolation_region" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
