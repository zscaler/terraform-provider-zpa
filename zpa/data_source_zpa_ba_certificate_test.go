package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBaCertificate_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceBaCertificateConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zpa_ba_certificate.foobar", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceBaCertificateConfig_basic = `
data "zpa_ba_certificate" "foobar" {
    name = "jenkins.securitygeek.io"
}`
