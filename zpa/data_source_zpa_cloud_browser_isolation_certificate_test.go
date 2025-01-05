package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
)

func TestAccDataSourceBaCertificates_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPACBICertificate)

	certPEM, err := generateCBIRootCACert()
	if err != nil {
		t.Fatalf("Error generating root CA certificate: %v", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBICertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCBICertificateConfigure(resourceTypeAndName, generatedName, string(certPEM)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "pem", resourceTypeAndName, "pem"),
				),
			},
		},
	})
}
