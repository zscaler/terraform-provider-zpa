package zpa

import (
	"testing"

	tag_key_controller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_key"
)

func TestResourceTagKey_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceTagKey(),
		config: map[string]any{
			"name":         "env",
			"namespace_id": "ns1",
			"enabled":      true,
			"tag_values":   []any{map[string]any{"name": "prod"}},
		},
		id: "tk1",
		response: tag_key_controller.TagKey{
			ID: "tk1", Name: "env", NamespaceID: "ns1", Enabled: true,
			TagValues: []tag_key_controller.TagValue{{ID: "tv1", Name: "prod"}},
		},
		// tag keys live under their namespace, and are imported as
		// <namespace_id>/<tag key ID or name> (both looked up through the API)
		importPrefix: "ns1/",
		importLookup: true,
		createPath:   v1Path("namespace", "ns1", "tagKey"),
		itemPath:     v1Path("namespace", "ns1", "tagKey", "tk1"),
		wantBody: map[string]any{
			"name": "env", "namespaceId": "ns1", "enabled": true,
			"tagValues": []any{map[string]any{"name": "prod"}},
		},
		wantState: map[string]any{
			"name": "env", "namespace_id": "ns1",
			"tag_values": []any{map[string]any{"id": "tv1", "name": "prod"}},
		},
	})
}
