package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccDataSourceAccessPolicyClientTypes_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAccessPolicyClientTypes_basic,
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
