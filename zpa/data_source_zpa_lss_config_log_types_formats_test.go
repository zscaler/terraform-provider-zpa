package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLSSLogTypeFormats_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceLSSLogTypeFormatsConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_lss_config_log_type_formats.zpn_trans_log", "log_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_lss_config_log_type_formats.zpn_auth_log", "log_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_lss_config_log_type_formats.zpn_ast_auth_log", "log_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_lss_config_log_type_formats.zpn_http_trans_log", "log_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_lss_config_log_type_formats.zpn_audit_log", "log_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_lss_config_log_type_formats.zpn_ast_comprehensive_stats", "log_type"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_lss_config_log_type_formats.zpn_sys_auth_log", "log_type"),
				),
			},
		},
	})
}

const testAccCheckDataSourceLSSLogTypeFormatsConfig_basic = `
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
	log_type = "zpn_trans_log"
}
data "zpa_lss_config_log_type_formats" "zpn_auth_log" {
	log_type = "zpn_auth_log"
}
data "zpa_lss_config_log_type_formats" "zpn_ast_auth_log" {
	log_type = "zpn_ast_auth_log"
}
data "zpa_lss_config_log_type_formats" "zpn_http_trans_log" {
	log_type = "zpn_http_trans_log"
}
data "zpa_lss_config_log_type_formats" "zpn_audit_log" {
	log_type = "zpn_audit_log"
}
data "zpa_lss_config_log_type_formats" "zpn_ast_comprehensive_stats" {
	log_type = "zpn_ast_comprehensive_stats"
}
data "zpa_lss_config_log_type_formats" "zpn_sys_auth_log" {
	log_type = "zpn_sys_auth_log"
}
`
