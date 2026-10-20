package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_group"
)

func TestDataSourceTagGroup(t *testing.T) {
	match := tag_group.TagGroup{
		ID: "tg1", Name: "Prod-Apps", Description: "Production",
		Tags: []tag_group.Tag{{
			TagKey:   &tag_group.TagKey{ID: "tk1", Name: "env"},
			TagValue: &tag_group.TagValue{ID: "tv1", Name: "prod"},
		}},
	}
	testDataSource(t, dataSourceSpec{
		resource:   dataSourceTagGroup(),
		listPath:   v1Path("tagGroup", "search"),
		listMethod: "POST",
		itemPath:   v1Path("tagGroup", "tg1"),
		items:      []tag_group.TagGroup{{ID: "tg0", Name: "prod-apps-legacy"}, match},
		item:       match,
		name:       "prod-apps",
		id:         "tg1",
		wantState: map[string]any{
			"description": "Production",
			"tags": []any{map[string]any{
				"tag_key":   []any{map[string]any{"id": "tk1", "name": "env"}},
				"tag_value": []any{map[string]any{"id": "tv1", "name": "prod"}},
			}},
		},
	})
}
