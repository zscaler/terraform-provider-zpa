package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAccessPolicyClientTypes_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: (testAccCheckDataSourceAccessPolicyClientTypes_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_exporter"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_exporter_noauth"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_browser_isolation"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_machine_tunnel"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_ip_anchoring"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_edge_connector"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_zapp"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_slogger"),
					resource.TestCheckResourceAttrSet("data.zpa_access_policy_client_types.this", "zpn_client_type_branch_connector"),
				),
			},
		},
	})
}

var testAccCheckDataSourceAccessPolicyClientTypes_basic = `
data "zpa_access_policy_client_types" "this" {}`
