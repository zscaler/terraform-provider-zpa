package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var idpNames = []string{
	"SGIO-Admin-Okta", "SGIO-User-Okta",
}

func TestAccDataSourceIdpController_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceIdpController_Basic(),
				Check: resource.ComposeTestCheckFunc(
					generateIdpControllerChecks()...,
				),
			},
		},
	})
}

func generateIdpControllerChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range idpNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_idp_controller.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourceIdpController_Basic() string {
	var configs string
	for _, name := range idpNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_idp_controller" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
