package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
)

func TestDataSourceApplicationSegmentBrowserAccess(t *testing.T) {
	match := applicationsegmentbrowseraccess.BrowserAccess{ID: "ba1", Name: "CRM portal", Enabled: true, SegmentGroupID: "sg1", DomainNames: []string{"crm.acme.com"}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceApplicationSegmentBrowserAccess(),
		listPath:  v1Path("application"),
		itemPath:  v1Path("application", "ba1"),
		items:     []applicationsegmentbrowseraccess.BrowserAccess{{ID: "other", Name: "CRM portal-old"}, match},
		item:      match,
		name:      "crm portal",
		id:        "ba1",
		wantState: map[string]any{"segment_group_id": "sg1", "domain_names": []any{"crm.acme.com"}},
	})
}
