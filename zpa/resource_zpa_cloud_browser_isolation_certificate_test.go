package zpa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func TestAccResourceCBICertificate_basic(t *testing.T) {
	var cbiCertificate cbicertificatecontroller.CBICertificate
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPACBICertificate)

	certPEM, err := generateCBIRootCACert()
	if err != nil {
		t.Fatalf("Error generating root CA certificate: %v", err)
	}

	initialCertName := "tf-acc-test-" + generatedName
	updatedCertName := "updated-tf-acc-test-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBICertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCBICertificateConfigure(resourceTypeAndName, initialCertName, string(certPEM)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBICertificateExists(resourceTypeAndName, &cbiCertificate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialCertName),
				),
			},
			{
				Config: testAccCheckCBICertificateConfigure(resourceTypeAndName, updatedCertName, string(certPEM)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBICertificateExists(resourceTypeAndName, &cbiCertificate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedCertName),
				),
			},
		},
	})
}

func generateCBIRootCACert() ([]byte, error) {
	// Generate a private key for root certificate
	priv, err := rsa.GenerateKey(rand.Reader, 4096) // 4096-bit key
	if err != nil {
		return nil, err
	}

	// Create root certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(1024 * 24 * time.Hour) // 1024 days validity

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"US"},
			Province:           []string{"California"},
			Locality:           []string{"San Jose"},
			Organization:       []string{"BD-HashiCorp"},
			OrganizationalUnit: []string{"ITDepartment"},
			CommonName:         "bd-hashicorp.com",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create root certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	// Encode root certificate in PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	return certPEM, nil
}

func testAccCheckCBICertificateExists(resource string, certificate *cbicertificatecontroller.CBICertificate) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedCertificate, _, err := apiClient.cbicertificatecontroller.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*certificate = *receivedCertificate

		return nil
	}
}

func testAccCheckCBICertificateDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPACBICertificate {
			continue
		}

		cbiCertificate, _, err := apiClient.cbicertificatecontroller.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if cbiCertificate != nil {
			return fmt.Errorf("cbi certificate with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCBICertificateConfigure(resourceTypeAndName, certName, certPEM string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
    name       = "%s"
    pem        = <<EOF
%s
EOF
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the certificate
		resourcetype.ZPACBICertificate, resourceName,
		certName, certPEM,

		// Data source type and name
		resourcetype.ZPACBICertificate, resourceName,

		// Reference to the resource
		resourcetype.ZPACBICertificate, resourceName,
	)
}
