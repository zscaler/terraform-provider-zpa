package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

func TestDataSourceUserPortalController(t *testing.T) {
	match := portal_controller.UserPortalController{
		ID: "up1", Name: "Portal01", Enabled: true, ExtDomain: "acme.com", ExtLabel: "portal01", ManagedByZS: true,
	}
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceUserPortalController(),
		listPath:  v1Path("userPortal"),
		itemPath:  v1Path("userPortal", "up1"),
		items:     []portal_controller.UserPortalController{{ID: "up0", Name: "Portal01-old"}, match},
		item:      match,
		name:      "portal01",
		id:        "up1",
		wantState: map[string]any{"name": "Portal01", "ext_domain": "acme.com", "managed_by_zs": true},
	})
}
