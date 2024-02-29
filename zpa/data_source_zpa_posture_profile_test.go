package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var profileNames = []string{
	"CrowdStrike_ZPA_Pre-ZTA", "CrowdStrike_ZPA_ZTA_40", "CrowdStrike_ZPA_ZTA_80",
}

func TestAccDataSourcePostureProfile_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourcePostureProfile_basic(),
				Check: resource.ComposeTestCheckFunc(
					generatePostureProfileChecks()...,
				),
			},
		},
	})
}

func generatePostureProfileChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range profileNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_posture_profile.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourcePostureProfile_basic() string {
	var configs string
	for _, name := range profileNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_posture_profile" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
