package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
)

func TestDataSourceAppConnectorGroup(t *testing.T) {
	match := appconnectorgroup.AppConnectorGroup{ID: "acg1", Name: "SJC Connectors", Enabled: true, CountryCode: "US", Location: "San Jose, CA, USA", UpgradeDay: "SUNDAY"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceAppConnectorGroup(),
		listPath:  v1Path("appConnectorGroup"),
		itemPath:  v1Path("appConnectorGroup", "acg1"),
		items:     []appconnectorgroup.AppConnectorGroup{{ID: "other", Name: "SJC Connectors (old)"}, match},
		item:      match,
		name:      "sjc connectors",
		id:        "acg1",
		wantState: map[string]any{"country_code": "US", "upgrade_day": "SUNDAY"},
	})
}
