package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceInspectionPredefinedControls_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceInspectionPredefinedControlsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.control01"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.control02"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.control03"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.control04"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.control05"),
				),
			},
		},
	})
}

func testAccDataSourceInspectionPredefinedControlsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "name"),
		resource.TestCheckResourceAttrSet(name, "version"),
	)
}

var testAccCheckDataSourceInspectionPredefinedControlsConfig_basic = `
data "zpa_inspection_predefined_controls" "control01" {
    name = "Failed to parse request body"
    version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "control02" {
    name = "Multipart request body failed strict validation"
    version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "control03" {
    name = "Multipart parser detected a possible unmatched boundary"
    version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "control04" {
    name = "Attempted multipart/form-data bypass"
    version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "control05" {
    name = "GET or HEAD Request with Body Content"
    version = "OWASP_CRS/3.3.0"
}
`
