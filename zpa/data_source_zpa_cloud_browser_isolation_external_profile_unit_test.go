package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

func TestDataSourceCBIExternalProfile(t *testing.T) {
	match := cbiprofilecontroller.IsolationProfile{
		ID: "p1", Name: "Isolation profile", Href: "https://isolation.acme.com/p1", IsDefault: true,
	}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceCBIExternalProfile(),
		listPath:  cbiPath("profiles"),
		plainList: true,
		itemPath:  cbiPath("profiles", "p1"),
		items:     []cbiprofilecontroller.IsolationProfile{{ID: "p0", Name: "Isolation profile old"}, match},
		item:      match,
		name:      "isolation profile",
		id:        "p1",
		wantState: map[string]any{"href": "https://isolation.acme.com/p1", "is_default": true},
	})
}
