package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
)

func TestDataSourceApplicationSegmentPRA(t *testing.T) {
	match := applicationsegmentpra.AppSegmentPRA{ID: "pra1", Name: "PRA servers", Enabled: true, SegmentGroupID: "sg1", DomainNames: []string{"rdp.acme.com"}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceApplicationSegmentPRA(),
		listPath:  v1Path("application"),
		itemPath:  v1Path("application", "pra1"),
		items:     []applicationsegmentpra.AppSegmentPRA{{ID: "other", Name: "PRA servers-old"}, match},
		item:      match,
		name:      "pra servers",
		id:        "pra1",
		wantState: map[string]any{"segment_group_id": "sg1", "domain_names": []any{"rdp.acme.com"}},
	})
}
