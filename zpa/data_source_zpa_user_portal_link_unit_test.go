package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

func TestDataSourceUserPortalLink(t *testing.T) {
	match := portal_link.UserPortalLink{
		ID: "ul1", Name: "server1.acme.com", Enabled: true, Link: "server1.acme.com", Protocol: "https://",
		UserPortals: []portal_controller.UserPortalController{{ID: "up1", Name: "Portal01", Enabled: true}},
	}
	testDataSource(t, dataSourceSpec{
		resource: dataSourceUserPortalLink(),
		listPath: v1Path("userPortalLink"),
		itemPath: v1Path("userPortalLink", "ul1"),
		items:    []portal_link.UserPortalLink{{ID: "ul0", Name: "server1.acme.com.old"}, match},
		item:     match,
		name:     "SERVER1.acme.com",
		id:       "ul1",
		wantState: map[string]any{
			"link":         "server1.acme.com",
			"user_portals": []any{map[string]any{"id": "up1", "name": "Portal01", "enabled": true}},
		},
	})
}
