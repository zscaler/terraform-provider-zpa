package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
)

func microtenantConfig() map[string]any {
	return map[string]any{
		"name":                      "Engineering",
		"enabled":                   true,
		"criteria_attribute":        "AuthDomain",
		"criteria_attribute_values": []any{"acme.com"},
	}
}

func TestResourceMicrotenantController_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceMicrotenantController(),
		config:   microtenantConfig(),
		id:       "mt1",
		response: microtenants.MicroTenant{
			ID: "mt1", Name: "Engineering", Enabled: true,
			CriteriaAttribute: "AuthDomain", CriteriaAttributeValues: []string{"acme.com"},
		},
		createPath: v1Path("microtenants"),
		itemPath:   v1Path("microtenants", "mt1"),
		wantBody: map[string]any{
			"name": "Engineering", "enabled": true,
			"criteriaAttribute": "AuthDomain", "criteriaAttributeValues": []any{"acme.com"},
		},
		wantState: map[string]any{"criteria_attribute_values": []any{"acme.com"}},
	})
}

// The microtenant admin user (and its password) is only returned by the create
// call; a later read without it must not clear it from state.
func TestResourceMicrotenantController_KeepsUserFromCreate(t *testing.T) {
	fx := newUnitFixture(t)
	created := microtenants.MicroTenant{
		ID: "mt1", Name: "Engineering",
		UserResource: &microtenants.UserResource{Username: "admin@mt1", Password: "s3cret", DisplayName: "Admin"},
	}
	fx.on("POST", fx.path("microtenants"), common.CreatedResponse(created))
	fx.on("GET", fx.path("microtenants", "mt1"), common.SuccessResponse(microtenants.MicroTenant{ID: "mt1", Name: "Engineering"}))

	d := fx.resourceData(resourceMicrotenantController(), "", microtenantConfig())
	fx.run(resourceMicrotenantCreate, d)
	assertState(t, d, map[string]any{
		"user": []any{map[string]any{"username": "admin@mt1", "password": "s3cret"}},
	})
}
