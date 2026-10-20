package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
)

func TestDataSourceUserPortalAUP(t *testing.T) {
	match := aup.UserPortalAup{ID: "aup1", Name: "Corp AUP", Enabled: true, Aup: "Use responsibly.", Email: "it@acme.com"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceUserPortalAUP(),
		listPath:  v1Path("userportal", "aup"),
		itemPath:  v1Path("userportal", "aup", "aup1"),
		items:     []aup.UserPortalAup{{ID: "other", Name: "Corp AUP Legacy"}, match},
		item:      match,
		name:      "Corp AUP",
		id:        "aup1",
		wantState: map[string]any{"aup": "Use responsibly.", "email": "it@acme.com"},
	})
}
