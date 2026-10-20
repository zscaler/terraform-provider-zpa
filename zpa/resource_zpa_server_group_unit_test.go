package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

func serverGroupConfig() map[string]any {
	return map[string]any{
		"name":                 "CRM servers",
		"description":          "Created with Terraform",
		"enabled":              true,
		"dynamic_discovery":    false,
		"servers":              []any{map[string]any{"id": []any{"as1"}}},
		"app_connector_groups": []any{map[string]any{"id": []any{"acg1"}}},
	}
}

func serverGroupResponse() servergroup.ServerGroup {
	return servergroup.ServerGroup{
		ID: "sg1", Name: "CRM servers", Description: "Created with Terraform", Enabled: true,
		Servers:            []appservercontroller.ApplicationServer{{ID: "as1"}},
		AppConnectorGroups: []appconnectorgroup.AppConnectorGroup{{ID: "acg1"}},
	}
}

func TestResourceServerGroup_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceServerGroup(),
		config:     serverGroupConfig(),
		id:         "sg1",
		response:   serverGroupResponse(),
		createPath: v1Path("serverGroup"),
		itemPath:   v1Path("serverGroup", "sg1"),
		listPath:   v1Path("serverGroup"),
		setup: func(fx *unitFixture) {
			fx.mockPolicySet("ACCESS_POLICY")
			fx.on("GET", fx.path("application"), common.SuccessResponse(common.ZPAList([]any{})))
		},
		updateConfig: func(c map[string]any) { c["description"] = "Updated" },
		wantBody: map[string]any{
			"dynamicDiscovery":   false,
			"servers":            []any{map[string]any{"id": "as1"}},
			"appConnectorGroups": []any{map[string]any{"id": "acg1"}},
		},
		wantState: map[string]any{
			"servers":              []any{map[string]any{"id": []any{"as1"}}},
			"app_connector_groups": []any{map[string]any{"id": []any{"acg1"}}},
		},
	})
}

func TestResourceServerGroup_ServersAndDynamicDiscoveryRules(t *testing.T) {
	cases := map[string]struct {
		dynamicDiscovery bool
		servers          []any
		wantErr          string
	}{
		"servers with dynamic discovery":       {true, []any{map[string]any{"id": []any{"as1"}}}, "only be attached to a server when DynamicDiscovery is disabled"},
		"no servers without dynamic discovery": {false, nil, "servers must not be empty when DynamicDiscovery is disabled"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fx := newUnitFixture(t)
			config := serverGroupConfig()
			config["dynamic_discovery"] = tc.dynamicDiscovery
			if tc.servers == nil {
				delete(config, "servers")
			} else {
				config["servers"] = tc.servers
			}
			d := fx.resourceData(resourceServerGroup(), "", config)
			if msg := fx.runErr(resourceServerGroupCreate, d); !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("got %q, want an error containing %q", msg, tc.wantErr)
			}
			if n := fx.calls("POST", fx.path("serverGroup")); n != 0 {
				t.Errorf("no create request may be sent, got %d", n)
			}
		})
	}
}

// Deleting a server group removes it from the application segments and access
// rules that reference it, without touching the ones that don't.
func TestResourceServerGroup_DeleteDetachesReferences(t *testing.T) {
	fx := newUnitFixture(t)
	ps := fx.mockPolicySet("ACCESS_POLICY")
	fx.on("GET", fx.path("serverGroup", "sg1"), common.SuccessResponse(servergroup.ServerGroup{ID: "sg1"}))
	fx.on("GET", fx.path("policySet", "rules", "policyType", "ACCESS_POLICY"), common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{
		{ID: "r1", AppServerGroups: []servergroup.ServerGroup{{ID: "sg1"}, {ID: "sg2"}}},
	})))
	fx.on("PUT", fx.path("policySet", ps, "rule", "r1"), common.NoContentResponse())
	fx.on("GET", fx.path("application"), common.SuccessResponse(common.ZPAList([]applicationsegment.ApplicationSegmentResource{
		{ID: "app1", Name: "uses-it", ServerGroups: []servergroup.ServerGroup{{ID: "sg1"}, {ID: "sg2"}}},
		{ID: "app2", Name: "unrelated", ServerGroups: []servergroup.ServerGroup{{ID: "sg2"}}},
	})))
	fx.on("PUT", fx.path("application", "app1"), common.NoContentResponse())
	fx.on("PUT", fx.path("application", "app2"), common.NoContentResponse())
	fx.on("DELETE", fx.path("serverGroup", "sg1"), common.NoContentResponse())

	d := fx.resourceData(resourceServerGroup(), "sg1", serverGroupConfig())
	fx.run(resourceServerGroupDelete, d)

	assertSubset(t, "access rule update", fx.body("PUT", fx.path("policySet", ps, "rule", "r1")),
		map[string]any{"appServerGroups": []any{map[string]any{"id": "sg2"}}})
	assertSubset(t, "application segment update", fx.body("PUT", fx.path("application", "app1")),
		map[string]any{"serverGroups": []any{map[string]any{"id": "sg2"}}})
	if n := fx.calls("PUT", fx.path("application", "app2")); n != 0 {
		t.Errorf("an application segment that does not use the server group must not be updated, got %d", n)
	}
	if fx.calls("DELETE", fx.path("serverGroup", "sg1")) != 1 {
		t.Error("the server group itself was not deleted")
	}
}
