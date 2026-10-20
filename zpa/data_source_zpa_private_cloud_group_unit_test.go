package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

func TestDataSourcePrivateCloudGroup(t *testing.T) {
	match := private_cloud_group.PrivateCloudGroup{ID: "pcg1", Name: "SJC Private Cloud", Enabled: true, CountryCode: "US", Location: "San Jose, CA, USA"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourcePrivateCloudGroup(),
		listPath:  v1Path("privateCloudControllerGroup"),
		itemPath:  v1Path("privateCloudControllerGroup", "pcg1"),
		items:     []private_cloud_group.PrivateCloudGroup{{ID: "other", Name: "SJC Private Cloud (old)"}, match},
		item:      match,
		name:      "sjc private cloud",
		id:        "pcg1",
		wantState: map[string]any{"country_code": "US", "location": "San Jose, CA, USA"},
	})
}
