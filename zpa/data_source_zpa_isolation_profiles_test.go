package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIsolationRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceIsolationRuleConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceIsolationRuleCheck("data.zpa_isolation_profile.bd_sa_profile1"),
					testAccDataSourceIsolationRuleCheck("data.zpa_isolation_profile.bd_sa_profile2"),
				),
			},
		},
	})
}

func testAccDataSourceIsolationRuleCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceIsolationRuleConfig_basic = `
data "zpa_isolation_profile" "bd_sa_profile1" {
    name = "BD_SA_Profile1"
}

data "zpa_isolation_profile" "bd_sa_profile2" {
    name = "BD_SA_Profile2"
}
`
*/
