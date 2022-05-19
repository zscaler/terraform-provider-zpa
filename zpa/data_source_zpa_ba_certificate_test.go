package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBaCertificate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceBaCertificateConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceBaCertificateCheck("data.zpa_ba_certificate.certificate"),
				),
			},
		},
	})
}

func testAccDataSourceBaCertificateCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceBaCertificateConfig_basic = `
data "zpa_ba_certificate" "certificate" {
    name = "jenkins.securitygeek.io"
}`
*/
