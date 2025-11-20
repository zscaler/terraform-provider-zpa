package datasources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

var idpNames = []string{
	"BD_Okta_Users",
}

func TestAccDataSourceIdpController_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
