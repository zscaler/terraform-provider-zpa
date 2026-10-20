package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func TestResourceCBICertificate_CRUD(t *testing.T) {
	// pem is stored normalized, ending in a single newline
	pem := "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n"
	testCRUD(t, crudSpec{
		resource:   resourceCBICertificates(),
		config:     map[string]any{"name": "Corp Root", "pem": pem},
		id:         "cc1",
		response:   cbicertificatecontroller.CBICertificate{ID: "cc1", Name: "Corp Root", PEM: pem},
		createPath: cbiPath("certificate"),
		itemPath:   cbiPath("certificates", "cc1"),
		wantBody:   map[string]any{"name": "Corp Root", "pem": pem},
		wantState:  map[string]any{"name": "Corp Root", "pem": pem},
	})
}
