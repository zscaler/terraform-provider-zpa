package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func TestDataSourceServiceEdgeGroup(t *testing.T) {
	match := serviceedgegroup.ServiceEdgeGroup{ID: "seg1", Name: "SJC Service Edges", Enabled: true, CountryCode: "US", IsPublic: "TRUE", UpgradeDay: "SUNDAY"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceServiceEdgeGroup(),
		listPath:  v1Path("serviceEdgeGroup"),
		itemPath:  v1Path("serviceEdgeGroup", "seg1"),
		items:     []serviceedgegroup.ServiceEdgeGroup{{ID: "other", Name: "SJC Service Edges (old)"}, match},
		item:      match,
		name:      "sjc service edges",
		id:        "seg1",
		wantState: map[string]any{"country_code": "US", "upgrade_day": "SUNDAY"},
	})
}
