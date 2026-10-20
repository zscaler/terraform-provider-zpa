package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

func customControlConfig() map[string]any {
	return map[string]any{
		"name":           "Block large headers",
		"default_action": "BLOCK",
		"paranoia_level": "1",
		"severity":       "CRITICAL",
		"type":           "REQUEST",
		"rules": []any{map[string]any{
			"type":       "REQUEST_HEADERS",
			"names":      []any{"x-test"},
			"conditions": []any{map[string]any{"lhs": "SIZE", "op": "GE", "rhs": "1000"}},
		}},
	}
}

// The API returns the rules as a JSON string in controlRuleJson.
func customControlResponse() inspection_custom_controls.InspectionCustomControl {
	return inspection_custom_controls.InspectionCustomControl{
		ID: "cc1", Name: "Block large headers", DefaultAction: "BLOCK", ParanoiaLevel: "1",
		Severity: "CRITICAL", Type: "REQUEST",
		ControlRuleJson: `[{"type":"REQUEST_HEADERS","names":["x-test"],"conditions":[{"lhs":"SIZE","op":"GE","rhs":"1000"}]}]`,
	}
}

func TestResourceInspectionCustomControls_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceInspectionCustomControls(),
		config:     customControlConfig(),
		id:         "cc1",
		response:   customControlResponse(),
		createPath: v1Path("inspectionControls", "custom"),
		itemPath:   v1Path("inspectionControls", "custom", "cc1"),
		// the API returns the rules as a JSON string
		fullFields: map[string]any{"controlRuleJson": `[{"type":"REQUEST_HEADERS","names":["x-test"],"conditions":[{"lhs":"SIZE","op":"GE","rhs":"1000"}]}]`},
		wantBody: map[string]any{
			"name": "Block large headers", "severity": "CRITICAL", "type": "REQUEST", "defaultAction": "BLOCK",
			"rules": []any{map[string]any{
				"type": "REQUEST_HEADERS", "names": []any{"x-test"},
				"conditions": []any{map[string]any{"lhs": "SIZE", "op": "GE", "rhs": "1000"}},
			}},
		},
		wantState: map[string]any{
			"rules": []any{map[string]any{
				"type": "REQUEST_HEADERS", "names": []any{"x-test"},
				"conditions": []any{map[string]any{"lhs": "SIZE", "op": "GE", "rhs": "1000"}},
			}},
		},
	})
}

func TestResourceInspectionCustomControls_InvalidRulesRejected(t *testing.T) {
	cases := map[string]struct {
		controlType string
		rule        map[string]any
		wantErr     string
	}{
		"header rule without names": {"REQUEST",
			map[string]any{"type": "REQUEST_HEADERS", "conditions": []any{map[string]any{"lhs": "SIZE", "op": "GE", "rhs": "10"}}},
			"rules.names must be set"},
		"uri rule with names": {"REQUEST",
			map[string]any{"type": "REQUEST_URI", "names": []any{"x"}},
			"rules.name is not allowed"},
		"request rule on a response control": {"RESPONSE",
			map[string]any{"type": "REQUEST_URI"},
			"rules.type must be: RESPONSE_HEADERS || RESPONSE_BODY"},
		"size compared to text": {"REQUEST",
			map[string]any{"type": "REQUEST_BODY", "conditions": []any{map[string]any{"lhs": "SIZE", "op": "GE", "rhs": "big"}}},
			"condition.rhs must be a number"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			config := customControlConfig()
			config["type"] = tc.controlType
			config["rules"] = []any{tc.rule}
			d := fx.resourceData(resourceInspectionCustomControls(), "", config)
			if msg := fx.runErr(resourceInspectionCustomControlsCreate, d); !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want an error containing %q", msg, tc.wantErr)
			}
			if n := len(fx.api.Server.AllRequests()); n != 0 {
				t.Errorf("no API call may be made for invalid rules, got %d", n)
			}
		})
	}
}

// Deleting a custom control first removes it from the inspection profiles that use it.
func TestResourceInspectionCustomControls_DeleteDetachesFromProfiles(t *testing.T) {
	fx := newUnitFixture(t)
	control := customControlResponse()
	control.AssociatedInspectionProfileNames = []zpacommon.AssociatedProfileNames{{ID: "ip1"}}
	fx.on("GET", fx.path("inspectionControls", "custom", "cc1"), common.SuccessResponse(control))
	fx.on("GET", fx.path("inspectionProfile", "ip1"), common.SuccessResponse(inspection_profile.InspectionProfile{
		ID: "ip1", Name: "profile",
		CustomControls: []inspection_profile.InspectionCustomControl{{ID: "cc1"}, {ID: "cc2"}},
	}))
	fx.on("PUT", fx.path("inspectionProfile", "ip1"), common.NoContentResponse())
	fx.on("DELETE", fx.path("inspectionControls", "custom", "cc1"), common.NoContentResponse())

	d := fx.resourceData(resourceInspectionCustomControls(), "cc1", customControlConfig())
	fx.run(resourceInspectionCustomControlsDelete, d)

	assertSubset(t, "profile update", fx.body("PUT", fx.path("inspectionProfile", "ip1")),
		map[string]any{"customControls": []any{map[string]any{"id": "cc2"}}})
	if fx.calls("DELETE", fx.path("inspectionControls", "custom", "cc1")) != 1 {
		t.Error("the custom control itself was not deleted")
	}
}
