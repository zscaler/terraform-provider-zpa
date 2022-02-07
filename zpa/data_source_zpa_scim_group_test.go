package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScimGroup_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceScimGroupConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_scim_groups.engineering", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_scim_groups.contractors", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceScimGroupConfig_basic = `
data "zpa_scim_groups" "engineering" {
    name = "Engineering"
	idp_name = "SGIO-User-Okta"
}
data "zpa_scim_groups" "contractors" {
    name = "Contractors"
	idp_name = "SGIO-User-Okta"
}
`
