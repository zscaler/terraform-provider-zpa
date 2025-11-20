package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

var networkNames = []string{
	"BD-TrustedNetwork03",
	"BDTrustedNetwork",
}

func TestAccDataSourceTrustedNetwork_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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

func createValidResourceName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}
