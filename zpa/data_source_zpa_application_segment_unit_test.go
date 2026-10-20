package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

func TestDataSourceApplicationSegment(t *testing.T) {
	match := applicationsegment.ApplicationSegmentResource{ID: "2001", Name: "crm", Enabled: true, SegmentGroupID: "1001", DomainNames: []string{"crm.acme.com"}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceApplicationSegment(),
		listPath:  v1Path("application"),
		itemPath:  v1Path("application", "2001"),
		items:     []applicationsegment.ApplicationSegmentResource{{ID: "other", Name: "crm-old"}, match},
		item:      match,
		name:      "crm",
		id:        "2001",
		wantState: map[string]any{"segment_group_id": "1001", "domain_names": []any{"crm.acme.com"}},
	})
}
