package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnrollmentCert_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceEnrollmentCertConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_enrollment_cert.root", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_enrollment_cert.client", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_enrollment_cert.connector", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_enrollment_cert.service_edge", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zpa_enrollment_cert.isolation_client", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceEnrollmentCertConfig_basic = `
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
data "zpa_enrollment_cert" "isolation_client" {
    name = "Isolation Client"
}
`
