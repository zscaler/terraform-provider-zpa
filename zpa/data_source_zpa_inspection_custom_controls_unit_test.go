package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
)

func TestDataSourceInspectionCustomControls(t *testing.T) {
	match := inspection_custom_controls.InspectionCustomControl{ID: "cc1", Name: "Block large headers", Severity: "CRITICAL", Type: "REQUEST", ControlRuleJson: `[{"type":"REQUEST_HEADERS","names":["x-test"],"conditions":[{"lhs":"SIZE","op":"GE","rhs":"1000"}]}]`}
	testDataSource(t, dataSourceSpec{
		resource: dataSourceInspectionCustomControls(),
		listPath: v1Path("inspectionControls", "custom"),
		itemPath: v1Path("inspectionControls", "custom", "cc1"),
		items:    []inspection_custom_controls.InspectionCustomControl{{ID: "other", Name: "Block large headers-old"}, match},
		item:     match,
		name:     "block large headers",
		// the API returns the rules as a JSON string
		fullFields: map[string]any{"controlRuleJson": `[{"type":"REQUEST_HEADERS","names":["x-test"],"conditions":[{"lhs":"SIZE","op":"GE","rhs":"1000"}]}]`},
		id:         "cc1",
		wantState:  map[string]any{"severity": "CRITICAL", "rules": []any{map[string]any{"type": "REQUEST_HEADERS", "names": []any{"x-test"}, "conditions": []any{map[string]any{"lhs": "SIZE", "op": "GE", "rhs": "1000"}}}}},
	})
}
