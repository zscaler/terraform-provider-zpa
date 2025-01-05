package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
)

func TestAccDataSourceBaCertificate_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPABACertificate)

	cert, privateKey, err := generateSelfSignedCert(generatedName)
	if err != nil {
		t.Fatalf("Error generating self-signed certificate: %v", err)
	}

	certPEM := pemEncode(cert, "CERTIFICATE")
	privateKeyPEM := pemEncode(privateKey, "RSA PRIVATE KEY")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccBaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBaCertificateConfigure(generatedName, certPEM, privateKeyPEM),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
				),
			},
		},
	})
}
