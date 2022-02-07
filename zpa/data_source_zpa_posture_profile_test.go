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
						"data.zpa_posture_profile.foobar", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourcePostureProfileConfig_basic = `
data "zpa_posture_profile" "foobar" {
    name = "CrowdStrike_ZPA_ZTA_40"
}`
