package zpa

import (
	"testing"

	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestDataSourceWorkloadTagGroup(t *testing.T) {
	match := zpacommon.CommonSummary{ID: "wtg1", Name: "Prod Workloads", Enabled: true}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceWorkloadTagGroup(),
		listPath:   v1Path("workloadTagGroup", "summary"),
		idFromList: true,
		items:      []zpacommon.CommonSummary{{ID: "other", Name: "Prod Workloads Legacy"}, match},
		item:       match,
		name:       "Prod Workloads",
		id:         "wtg1",
		wantState:  map[string]any{"enabled": true},
	})
}
