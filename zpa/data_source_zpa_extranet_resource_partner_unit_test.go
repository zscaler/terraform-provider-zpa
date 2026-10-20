package zpa

import (
	"testing"

	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestDataSourceExtranetResourcePartner(t *testing.T) {
	match := zpacommon.CommonSummary{ID: "er1", Name: "Partner Network", Enabled: true}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceExtranetResourcePartner(),
		listPath:   v1Path("extranetResource", "partner"),
		idFromList: true,
		items:      []zpacommon.CommonSummary{{ID: "other", Name: "Partner Network Legacy"}, match},
		item:       match,
		name:       "Partner Network",
		id:         "er1",
		wantState:  map[string]any{"enabled": true},
	})
}
