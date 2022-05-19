package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePostureProfile_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourcePostureProfileConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePostureProfileCheck("data.zpa_posture_profile.pre_zta"),
					testAccDataSourcePostureProfileCheck("data.zpa_posture_profile.zta_40"),
					testAccDataSourcePostureProfileCheck("data.zpa_posture_profile.zta_80"),
				),
			},
		},
	})
}

func testAccDataSourcePostureProfileCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourcePostureProfileConfig_basic = `
data "zpa_posture_profile" "pre_zta" {
    name = "CrowdStrike_ZPA_Pre-ZTA (zscalertwo.net)"
}
data "zpa_posture_profile" "zta_40" {
    name = "CrowdStrike_ZPA_ZTA_40 (zscalertwo.net)"
}
data "zpa_posture_profile" "zta_80" {
    name = "CrowdStrike_ZPA_ZTA_80 (zscalertwo.net)"
}`
