package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/enrollmentcert"
)

func TestDataSourceEnrollmentCert(t *testing.T) {
	match := enrollmentcert.EnrollmentCert{ID: "ec1", Name: "Connector", AllowSigning: true, ClientCertType: "ZAPP_CLIENT"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceEnrollmentCert(),
		listPath:  v2Path("enrollmentCert"),
		itemPath:  v1Path("enrollmentCert", "ec1"),
		items:     []enrollmentcert.EnrollmentCert{{ID: "other", Name: "Connector-old"}, match},
		item:      match,
		name:      "connector",
		id:        "ec1",
		wantState: map[string]any{"allow_signing": true, "client_cert_type": "ZAPP_CLIENT"},
	})
}
