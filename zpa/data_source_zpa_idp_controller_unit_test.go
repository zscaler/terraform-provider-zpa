package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
)

func TestDataSourceIdpController(t *testing.T) {
	match := idpcontroller.IdpController{ID: "idp1", Name: "Okta_Users", Enabled: true, SsoType: []string{"USER"}, Domainlist: []string{"acme.com"}, CreationTime: "1700000000", ModifiedTime: "1700000000"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceIdpController(),
		listPath:  v2Path("idp"),
		itemPath:  v1Path("idp", "idp1"),
		items:     []idpcontroller.IdpController{{ID: "other", Name: "Okta_Users-old"}, match},
		item:      match,
		name:      "Okta_Users", // the SDK matches IdP names exactly
		id:        "idp1",
		wantState: map[string]any{"sso_type": []any{"USER"}, "domain_list": []any{"acme.com"}},
	})
}
