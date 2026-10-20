package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
)

func TestDataSourceServiceEdgeController(t *testing.T) {
	match := serviceedgecontroller.ServiceEdgeController{ID: "se1", Name: "pse-01", ServiceEdgeGroupID: "seg1", ControlChannelStatus: "ZPN_STATUS_AUTHENTICATED"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceServiceEdgeController(),
		listPath:  v1Path("serviceEdge"),
		itemPath:  v1Path("serviceEdge", "se1"),
		items:     []serviceedgecontroller.ServiceEdgeController{{ID: "other", Name: "pse-01 Legacy"}, match},
		item:      match,
		name:      "pse-01",
		id:        "se1",
		wantState: map[string]any{"service_edge_group_id": "seg1", "control_channel_status": "ZPN_STATUS_AUTHENTICATED"},
	})
}
