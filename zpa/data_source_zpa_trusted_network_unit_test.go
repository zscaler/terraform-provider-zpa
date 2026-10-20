package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
)

func TestDataSourceTrustedNetwork(t *testing.T) {
	match := trustednetwork.TrustedNetwork{ID: "tn1", Name: "Corp-Trusted-Networks", NetworkID: "net-udid-1", ZscalerCloud: "zscalerthree"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceTrustedNetwork(),
		listPath:  v2Path("network"),
		itemPath:  v1Path("network", "tn1"),
		items:     []trustednetwork.TrustedNetwork{{ID: "other", Name: "Corp-Trusted-Networks-old"}, match},
		item:      match,
		name:      "corp-trusted-networks",
		id:        "tn1",
		wantState: map[string]any{"network_id": "net-udid-1", "zscaler_cloud": "zscalerthree"},
	})
}
