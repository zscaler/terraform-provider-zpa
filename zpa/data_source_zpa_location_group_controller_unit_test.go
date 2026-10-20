package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

// The location group is found through one of its ZIA locations.
func TestDataSourceLocationGroupController(t *testing.T) {
	fx := newUnitFixture(t)
	mockExtranet(fx)
	fx.on("GET", fx.path("locationGroup", "extranetResource", "er1"), common.SuccessResponse(common.ZPAList([]zpacommon.LocationGroupDTO{
		{ID: "lg0", Name: "Branches", ZiaLocations: []zpacommon.CommonSummary{{ID: "loc0", Name: "Branch"}}},
		{ID: "lg1", Name: "Headquarters", ZiaLocations: []zpacommon.CommonSummary{{ID: "loc1", Name: "HQ", Enabled: true}}},
	})))

	d := fx.resourceData(dataSourceLocationGroupController(), "", map[string]any{"location_name": "HQ", "zia_er_name": "ZIA ER"})
	fx.run(dataSourceLocationGroupControllerRead, d)
	if d.Id() != "lg1" {
		t.Errorf("id = %q, want lg1", d.Id())
	}
	assertState(t, d, map[string]any{
		"location_group_name": "Headquarters",
		"zia_locations":       []any{map[string]any{"id": "loc1", "name": "HQ", "enabled": true}},
	})

	t.Run("NoLocationGroups", func(t *testing.T) {
		fx := newUnitFixture(t)
		mockExtranet(fx)
		fx.on("GET", fx.path("locationGroup", "extranetResource", "er1"), common.SuccessResponse(common.ZPAList([]zpacommon.LocationGroupDTO{})))
		d := fx.resourceData(dataSourceLocationGroupController(), "", map[string]any{"location_name": "HQ", "zia_er_name": "ZIA ER"})
		if msg := fx.runErr(dataSourceLocationGroupControllerRead, d); !strings.Contains(msg, "no location groups found") {
			t.Errorf("got %q", msg)
		}
	})
}
