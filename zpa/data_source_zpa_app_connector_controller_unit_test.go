package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorcontroller"
)

func TestDataSourceAppConnectorController(t *testing.T) {
	match := appconnectorcontroller.AppConnector{ID: "ac1", Name: "SJC-Connector-01", Enabled: true, AppConnectorGroupID: "acg1", ControlChannelStatus: "ZPN_STATUS_AUTHENTICATED"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceAppConnectorController(),
		listPath:  v1Path("connector"),
		itemPath:  v1Path("connector", "ac1"),
		items:     []appconnectorcontroller.AppConnector{{ID: "other", Name: "SJC-Connector-01 (old)"}, match},
		item:      match,
		name:      "sjc-connector-01",
		id:        "ac1",
		wantState: map[string]any{"app_connector_group_id": "acg1", "control_channel_status": "ZPN_STATUS_AUTHENTICATED"},
	})
}
