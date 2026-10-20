package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_group"
)

func TestResourceTagGroup_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceTagGroup(),
		config: map[string]any{
			"name":        "prod-apps",
			"description": "Production",
			"tags":        []any{"tv1"},
		},
		id: "tg1",
		response: tag_group.TagGroup{
			ID: "tg1", Name: "prod-apps", Description: "Production",
			Tags: []tag_group.Tag{{TagValue: &tag_group.TagValue{ID: "tv1"}}},
		},
		createPath: v1Path("tagGroup"),
		itemPath:   v1Path("tagGroup", "tg1"),
		// tags are tag value IDs, sent as tags[].tagValue.id
		wantBody: map[string]any{
			"name": "prod-apps", "description": "Production",
			"tags": []any{map[string]any{"tagValue": map[string]any{"id": "tv1"}}},
		},
		wantState: map[string]any{"name": "prod-apps", "tags": []any{"tv1"}},
	})
}
