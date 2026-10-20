package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

func inspectionProfileConfig() map[string]any {
	return map[string]any{
		"name":                   "AppProtection profile",
		"paranoia_level":         "1",
		"global_control_actions": []any{"PREDEFINED:PASS", "CUSTOM:NONE", "OVERRIDE_ACTION:COMMON"},
		"predefined_controls":    []any{map[string]any{"id": "pc1", "action": "BLOCK"}},
		"custom_controls":        []any{map[string]any{"id": "cc1", "action": "PASS"}},
	}
}

func inspectionProfileResponse() inspection_profile.InspectionProfile {
	return inspection_profile.InspectionProfile{
		ID: "ip1", Name: "AppProtection profile", ParanoiaLevel: "1", PredefinedControlsVersion: "OWASP_CRS/3.3.0",
		GlobalControlActions: []string{"PREDEFINED:PASS", "CUSTOM:NONE", "OVERRIDE_ACTION:COMMON"},
		PredefinedControls:   []zpacommon.CustomCommonControls{{ID: "pc1", Action: "BLOCK"}},
		CustomControls:       []inspection_profile.InspectionCustomControl{{ID: "cc1", Action: "PASS"}},
	}
}

func TestResourceInspectionProfile_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceInspectionProfile(),
		config:     inspectionProfileConfig(),
		id:         "ip1",
		response:   inspectionProfileResponse(),
		createPath: v1Path("inspectionProfile"),
		itemPath:   v1Path("inspectionProfile", "ip1"),
		// predefined_controls_version defaults to OWASP_CRS/3.3.0
		wantBody: map[string]any{
			"name": "AppProtection profile", "paranoiaLevel": "1", "predefinedControlsVersion": "OWASP_CRS/3.3.0",
			"predefinedControls": []any{map[string]any{"id": "pc1", "action": "BLOCK"}},
			"customControls":     []any{map[string]any{"id": "cc1", "action": "PASS"}},
		},
		wantState: map[string]any{
			"predefined_controls_version": "OWASP_CRS/3.3.0",
			"global_control_actions":      []any{"PREDEFINED:PASS", "CUSTOM:NONE", "OVERRIDE_ACTION:COMMON"},
		},
	})
}

func TestResourceInspectionProfile_RedirectWithoutValueRejected(t *testing.T) {
	fx := newUnitFixture(t)
	config := inspectionProfileConfig()
	config["predefined_controls"] = []any{map[string]any{"id": "pc1", "action": "REDIRECT"}}

	d := fx.resourceData(resourceInspectionProfile(), "", config)
	if msg := fx.runErr(resourceInspectionProfileCreate, d); !strings.Contains(msg, "action value must be set") {
		t.Fatalf("expected REDIRECT without action_value to be rejected, got %q", msg)
	}
	if n := len(fx.api.Server.AllRequests()); n != 0 {
		t.Errorf("no API call may be made, got %d", n)
	}
}

// associate_all_controls attaches every predefined control after the profile is created.
func TestResourceInspectionProfile_AssociateAllControls(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("POST", fx.path("inspectionProfile"), common.CreatedResponse(inspectionProfileResponse()))
	fx.on("GET", fx.path("inspectionProfile", "ip1"), common.SuccessResponse(inspectionProfileResponse()))
	fx.on("PUT", fx.path("inspectionProfile", "ip1", "associateAllPredefinedControls"), common.NoContentResponse())

	config := inspectionProfileConfig()
	config["associate_all_controls"] = true
	d := fx.resourceData(resourceInspectionProfile(), "", config)
	fx.run(resourceInspectionProfileCreate, d)

	if n := fx.calls("PUT", fx.path("inspectionProfile", "ip1", "associateAllPredefinedControls")); n != 1 {
		t.Errorf("expected 1 associateAllPredefinedControls request, got %d", n)
	}
}

// The API reports how many controls of each type the profile has (controlsInfo);
// reading it back must not fail.
func TestResourceInspectionProfile_ReadWithControlsInfo(t *testing.T) {
	fx := newUnitFixture(t)
	resp := inspectionProfileResponse()
	resp.ControlInfoResource = []inspection_profile.ControlInfoResource{{ControlType: "CUSTOM", Count: "1"}}
	fx.on("GET", fx.path("inspectionProfile", "ip1"), common.SuccessResponse(resp))

	d := fx.resourceData(resourceInspectionProfile(), "ip1", inspectionProfileConfig())
	fx.run(resourceInspectionProfileRead, d)
	assertState(t, d, map[string]any{"controls_info": []any{map[string]any{"control_type": "CUSTOM"}}})
}
