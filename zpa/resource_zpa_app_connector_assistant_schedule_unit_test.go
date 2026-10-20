package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func TestResourceAppConnectorAssistantSchedule(t *testing.T) {
	testAssistantSchedule(t, resourceAppConnectorAssistantSchedule(), "connectorSchedule")
}

// testAssistantSchedule covers the app connector and service edge cleanup
// schedules: one schedule per tenant, created once and updated afterwards.
func testAssistantSchedule(t *testing.T, r *schema.Resource, endpoint string) {
	config := map[string]any{
		"customer_id":        common.TestCustomerID,
		"enabled":            true,
		"delete_disabled":    true,
		"frequency":          "days",
		"frequency_interval": "5",
	}
	schedule := map[string]any{
		"id": "s1", "customerId": common.TestCustomerID, "enabled": true,
		"deleteDisabled": true, "frequency": "days", "frequencyInterval": "5",
	}

	t.Run("Create", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("POST", fx.path(endpoint), common.CreatedResponse(schedule))
		fx.on("GET", fx.path(endpoint), common.SuccessResponse(schedule))
		d := fx.resourceData(r, "", config)
		fx.run(r.CreateContext, d)
		// the ID comes from the create response
		if d.Id() != "s1" {
			t.Errorf("id = %q, want s1", d.Id())
		}
		assertSubset(t, "create request", fx.body("POST", fx.path(endpoint)), map[string]any{
			"customerId": common.TestCustomerID, "enabled": true, "frequency": "days", "frequencyInterval": "5",
		})
	})

	// The tenant already has a schedule: create adopts and updates it.
	t.Run("CreateWhenExistsUpdates", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("POST", fx.path(endpoint), common.SuccessResponseWithStatus(400,
			map[string]any{"id": "resource.already.exist", "reason": "Schedule already exists"}))
		fx.on("GET", fx.path(endpoint), common.SuccessResponse(schedule))
		fx.on("PUT", fx.path(endpoint, "s1"), common.NoContentResponse())
		d := fx.resourceData(r, "", config)
		fx.run(r.CreateContext, d)
		if d.Id() != "s1" {
			t.Errorf("id = %q, want s1", d.Id())
		}
		if fx.calls("PUT", fx.path(endpoint, "s1")) != 1 {
			t.Error("the existing schedule was not updated")
		}
	})

	t.Run("ReadNotFoundRemovesFromState", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.on("GET", fx.path(endpoint), common.NotFoundResponse())
		d := fx.resourceData(r, "s1", config)
		fx.run(r.ReadContext, d)
		if d.Id() != "" {
			t.Errorf("id = %q, want it removed from state", d.Id())
		}
	})

	// The schedule cannot be deleted: destroy only removes it from state.
	t.Run("DeleteMakesNoAPICall", func(t *testing.T) {
		fx := newUnitFixture(t)
		d := fx.resourceData(r, "s1", config)
		fx.run(r.DeleteContext, d)
		if n := len(fx.api.Server.AllRequests()); n != 0 {
			t.Errorf("expected no API calls, got %d", n)
		}
	})
}
