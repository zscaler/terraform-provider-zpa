package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/isolationprofile"
)

var isolationProfiles = []isolationprofile.IsolationProfile{
	{ID: "iso0", Name: "Isolation Profile Legacy"},
	{ID: "iso1", Name: "Isolation Profile", Enabled: true, IsolationURL: "https://isolation.acme.com"},
}

func TestDataSourceIsolationProfile(t *testing.T) {
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceIsolationProfile(),
		listPath:  v1Path("isolation", "profiles"),
		items:     isolationProfiles,
		name:      "isolation profile",
		id:        "iso1",
		wantState: map[string]any{"isolation_url": "https://isolation.acme.com"},
	})
}

// There is no single-profile endpoint: the ID lookup searches the list.
func TestDataSourceIsolationProfile_ByID(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("isolation", "profiles"), common.SuccessResponse(common.ZPAList(isolationProfiles)))

	d := fx.resourceData(dataSourceIsolationProfile(), "iso1", map[string]any{})
	fx.run(dataSourceIsolationProfileRead, d)
	assertState(t, d, map[string]any{"name": "Isolation Profile", "isolation_url": "https://isolation.acme.com"})
}
