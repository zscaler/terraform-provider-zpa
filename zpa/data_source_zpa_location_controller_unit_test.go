package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

// mockExtranet registers the extranet resource "ZIA ER" (er1).
func mockExtranet(fx *unitFixture) {
	fx.on("GET", fx.path("extranetResource", "partner"), common.SuccessResponse(common.ZPAList([]zpacommon.CommonSummary{
		{ID: "er0", Name: "Other ER"}, {ID: "er1", Name: "ZIA ER"},
	})))
}

// ZIA locations are looked up by name within an extranet resource.
func TestDataSourceLocationController(t *testing.T) {
	fx := newUnitFixture(t)
	mockExtranet(fx)
	fx.on("GET", fx.path("location", "extranetResource", "er1"), common.SuccessResponse(common.ZPAList([]zpacommon.CommonSummary{
		{ID: "loc0", Name: "Branch"}, {ID: "loc1", Name: "HQ", Enabled: true},
	})))

	d := fx.resourceData(dataSourceLocationController(), "", map[string]any{"name": "HQ", "zia_er_name": "ZIA ER"})
	fx.run(dataSourceLocationControllerRead, d)
	if d.Id() != "loc1" {
		t.Errorf("id = %q, want loc1", d.Id())
	}
	assertState(t, d, map[string]any{"enabled": true})

	t.Run("LocationNotFound", func(t *testing.T) {
		fx := newUnitFixture(t)
		mockExtranet(fx)
		fx.on("GET", fx.path("location", "extranetResource", "er1"), common.SuccessResponse(common.ZPAList([]zpacommon.CommonSummary{{ID: "loc0", Name: "Branch"}})))
		d := fx.resourceData(dataSourceLocationController(), "", map[string]any{"name": "HQ", "zia_er_name": "ZIA ER"})
		if msg := fx.runErr(dataSourceLocationControllerRead, d); !strings.Contains(msg, "location with name 'HQ' not found") {
			t.Errorf("got %q", msg)
		}
	})

	t.Run("ExtranetNotFound", func(t *testing.T) {
		fx := newUnitFixture(t)
		mockExtranet(fx)
		d := fx.resourceData(dataSourceLocationController(), "", map[string]any{"name": "HQ", "zia_er_name": "Missing ER"})
		if msg := fx.runErr(dataSourceLocationControllerRead, d); !strings.Contains(msg, "error fetching extranet resource") {
			t.Errorf("got %q", msg)
		}
	})
}
