package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

func TestResourceC2CIPRanges_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceC2CIPRanges(),
		config: map[string]any{
			"name":           "Site A",
			"enabled":        true,
			"ip_range_begin": "192.168.1.1",
			"ip_range_end":   "192.168.1.254",
			"subnet_cidr":    "/24",
			"location":       "San Jose, CA, USA",
			"country_code":   "US",
		},
		id: "ipr1",
		response: c2c_ip_ranges.IPRanges{
			ID: "ipr1", Name: "Site A", Enabled: true, IpRangeBegin: "192.168.1.1", IpRangeEnd: "192.168.1.254",
			SubnetCidr: "/24", Location: "San Jose, CA, USA", CountryCode: "US",
		},
		// the C2C IP ranges API is the v2 path under the v1 customer base
		createPath: v1Path("v2", "ipRanges"),
		itemPath:   v1Path("v2", "ipRanges", "ipr1"),
		wantBody: map[string]any{
			"name": "Site A", "ipRangeBegin": "192.168.1.1", "ipRangeEnd": "192.168.1.254",
			"subnetCidr": "/24", "countryCode": "US",
		},
		wantState: map[string]any{"ip_range_begin": "192.168.1.1", "location": "San Jose, CA, USA"},
	})
}
