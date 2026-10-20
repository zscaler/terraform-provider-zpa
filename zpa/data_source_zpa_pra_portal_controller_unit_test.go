package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

func TestDataSourcePRAPortalController(t *testing.T) {
	match := praportal.PRAPortal{ID: "pp1", Name: "pra.acme.com", Enabled: true, Domain: "pra.acme.com", CertificateID: "c1"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePRAPortalController(),
		listPath:  v1Path("praPortal"),
		itemPath:  v1Path("praPortal", "pp1"),
		items:     []praportal.PRAPortal{{ID: "other", Name: "pra.acme.com-old"}, match},
		item:      match,
		name:      "pra.acme.com",
		id:        "pp1",
		wantState: map[string]any{"domain": "pra.acme.com", "certificate_id": "c1"},
	})
}
