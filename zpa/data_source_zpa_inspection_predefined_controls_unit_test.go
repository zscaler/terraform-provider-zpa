package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_predefined_controls"
)

var sqlInjectionControl = inspection_predefined_controls.PredefinedControls{
	ID: "pc1", Name: "Failed to parse request body", ControlGroup: "Protocol Issues", Severity: "CRITICAL", Version: "OWASP_CRS/3.3.0",
}

// The name lookup returns controls grouped by control group, for one version.
func TestDataSourceInspectionPredefinedControls_ByName(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("inspectionControls", "predefined"), common.SuccessResponse([]inspection_predefined_controls.ControlGroupItem{{
		ControlGroup:                 "Protocol Issues",
		PredefinedInspectionControls: []inspection_predefined_controls.PredefinedControls{{ID: "pc0", Name: "Other"}, sqlInjectionControl},
	}}))

	d := fx.resourceData(dataSourceInspectionPredefinedControls(), "", map[string]any{"name": "failed to parse request body", "version": "OWASP_CRS/3.3.0"})
	fx.run(dataSourceInspectionPredefinedControlsRead, d)
	if d.Id() != "pc1" {
		t.Errorf("id = %q, want pc1", d.Id())
	}
	assertState(t, d, map[string]any{"control_group": "Protocol Issues", "severity": "CRITICAL"})
}

func TestDataSourceInspectionPredefinedControls_ByID(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("inspectionControls", "predefined", "pc1"), common.SuccessResponse(sqlInjectionControl))

	d := fx.resourceData(dataSourceInspectionPredefinedControls(), "pc1", map[string]any{})
	fx.run(dataSourceInspectionPredefinedControlsRead, d)
	assertState(t, d, map[string]any{"name": "Failed to parse request body"})
}

func TestDataSourceInspectionPredefinedControls_NameRequiresVersion(t *testing.T) {
	fx := newUnitFixture(t)
	d := fx.resourceData(dataSourceInspectionPredefinedControls(), "", map[string]any{"name": "Failed to parse request body"})
	if msg := fx.runErr(dataSourceInspectionPredefinedControlsRead, d); !strings.Contains(msg, "version must be set") {
		t.Fatalf("got %q, want a version required error", msg)
	}
}
