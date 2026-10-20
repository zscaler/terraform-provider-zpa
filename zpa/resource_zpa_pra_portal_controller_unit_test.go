package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

func praPortalConfig() map[string]any {
	return map[string]any{
		"name":                      "pra.acme.com",
		"enabled":                   true,
		"domain":                    "pra.acme.com",
		"certificate_id":            "c1",
		"user_notification":         "Created with Terraform",
		"user_notification_enabled": true,
	}
}

func TestResourcePRAPortalController_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourcePRAPortalController(),
		config:   praPortalConfig(),
		id:       "pp1",
		response: praportal.PRAPortal{
			ID: "pp1", Name: "pra.acme.com", Enabled: true, Domain: "pra.acme.com", CertificateID: "c1",
			UserNotification: "Created with Terraform", UserNotificationEnabled: true,
		},
		createPath: v1Path("praPortal"),
		itemPath:   v1Path("praPortal", "pp1"),
		setup: func(fx *unitFixture) {
			fx.on("GET", fx.path("praConsole"), common.SuccessResponse(common.ZPAList([]praconsole.PRAConsole{})))
		},
		wantBody: map[string]any{
			"name": "pra.acme.com", "domain": "pra.acme.com", "certificateId": "c1", "userNotificationEnabled": true,
		},
		wantState: map[string]any{"domain": "pra.acme.com", "certificate_id": "c1"},
	})
}

func TestResourcePRAPortalController_CertificateConflictsWithExternalDomain(t *testing.T) {
	fx := newUnitFixture(t)
	config := praPortalConfig()
	config["ext_label"] = "pra"
	config["ext_domain"] = "acme.com"

	err := fx.diff(resourcePRAPortalController(), "", map[string]any{}, config)
	if err == nil || !strings.Contains(err.Error(), "'certificate_id' cannot be set") {
		t.Fatalf("expected certificate_id with an external domain to be rejected at plan time, got %v", err)
	}
}

// Deleting a portal removes it from the consoles that reference it, and deletes
// a console left without portals.
func TestResourcePRAPortalController_DeleteDetachesFromConsoles(t *testing.T) {
	fx := newUnitFixture(t)
	consoles := []praconsole.PRAConsole{
		{ID: "con1", Name: "shared", PRAPortals: []praconsole.PRAPortals{{ID: "pp1"}, {ID: "pp2"}}},
		{ID: "con2", Name: "only-this-portal", PRAPortals: []praconsole.PRAPortals{{ID: "pp1"}}},
		{ID: "con3", Name: "unrelated", PRAPortals: []praconsole.PRAPortals{{ID: "pp9"}}},
	}
	fx.on("GET", fx.path("praConsole"), common.SuccessResponse(common.ZPAList(consoles)))
	fx.on("PUT", fx.path("praConsole", "con1"), common.NoContentResponse())
	fx.on("DELETE", fx.path("praConsole", "con2"), common.NoContentResponse())
	fx.on("DELETE", fx.path("praPortal", "pp1"), common.NoContentResponse())

	d := fx.resourceData(resourcePRAPortalController(), "pp1", praPortalConfig())
	fx.run(resourcePRAPortalControllerDelete, d)

	assertSubset(t, "console update", fx.body("PUT", fx.path("praConsole", "con1")),
		map[string]any{"praPortals": []any{map[string]any{"id": "pp2"}}})
	if fx.calls("DELETE", fx.path("praConsole", "con2")) != 1 {
		t.Error("a console left without portals must be deleted")
	}
	if fx.calls("PUT", fx.path("praConsole", "con3"))+fx.calls("DELETE", fx.path("praConsole", "con3")) != 0 {
		t.Error("consoles that do not reference the portal must not be changed")
	}
	if fx.calls("DELETE", fx.path("praPortal", "pp1")) != 1 {
		t.Error("the portal itself was not deleted")
	}
}
