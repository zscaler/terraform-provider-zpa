package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"
)

func TestDataSourceApplicationSegmentInspection(t *testing.T) {
	match := applicationsegmentinspection.AppSegmentInspection{ID: "ai1", Name: "jenkins.acme.com", Enabled: true, SegmentGroupID: "sg1", DomainNames: []string{"jenkins.acme.com"}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceApplicationSegmentInspection(),
		listPath:  v1Path("application"),
		itemPath:  v1Path("application", "ai1"),
		items:     []applicationsegmentinspection.AppSegmentInspection{{ID: "other", Name: "jenkins.acme.com-old"}, match},
		item:      match,
		name:      "jenkins.acme.com",
		id:        "ai1",
		wantState: map[string]any{"segment_group_id": "sg1", "domain_names": []any{"jenkins.acme.com"}},
	})
}
