package zpa

import (
	"testing"

	tag_key_controller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_key"
)

func TestDataSourceTagKey(t *testing.T) {
	match := tag_key_controller.TagKey{
		ID: "tk1", Name: "Env", NamespaceID: "ns1", Enabled: true,
		TagValues: []tag_key_controller.TagValue{{ID: "tv1", Name: "prod"}},
	}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceTagKey(),
		listPath:   v1Path("namespace", "ns1", "tagKey", "search"),
		listMethod: "POST",
		itemPath:   v1Path("namespace", "ns1", "tagKey", "tk1"),
		items:      []tag_key_controller.TagKey{{ID: "tk0", Name: "environment", NamespaceID: "ns1"}, match},
		item:       match,
		name:       "env",
		id:         "tk1",
		// tag keys are looked up within a namespace
		extra: map[string]any{"namespace_id": "ns1"},
		wantState: map[string]any{
			"namespace_id": "ns1", "enabled": true,
			"tag_values": []any{map[string]any{"id": "tv1", "name": "prod"}},
		},
	})
}
