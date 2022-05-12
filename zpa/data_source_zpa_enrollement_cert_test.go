package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnrollmentCert_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceEnrollmentCertConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceEnrollmentCertCheck("data.zpa_enrollment_cert.root"),
					testAccDataSourceEnrollmentCertCheck("data.zpa_enrollment_cert.client"),
					testAccDataSourceEnrollmentCertCheck("data.zpa_enrollment_cert.connector"),
					testAccDataSourceEnrollmentCertCheck("data.zpa_enrollment_cert.service_edge"),
				),
			},
		},
	})
}

func testAccDataSourceEnrollmentCertCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceEnrollmentCertConfig_basic = `
data "zpa_enrollment_cert" "root" {
    name = "Root"
}

data "zpa_enrollment_cert" "client" {
    name = "Client"
}

data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}

data "zpa_enrollment_cert" "service_edge" {
    name = "Service Edge"
}
`
