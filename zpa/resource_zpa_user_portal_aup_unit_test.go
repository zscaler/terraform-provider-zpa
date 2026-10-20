package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
)

func TestResourceUserPortalAUP_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceUserPortalAUP(),
		config: map[string]any{
			"name": "Corp AUP", "enabled": true, "aup": "Use responsibly.", "email": "it@acme.com", "phone_num": "+15550100",
		},
		id: "aup1",
		response: aup.UserPortalAup{
			ID: "aup1", Name: "Corp AUP", Enabled: true, Aup: "Use responsibly.", Email: "it@acme.com", PhoneNum: "+15550100",
		},
		createPath: v1Path("userportal", "aup"),
		itemPath:   v1Path("userportal", "aup", "aup1"),
		wantBody:   map[string]any{"name": "Corp AUP", "aup": "Use responsibly.", "email": "it@acme.com", "phoneNum": "+15550100"},
		wantState:  map[string]any{"aup": "Use responsibly.", "phone_num": "+15550100"},
	})
}
