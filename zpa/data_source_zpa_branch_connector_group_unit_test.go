package zpa

import (
	"testing"

	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestDataSourceBranchConnectorGroup(t *testing.T) {
	match := zpacommon.CommonSummary{ID: "bcg1", Name: "Branch Connectors", Enabled: true}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceBranchConnectorGroup(),
		listPath:   v1Path("branchConnectorGroup", "summary"),
		idFromList: true,
		items:      []zpacommon.CommonSummary{{ID: "other", Name: "Branch Connectors Legacy"}, match},
		item:       match,
		name:       "Branch Connectors",
		id:         "bcg1",
		wantState:  map[string]any{"enabled": true},
	})
}
