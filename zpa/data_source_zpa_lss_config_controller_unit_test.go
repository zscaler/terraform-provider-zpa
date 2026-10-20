package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

// LSS configurations are matched on the name inside their config block.
func TestDataSourceLSSConfigController(t *testing.T) {
	lss := func(id, name string) lssconfigcontroller.LSSResource {
		return lssconfigcontroller.LSSResource{
			ID:              id,
			LSSConfig:       &lssconfigcontroller.LSSConfig{ID: id, Name: name, LSSHost: "splunk.acme.com", LSSPort: "5001", SourceLogType: "zpn_trans_log"},
			ConnectorGroups: []lssconfigcontroller.ConnectorGroups{{ID: "acg1"}},
		}
	}
	match := lss("lss1", "SIEM user activity")
	testDataSource(t, dataSourceSpec{
		resource: dataSourceLSSConfigController(),
		listPath: v2Path("lssConfig"),
		itemPath: v2Path("lssConfig", "lss1"),
		items:    []lssconfigcontroller.LSSResource{lss("lss0", "SIEM user activity old"), match},
		item:     match,
		name:     "siem user activity",
		// LSS configurations are named by their config block
		fullFields: map[string]any{"config": map[string]any{"name": "siem user activity", "lssHost": "splunk.acme.com"}},
		id:         "lss1",
		wantState: map[string]any{
			"config": []any{map[string]any{"lss_host": "splunk.acme.com", "source_log_type": "zpn_trans_log"}},
		},
	})
}
