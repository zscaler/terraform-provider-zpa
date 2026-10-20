package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_controller"
)

func TestDataSourcePrivateCloudController(t *testing.T) {
	match := private_cloud_controller.PrivateCloudController{ID: "pcc1", Name: "pcc-01", ControlChannelStatus: "ZPN_STATUS_AUTHENTICATED", CurrentVersion: "24.1"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePrivateCloudController(),
		listPath:  v1Path("privateCloudController"),
		itemPath:  v1Path("privateCloudController", "pcc1"),
		items:     []private_cloud_controller.PrivateCloudController{{ID: "other", Name: "pcc-01 Legacy"}, match},
		item:      match,
		name:      "pcc-01",
		id:        "pcc1",
		wantState: map[string]any{"control_channel_status": "ZPN_STATUS_AUTHENTICATED", "current_version": "24.1"},
	})
}
