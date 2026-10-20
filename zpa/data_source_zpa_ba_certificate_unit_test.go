package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
)

// The name lookup searches the issued certificates (v2); the ID lookup reads the certificate (v1).
func TestDataSourceBaCertificate(t *testing.T) {
	match := bacertificate.BaCertificate{
		ID: "c1", Name: "portal.acme.com", CName: "portal.acme.com", IssuedTo: "CN=portal.acme.com",
		CreationTime: "1700000000", ModifiedTime: "1700000000", ValidFromInEpochSec: "1700000000", ValidToInEpochSec: "1731536000",
	}
	other := match
	other.ID, other.Name = "c0", "portal.acme.com.old"
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceBaCertificate(),
		listPath:  v2Path("clientlessCertificate", "issued"),
		itemPath:  v1Path("certificate", "c1"),
		items:     []bacertificate.BaCertificate{other, match},
		item:      match,
		name:      "portal.acme.com",
		id:        "c1",
		wantState: map[string]any{"cname": "portal.acme.com", "issued_to": "CN=portal.acme.com"},
	})
}
