package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
)

func TestDataSourceApplicationServer(t *testing.T) {
	match := appservercontroller.ApplicationServer{ID: "as1", Name: "server1", Address: "server1.acme.com", Enabled: true, AppServerGroupIds: []string{"sg1"}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceApplicationServer(),
		listPath:  v1Path("server"),
		itemPath:  v1Path("server", "as1"),
		items:     []appservercontroller.ApplicationServer{{ID: "other", Name: "server1 (old)"}, match},
		item:      match,
		name:      "server1",
		id:        "as1",
		wantState: map[string]any{"address": "server1.acme.com", "app_server_group_ids": []any{"sg1"}},
	})
}
