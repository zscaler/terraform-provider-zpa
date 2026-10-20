package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

func TestDataSourceInspectionProfile(t *testing.T) {
	match := inspection_profile.InspectionProfile{ID: "ip1", Name: "AppProtection profile", ParanoiaLevel: "1", GlobalControlActions: []string{"PREDEFINED:PASS"}}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceInspectionProfile(),
		listPath:  v1Path("inspectionProfile"),
		itemPath:  v1Path("inspectionProfile", "ip1"),
		items:     []inspection_profile.InspectionProfile{{ID: "other", Name: "AppProtection profile-old"}, match},
		item:      match,
		name:      "appprotection profile",
		id:        "ip1",
		wantState: map[string]any{"paranoia_level": "1", "global_control_actions": []any{"PREDEFINED:PASS"}},
	})
}
