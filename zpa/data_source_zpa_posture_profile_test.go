package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePostureProfile_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourcePostureProfileConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_posture_profile.pre_zta", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_posture_profile.zta_40", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_posture_profile.zta_80", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourcePostureProfileConfig_basic = `
data "zpa_posture_profile" "pre_zta" {
    name = "CrowdStrike_ZPA_Pre-ZTA (zscalerthree.net)"
}
data "zpa_posture_profile" "zta_40" {
    name = "CrowdStrike_ZPA_ZTA_40 (zscalerthree.net)"
}
data "zpa_posture_profile" "zta_80" {
    name = "CrowdStrike_ZPA_ZTA_80 (zscalerthree.net)"
}`
