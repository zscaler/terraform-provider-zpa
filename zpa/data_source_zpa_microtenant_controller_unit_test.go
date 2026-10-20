package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
)

func TestDataSourceMicrotenantController(t *testing.T) {
	match := microtenants.MicroTenant{
		ID: "mt1", Name: "Engineering", Enabled: true,
		CriteriaAttribute: "AuthDomain", CriteriaAttributeValues: []string{"acme.com"},
		Roles: []microtenants.Roles{{ID: "r1", Name: "Admin", CustomRole: true}},
	}
	testDataSource(t, dataSourceSpec{
		resource: dataSourceMicrotenantController(),
		listPath: v1Path("microtenants"),
		itemPath: v1Path("microtenants", "mt1"),
		items:    []microtenants.MicroTenant{{ID: "mt0", Name: "Engineering-Legacy"}, match},
		item:     match,
		name:     "engineering",
		id:       "mt1",
		wantState: map[string]any{
			"criteria_attribute": "AuthDomain",
			"roles":              []any{map[string]any{"id": "r1", "name": "Admin", "custom_role": true}},
		},
	})
}

// The microtenant admin user is exported in full.
func TestDataSourceMicrotenantController_User(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("microtenants", "mt1"), common.SuccessResponse(microtenants.MicroTenant{
		ID: "mt1", Name: "Engineering",
		UserResource: &microtenants.UserResource{
			ID: "u1", Username: "admin@mt1", DisplayName: "Admin", RoleID: "r1", PhoneNumber: "+15550100",
			PinSession: true, OperationType: "CREATE", ModifiedBy: "72058304855015424", ModifiedTime: "1700000000",
		},
	}))

	d := fx.resourceData(dataSourceMicrotenantController(), "mt1", map[string]any{})
	fx.run(dataSourceMicrotenantControllerRead, d)
	assertState(t, d, map[string]any{"user": []any{map[string]any{
		"username": "admin@mt1", "role_id": "r1", "phone_number": "+15550100", "pin_session": true, "operation_type": "CREATE",
	}}})
}
