package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
)

func TestResourceApplicationServer_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceApplicationServer(),
		config: map[string]any{
			"name":        "server1",
			"description": "Created with Terraform",
			"address":     "server1.acme.com",
			"enabled":     true,
		},
		id: "as1",
		response: appservercontroller.ApplicationServer{
			ID: "as1", Name: "server1", Description: "Created with Terraform", Address: "server1.acme.com", Enabled: true,
		},
		createPath: v1Path("server"),
		itemPath:   v1Path("server", "as1"),
		// refresh reads the whole server list once and serves every server from it
		listPath:  v1Path("server"),
		wantBody:  map[string]any{"name": "server1", "address": "server1.acme.com", "enabled": true},
		wantState: map[string]any{"address": "server1.acme.com"},
	})
}

// A refresh reads servers from the list endpoint, not one GET per server.
func TestResourceApplicationServer_ReadFromList(t *testing.T) {
	fx := newUnitFixture(t)
	servers := []appservercontroller.ApplicationServer{
		{ID: "as1", Name: "server1", Address: "server1.acme.com", AppServerGroupIds: []string{"sg1"}},
		{ID: "as2", Name: "server2", Address: "server2.acme.com"},
	}
	fx.on("GET", fx.path("server"), common.SuccessResponse(common.ZPAList(servers)))

	for _, s := range servers {
		d := fx.resourceData(resourceApplicationServer(), s.ID, map[string]any{"name": s.Name, "address": "x"})
		fx.run(resourceApplicationServerRead, d)
		assertState(t, d, map[string]any{"address": s.Address})
	}
	if n := fx.calls("GET", fx.path("server")); n != 1 {
		t.Errorf("server list fetched %d times, want 1", n)
	}
	if n := fx.calls("GET", fx.path("server", "as1")) + fx.calls("GET", fx.path("server", "as2")); n != 0 {
		t.Errorf("expected no per-server GET, got %d", n)
	}
}
