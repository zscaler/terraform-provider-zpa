package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func TestDataSourceCBICertificate(t *testing.T) {
	match := cbicertificatecontroller.CBICertificate{ID: "cc1", Name: "Zscaler Root Certificate", PEM: "-----BEGIN CERTIFICATE-----", IsDefault: true}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceCBICertificates(),
		listPath:  cbiPath("certificates"),
		plainList: true,
		itemPath:  cbiPath("certificates", "cc1"),
		items:     []cbicertificatecontroller.CBICertificate{{ID: "cc0", Name: "Zscaler Root Certificate (old)"}, match},
		item:      match,
		name:      "zscaler root certificate",
		id:        "cc1",
		wantState: map[string]any{"is_default": true},
	})
}
