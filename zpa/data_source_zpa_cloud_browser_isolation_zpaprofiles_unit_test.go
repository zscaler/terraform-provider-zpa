package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbizpaprofile"
)

// ZPA CBI profiles are looked up by name only.
func TestDataSourceCBIZPAProfile(t *testing.T) {
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceCBIZPAProfiles(),
		listPath:  cbiPath("zpaprofiles"),
		plainList: true,
		items: []cbizpaprofile.ZPAProfiles{
			{ID: "z0", Name: "ZPA Profile Old"},
			{ID: "z1", Name: "ZPA Profile", Enabled: true, CBIURL: "https://cbi.acme.com"},
		},
		name:      "zpa profile",
		id:        "z1",
		wantState: map[string]any{"enabled": true, "cbi_url": "https://cbi.acme.com"},
	})
}
