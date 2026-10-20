package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/emergencyaccess"
)

// Emergency access users are identified by userId and deactivated, not deleted.
func TestResourceEmergencyAccess_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceEmergencyAccess(),
		config:   map[string]any{"email_id": "breakglass@acme.com", "first_name": "Break", "last_name": "Glass"},
		id:       "u1",
		response: emergencyaccess.EmergencyAccess{
			UserID: "u1", EmailID: "breakglass@acme.com", FirstName: "Break", LastName: "Glass",
		},
		createPath:   v1Path("emergencyAccess", "user"),
		itemPath:     v1Path("emergencyAccess", "user", "u1"),
		deletePath:   v1Path("emergencyAccess", "user", "u1", "deactivate"),
		deleteMethod: "PUT",
		fullFields:   map[string]any{"userId": "u1"},
		updateConfig: func(c map[string]any) { c["last_name"] = "Glass II" },
		wantBody:     map[string]any{"emailId": "breakglass@acme.com", "firstName": "Break"},
		wantState:    map[string]any{"user_id": "u1", "email_id": "breakglass@acme.com"},
	})
}
