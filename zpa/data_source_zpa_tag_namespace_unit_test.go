package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_namespace"
)

func TestDataSourceTagNamespace(t *testing.T) {
	match := tag_namespace.Namespace{ID: "ns1", Name: "Corp", Description: "Corporate", Enabled: true}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceTagNamespace(),
		listPath:   v1Path("namespace", "search"),
		listMethod: "POST",
		itemPath:   v1Path("namespace", "ns1"),
		items:      []tag_namespace.Namespace{{ID: "ns0", Name: "corporate"}, match},
		item:       match,
		name:       "corp",
		id:         "ns1",
		wantState:  map[string]any{"description": "Corporate", "enabled": true},
	})
}
