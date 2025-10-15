package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var lssConfigLogTypeFormatsNames = []string{
	"zpn_ast_comprehensive_stats", "zpn_auth_log", "zpn_pbroker_comprehensive_stats", "zpn_ast_auth_log", "zpn_audit_log",
	"zpn_trans_log", "zpn_waf_http_exchanges_log",
}

func TestAccDataSourceLSSLogTypeFormats_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceLSSLogTypeFormats_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateLSSConfigLogTypeFormatChecks()...,
				),
			},
		},
	})
}

func generateLSSConfigLogTypeFormatChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, log_type := range lssConfigLogTypeFormatsNames {
		resourceName := createValidResourceName(log_type)
		checkName := fmt.Sprintf("data.zpa_lss_config_log_type_formats.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "log_type"),
		))
	}
	return checks
}

func testAccCheckDataSourceLSSLogTypeFormats_basic() string {
	var configs string
	for _, log_type := range lssConfigLogTypeFormatsNames {
		resourceName := createValidResourceName(log_type)
		configs += fmt.Sprintf(`
data "zpa_lss_config_log_type_formats" "%s" {
    log_type = "%s"
}
`, resourceName, log_type)
	}
	return configs
}
