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
				Config: testAccCheckDataSourceInspectionPredefinedControls_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.example1"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.example2"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.example3"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.example4"),
					testAccDataSourceInspectionPredefinedControlsCheck("data.zpa_inspection_predefined_controls.example5"),
				),
			},
		},
	})
}

func testAccDataSourceInspectionPredefinedControlsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceInspectionPredefinedControls_basic = `
data "zpa_inspection_predefined_controls" "example1" {
    name = "Failed to parse request body"
	version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "example2" {
    name = "Multipart request body failed strict validation"
	version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "example3" {
    name = "Multipart parser detected a possible unmatched boundary"
	version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "example4" {
    name = "Attempted multipart/form-data bypass"
	version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_predefined_controls" "example5" {
    name = "GET or HEAD Request with Body Content"
	version = "OWASP_CRS/3.3.0"
}
`
