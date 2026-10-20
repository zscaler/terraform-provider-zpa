package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

func TestResourceUserPortalController_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceUserPortalController(),
		config: map[string]any{
			"name":                      "Portal01",
			"enabled":                   true,
			"ext_domain":                "acme.com",
			"ext_label":                 "portal01",
			"user_notification":         "Welcome",
			"user_notification_enabled": true,
		},
		id: "up1",
		response: portal_controller.UserPortalController{
			ID: "up1", Name: "Portal01", Enabled: true, ExtDomain: "acme.com", ExtLabel: "portal01",
			UserNotification: "Welcome", UserNotificationEnabled: true,
		},
		createPath: v1Path("userPortal"),
		itemPath:   v1Path("userPortal", "up1"),
		wantBody: map[string]any{
			"name": "Portal01", "enabled": true, "extDomain": "acme.com", "extLabel": "portal01",
			"userNotification": "Welcome", "userNotificationEnabled": true,
		},
		wantState: map[string]any{"name": "Portal01", "ext_domain": "acme.com", "user_notification_enabled": true},
	})
}
