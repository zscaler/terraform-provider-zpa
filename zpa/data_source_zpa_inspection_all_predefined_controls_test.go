package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)


func TestAccDataSourceInspectionAllPredefinedControls_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceInspectionAllPredefinedControlsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceInspectionAllPredefinedControlsCheck("data.zpa_inspection_all_predefined_controls.preprocessors"),
					testAccDataSourceInspectionAllPredefinedControlsCheck("data.zpa_inspection_all_predefined_controls.protocol_issues"),
					testAccDataSourceInspectionAllPredefinedControlsCheck("data.zpa_inspection_all_predefined_controls.php_injection"),
				),
			},
		},
	})
}

func testAccDataSourceInspectionAllPredefinedControlsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "version"),
		resource.TestCheckResourceAttrSet(name, "group_name"),
	)
}

var testAccCheckDataSourceInspectionAllPredefinedControlsConfig_basic = `
data "zpa_inspection_all_predefined_controls" "preprocessors" {
	version    = "OWASP_CRS/3.3.0"
	group_name = "Preprocessors"
}

data "zpa_inspection_all_predefined_controls" "protocol_issues" {
	version    = "OWASP_CRS/3.3.0"
	group_name = "Protocol Issues"
}

data "zpa_inspection_all_predefined_controls" "php_injection" {
	version    = "OWASP_CRS/3.3.0"
	group_name = "PHP Injection"
}
`
*/