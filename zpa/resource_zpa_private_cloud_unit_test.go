package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud"
)

func privateCloudConfig() map[string]any {
	return map[string]any{
		"name":                 "SJC Private Cloud",
		"description":          "Created with Terraform",
		"enabled":              true,
		"re_enroll_period":     "30",
		"assistant_groups_ids": []any{map[string]any{"id": []any{"acg1"}}},
	}
}

func privateCloudResponse() private_cloud.PrivateCloudController {
	return private_cloud.PrivateCloudController{
		ID: "pc1", Name: "SJC Private Cloud", Description: "Created with Terraform", Enabled: true, ReEnrollPeriod: "30",
		AssistantGroupsIDs: []zpacommon.CommonSummary{{ID: "acg1"}},
	}
}

func TestResourcePrivateCloud_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourcePrivateCloud(),
		config:     privateCloudConfig(),
		id:         "pc1",
		response:   privateCloudResponse(),
		createPath: v1Path("privateCloud"),
		itemPath:   v1Path("privateCloud", "pc1"),
		// associated groups are sent as ID-only objects
		wantBody: map[string]any{
			"name": "SJC Private Cloud", "reEnrollPeriod": "30",
			"assistantGroupsIds": []any{map[string]any{"id": "acg1"}},
		},
		wantState: map[string]any{"assistant_groups_ids": []any{map[string]any{"id": []any{"acg1"}}}},
	})
}

// The API rejects fire_drill_enabled on create: the private cloud is created
// without it and fire drill is enabled by a follow-up update.
func TestResourcePrivateCloud_FireDrillEnabledAfterCreate(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("POST", fx.path("privateCloud"), common.CreatedResponse(privateCloudResponse()))
	fx.on("PUT", fx.path("privateCloud", "pc1"), common.NoContentResponse())
	fx.on("GET", fx.path("privateCloud", "pc1"), common.SuccessResponse(privateCloudResponse()))

	config := privateCloudConfig()
	config["fire_drill_enabled"] = true
	config["zpn_fire_drill_site"] = []any{map[string]any{"fire_drill_interval": "2", "fire_drill_interval_time_unit": "HOURS"}}
	d := fx.resourceData(resourcePrivateCloud(), "", config)
	fx.run(resourcePrivateCloudCreate, d)

	created := fx.body("POST", fx.path("privateCloud"))
	if created["fireDrillEnabled"] == true || created["zpnFireDrillSite"] != nil {
		t.Errorf("create must not send fire drill settings: %v", created)
	}
	assertSubset(t, "follow-up update", fx.body("PUT", fx.path("privateCloud", "pc1")), map[string]any{
		"fireDrillEnabled": true,
		"zpnFireDrillSite": map[string]any{"fireDrillInterval": "2", "fireDrillIntervalTimeUnit": "HOURS"},
	})
}

func TestResourcePrivateCloud_FireDrillSiteRequiresFireDrill(t *testing.T) {
	fx := newUnitFixture(t)
	config := privateCloudConfig()
	config["zpn_fire_drill_site"] = []any{map[string]any{"fire_drill_interval": "2", "fire_drill_interval_time_unit": "HOURS"}}
	d := fx.resourceData(resourcePrivateCloud(), "", config)
	if msg := fx.runErr(resourcePrivateCloudCreate, d); !strings.Contains(msg, "requires fire_drill_enabled") {
		t.Fatalf("got %q, want a fire_drill_enabled error", msg)
	}
	if n := len(fx.api.Server.AllRequests()); n != 0 {
		t.Errorf("no API call may be made, got %d", n)
	}
}
