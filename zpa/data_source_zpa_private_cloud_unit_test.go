package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud"
)

func TestDataSourcePrivateCloud(t *testing.T) {
	match := private_cloud.PrivateCloudController{ID: "pc1", Name: "SJC Private Cloud", Enabled: true, ReEnrollPeriod: "30"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePrivateCloud(),
		listPath:  v1Path("privateCloud"),
		itemPath:  v1Path("privateCloud", "pc1"),
		items:     []private_cloud.PrivateCloudController{{ID: "other", Name: "SJC Private Cloud Legacy"}, match},
		item:      match,
		name:      "SJC Private Cloud",
		id:        "pc1",
		wantState: map[string]any{"re_enroll_period": "30", "enabled": true},
	})
}
