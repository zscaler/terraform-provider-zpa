package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

// The C2C IP ranges list is a plain JSON array; the name match is exact.
func TestDataSourceC2CIPRanges(t *testing.T) {
	match := c2c_ip_ranges.IPRanges{ID: "ipr1", Name: "Site A", IpRangeBegin: "192.168.1.1", IpRangeEnd: "192.168.1.254", CountryCode: "US"}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceC2CIPRanges(),
		listPath:  v1Path("v2", "ipRanges"),
		plainList: true,
		itemPath:  v1Path("v2", "ipRanges", "ipr1"),
		items:     []c2c_ip_ranges.IPRanges{{ID: "ipr0", Name: "Site A (old)"}, match},
		item:      match,
		name:      "Site A",
		id:        "ipr1",
		wantState: map[string]any{"ip_range_begin": "192.168.1.1", "country_code": "US"},
	})
}
