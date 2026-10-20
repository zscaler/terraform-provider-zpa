package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_namespace"
)

func TestResourceTagNamespace_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceTagNamespace(),
		config:   map[string]any{"name": "corp", "description": "Corporate", "enabled": true},
		id:       "ns1",
		response: tag_namespace.Namespace{ID: "ns1", Name: "corp", Description: "Corporate", Enabled: true},

		createPath: v1Path("namespace"),
		itemPath:   v1Path("namespace", "ns1"),
		wantBody:   map[string]any{"name": "corp", "enabled": true},
		wantState:  map[string]any{"name": "corp", "description": "Corporate", "enabled": true},
	})
}
