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
						"data.zpa_ba_certificate.foobar", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceEnrollmentCertConfig_basic = `
data "zpa_ba_certificate" "foobar" {
    name = "jenkins.securitygeek.io"
}`
