package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var enrollementCertNames = []string{
	"Root", "Client", "Connector", "Service Edge",
}

func TestAccDataSourceEnrollmentCert_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceEnrollmentCert_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateEnrollmentCertChecks()...,
				),
			},
		},
	})
}

func generateEnrollmentCertChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range enrollementCertNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zpa_enrollment_cert.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourceEnrollmentCert_basic() string {
	var configs string
	for _, name := range enrollementCertNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zpa_enrollment_cert" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
