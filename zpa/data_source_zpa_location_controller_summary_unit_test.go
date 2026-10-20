package zpa

import (
	"testing"

	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestDataSourceLocationControllerSummary(t *testing.T) {
	match := zpacommon.CommonSummary{ID: "loc1", Name: "HQ Location", Enabled: true}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceLocationControllerSummary(),
		listPath:   v1Path("location", "summary"),
		idFromList: true,
		items:      []zpacommon.CommonSummary{{ID: "other", Name: "HQ Location Legacy"}, match},
		item:       match,
		name:       "HQ Location",
		id:         "loc1",
		wantState:  map[string]any{"enabled": true},
	})
}
