package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

var weightedLBPath = v1Path("application", "2001", "weightedLbConfig")

func weightedLBConfig() map[string]any {
	return map[string]any{
		"application_id":          "2001",
		"weighted_load_balancing": true,
		"application_to_server_group_mappings": []any{
			map[string]any{"id": "sg1", "weight": "10"},
			map[string]any{"name": "Backup servers", "passive": true},
		},
	}
}

func weightedLBResponse() applicationsegment.WeightedLoadBalancerConfig {
	return applicationsegment.WeightedLoadBalancerConfig{
		ApplicationID: "2001", WeightedLoadBalancing: true,
		ApplicationToServerGroupMaps: []applicationsegment.ApplicationToServerGroupMapping{
			{ID: "sg1", Name: "Primary servers", Weight: "10"},
			{ID: "sg2", Name: "Backup servers", Passive: true, Weight: "0"},
		},
	}
}

// mockWeightedLB registers the config, the application and the server group lookup.
func mockWeightedLB(fx *unitFixture) {
	fx.on("GET", weightedLBPath, common.SuccessResponse(weightedLBResponse()))
	fx.on("PUT", weightedLBPath, common.NoContentResponse())
	fx.on("GET", fx.path("application", "2001"), common.SuccessResponse(applicationsegment.ApplicationSegmentResource{ID: "2001", Name: "crm"}))
	fx.on("GET", fx.path("application"), common.SuccessResponse(common.ZPAList([]applicationsegment.ApplicationSegmentResource{{ID: "2001", Name: "crm"}})))
	fx.on("GET", fx.path("serverGroup"), common.SuccessResponse(common.ZPAList([]servergroup.ServerGroup{{ID: "sg2", Name: "Backup servers"}})))
}

// Server groups given by name are resolved to IDs; a missing weight is sent as "0".
func TestResourceApplicationSegmentWeightedLBConfig_Create(t *testing.T) {
	fx := newUnitFixture(t)
	mockWeightedLB(fx)

	d := fx.resourceData(resourceApplicationSegmentWeightedLBConfig(), "", weightedLBConfig())
	fx.run(resourceApplicationSegmentWeightedLBConfigCreate, d)
	if d.Id() != "2001" {
		t.Errorf("id = %q, want the application ID 2001", d.Id())
	}
	assertSubset(t, "config request", fx.body("PUT", weightedLBPath), map[string]any{
		"applicationId": "2001", "weightedLoadBalancing": true,
		"applicationToServerGroupMappings": []any{
			map[string]any{"id": "sg1", "weight": "10", "passive": false},
			map[string]any{"id": "sg2", "name": "Backup servers", "weight": "0", "passive": true},
		},
	})
	assertState(t, d, map[string]any{"application_name": "crm"})
}

// The application can be given by name instead of ID.
func TestResourceApplicationSegmentWeightedLBConfig_CreateByApplicationName(t *testing.T) {
	fx := newUnitFixture(t)
	mockWeightedLB(fx)
	config := weightedLBConfig()
	delete(config, "application_id")
	config["application_name"] = "crm"

	d := fx.resourceData(resourceApplicationSegmentWeightedLBConfig(), "", config)
	fx.run(resourceApplicationSegmentWeightedLBConfigCreate, d)
	if d.Id() != "2001" {
		t.Errorf("id = %q, want 2001", d.Id())
	}
}

func TestResourceApplicationSegmentWeightedLBConfig_RequiresApplication(t *testing.T) {
	fx := newUnitFixture(t)
	config := weightedLBConfig()
	delete(config, "application_id")
	d := fx.resourceData(resourceApplicationSegmentWeightedLBConfig(), "", config)
	if msg := fx.runErr(resourceApplicationSegmentWeightedLBConfigCreate, d); !strings.Contains(msg, "Missing application identifier") {
		t.Fatalf("got %q", msg)
	}
}

func TestResourceApplicationSegmentWeightedLBConfig_MappingRequiresServerGroup(t *testing.T) {
	fx := newUnitFixture(t)
	config := weightedLBConfig()
	config["application_to_server_group_mappings"] = []any{map[string]any{"weight": "10"}}
	d := fx.resourceData(resourceApplicationSegmentWeightedLBConfig(), "", config)
	if msg := fx.runErr(resourceApplicationSegmentWeightedLBConfigCreate, d); !strings.Contains(msg, "Missing server group identifier") {
		t.Fatalf("got %q", msg)
	}
	if n := fx.calls("PUT", weightedLBPath); n != 0 {
		t.Errorf("no config request may be sent, got %d", n)
	}
}

func TestResourceApplicationSegmentWeightedLBConfig_Update(t *testing.T) {
	fx := newUnitFixture(t)
	mockWeightedLB(fx)
	config := weightedLBConfig()
	config["weighted_load_balancing"] = false
	d := fx.updateData(resourceApplicationSegmentWeightedLBConfig(), "2001", weightedLBConfig(), config)
	fx.run(resourceApplicationSegmentWeightedLBConfigUpdate, d)
	assertSubset(t, "config request", fx.body("PUT", weightedLBPath), map[string]any{"weightedLoadBalancing": false})
}

// Removing the resource turns weighted load balancing off; there is nothing to delete.
func TestResourceApplicationSegmentWeightedLBConfig_DeleteDisables(t *testing.T) {
	fx := newUnitFixture(t)
	mockWeightedLB(fx)
	d := fx.resourceData(resourceApplicationSegmentWeightedLBConfig(), "2001", weightedLBConfig())
	fx.run(resourceApplicationSegmentWeightedLBConfigDelete, d)
	assertSubset(t, "config request", fx.body("PUT", weightedLBPath), map[string]any{"applicationId": "2001", "weightedLoadBalancing": false})
}

// A deleted application segment takes its weighted LB config with it.
func TestResourceApplicationSegmentWeightedLBConfig_ApplicationGone(t *testing.T) {
	for name, step := range map[string]crudFunc{
		"Read":   resourceApplicationSegmentWeightedLBConfigRead,
		"Delete": resourceApplicationSegmentWeightedLBConfigDelete,
	} {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			fx.on("GET", weightedLBPath, common.NotFoundResponse())
			fx.on("PUT", weightedLBPath, common.NotFoundResponse())
			d := fx.resourceData(resourceApplicationSegmentWeightedLBConfig(), "2001", weightedLBConfig())
			fx.run(step, d)
			if d.Id() != "" {
				t.Errorf("id = %q, want it removed from state", d.Id())
			}
		})
	}
}

func TestResourceApplicationSegmentWeightedLBConfig_APIErrors(t *testing.T) {
	for name, step := range map[string]crudFunc{
		"Read":   resourceApplicationSegmentWeightedLBConfigRead,
		"Update": resourceApplicationSegmentWeightedLBConfigUpdate,
		"Delete": resourceApplicationSegmentWeightedLBConfigDelete,
	} {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			mockWeightedLB(fx)
			fx.on("GET", weightedLBPath, apiError())
			fx.on("PUT", weightedLBPath, apiError())
			d := fx.resourceData(resourceApplicationSegmentWeightedLBConfig(), "2001", weightedLBConfig())
			if msg := fx.runErr(step, d); !strings.Contains(msg, apiErrorMarker) {
				t.Errorf("expected the API error to be returned, got %q", msg)
			}
		})
	}
}
