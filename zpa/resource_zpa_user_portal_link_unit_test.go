package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

func TestResourceUserPortalLink_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceUserPortalLink(),
		config: map[string]any{
			"name":         "server1.acme.com",
			"enabled":      true,
			"link":         "server1.acme.com",
			"protocol":     "https://",
			"icon_text":    "",
			"user_portals": []any{map[string]any{"id": []any{"up1"}}},
		},
		id: "ul1",
		response: portal_link.UserPortalLink{
			ID: "ul1", Name: "server1.acme.com", Enabled: true, Link: "server1.acme.com", Protocol: "https://",
			UserPortals: []portal_controller.UserPortalController{{ID: "up1"}},
		},
		createPath: v1Path("userPortalLink"),
		itemPath:   v1Path("userPortalLink", "ul1"),
		// user_portals { id = [...] } is sent as userPortals[].id
		wantBody: map[string]any{
			"name": "server1.acme.com", "link": "server1.acme.com", "protocol": "https://",
			"userPortals": []any{map[string]any{"id": "up1"}},
		},
		wantState: map[string]any{
			"link":         "server1.acme.com",
			"user_portals": []any{map[string]any{"id": []any{"up1"}}},
		},
	})
}
