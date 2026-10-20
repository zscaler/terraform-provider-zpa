package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_predefined_controls"
)

var predefinedControlGroups = []inspection_predefined_controls.ControlGroupItem{
	{ControlGroup: "Protocol Issues", PredefinedInspectionControls: []inspection_predefined_controls.PredefinedControls{{ID: "pc1", Name: "a"}, {ID: "pc2", Name: "b"}}},
	{ControlGroup: "Preprocessors", PredefinedInspectionControls: []inspection_predefined_controls.PredefinedControls{{ID: "pc3", Name: "c"}}},
}

func TestDataSourceInspectionAllPredefinedControls(t *testing.T) {
	cases := map[string]struct {
		config  map[string]any
		wantIDs []string
	}{
		"all groups":   {map[string]any{"version": "OWASP_CRS/3.3.0"}, []string{"pc1", "pc2", "pc3"}},
		"single group": {map[string]any{"version": "OWASP_CRS/3.3.0", "group_name": "Preprocessors"}, []string{"pc3"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			fx.on("GET", fx.path("inspectionControls", "predefined"), common.SuccessResponse(predefinedControlGroups))

			d := fx.resourceData(dataSourceInspectionAllPredefinedControls(), "", tc.config)
			fx.run(dataSourceInspectionAllPredefinedControlsRead, d)
			var got []string
			for _, c := range d.Get("list").([]any) {
				got = append(got, c.(map[string]any)["id"].(string))
			}
			if strings.Join(got, ",") != strings.Join(tc.wantIDs, ",") {
				t.Errorf("controls = %v, want %v", got, tc.wantIDs)
			}
		})
	}
}

func TestDataSourceInspectionAllPredefinedControls_RequiresVersion(t *testing.T) {
	fx := newUnitFixture(t)
	d := fx.resourceData(dataSourceInspectionAllPredefinedControls(), "", map[string]any{})
	if msg := fx.runErr(dataSourceInspectionAllPredefinedControlsRead, d); msg == "" {
		t.Fatal("expected an error without version")
	}
}
