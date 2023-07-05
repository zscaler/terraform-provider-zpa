package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLSSLogTypeFormats_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceLSSLogTypeFormats_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_trans_log"),
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_auth_log"),
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_ast_auth_log"),
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_http_trans_log"),
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_audit_log"),
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_ast_comprehensive_stats"),
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_sys_auth_log"),
					testAccDataSourceLSSLogTypeFormatsCheck("data.zpa_lss_config_log_type_formats.zpn_waf_http_exchanges_log"),
				),
			},
		},
	})
}

func testAccDataSourceLSSLogTypeFormatsCheck(log_type string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(log_type, "log_type"),
	)
}

var testAccCheckDataSourceLSSLogTypeFormats_basic = `
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

data "zpa_lss_config_log_type_formats" "zpn_waf_http_exchanges_log" {
	log_type = "zpn_waf_http_exchanges_log"
}
`
