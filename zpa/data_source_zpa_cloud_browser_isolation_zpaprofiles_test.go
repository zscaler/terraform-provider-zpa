package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var zpaCbiProfileNames = []string{
	"BD_SA_Profile1", "BD_SA_Profile2",
}

func TestAccDataSourceCBIZPAProfiles_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCBIZPAProfiles_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateCBIZPAProfilesChecks()...,
				),
			},
		},
	})
}

func generateCBIZPAProfilesChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range zpaCbiProfileNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_cloud_browser_isolation_zpa_profile.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourceCBIZPAProfiles_basic() string {
	var configs string
	for _, name := range zpaCbiProfileNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_cloud_browser_isolation_zpa_profile" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
