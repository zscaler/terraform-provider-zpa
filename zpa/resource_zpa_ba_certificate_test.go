package zpa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
)

func TestAccResourceBaCertificate_basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPABACertificate) // Random certificate name
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
					testAccCheckBaCertificateExists(resourceTypeAndName),
				),
			},
			{
				ResourceName:            resourceTypeAndName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert_blob"}, // Ignore cert_blob during import verification
			},
		},
	})
}

func generateSelfSignedCert(certName string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"US"},
			Province:           []string{"California"},
			Locality:           []string{"San Jose"},
			Organization:       []string{"BD-HashiCorp"},
			OrganizationalUnit: []string{"ITDepartment"},
			CommonName:         certName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	return derBytes, x509.MarshalPKCS1PrivateKey(priv), nil
}

func pemEncode(derBytes []byte, pemType string) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: pemType, Bytes: derBytes}))
}

func testAccCheckBaCertificateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Find the resource in the state
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		// Assume you have an API client set up and it has a method to get a certificate by ID
		apiClient := testAccProvider.Meta().(*Client)
		_, _, err := apiClient.bacertificate.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching certificate with resource ID [%s] from API: %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccBaCertificateDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPABACertificate {
			continue
		}

		baCertificate, _, err := apiClient.bacertificate.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if baCertificate != nil {
			return fmt.Errorf("browser access certificate with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccBaCertificateConfigure(generatedName, certPEM, privateKeyPEM string) string {
	return fmt.Sprintf(`
resource "zpa_ba_certificate" "%s" {
    name       = "tf-acc-test-%s"
    cert_blob  = <<EOF
%s
%s
EOF
    description = "Self-signed certificate for testing"
}

data "zpa_ba_certificate" "%s" {
    id = zpa_ba_certificate.%s.id
}
`, generatedName, generatedName, privateKeyPEM, certPEM, generatedName, generatedName)
}
