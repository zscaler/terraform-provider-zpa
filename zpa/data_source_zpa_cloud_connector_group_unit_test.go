package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector_group"
)

func TestDataSourceCloudConnectorGroup(t *testing.T) {
	match := cloud_connector_group.CloudConnectorGroup{ID: "ccg1", Name: "AWS Cloud Connectors", Enabled: true, ZiaCloud: "zscalerthree", CloudConnectors: []cloud_connector_group.CloudConnectors{{ID: "cc1", Name: "cc-01"}}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceCloudConnectorGroup(),
		listPath:  v1Path("cloudConnectorGroup"),
		itemPath:  v1Path("cloudConnectorGroup", "ccg1"),
		items:     []cloud_connector_group.CloudConnectorGroup{{ID: "other", Name: "AWS Cloud Connectors-old"}, match},
		item:      match,
		name:      "aws cloud connectors",
		id:        "ccg1",
		wantState: map[string]any{"zia_cloud": "zscalerthree", "cloud_connectors": []any{map[string]any{"id": "cc1", "name": "cc-01"}}},
	})
}
