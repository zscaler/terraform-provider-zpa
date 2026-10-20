package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

func TestDataSourceAppConnectorAssistantSchedule(t *testing.T) {
	testAssistantScheduleDataSource(t, dataSourceAppConnectorAssistantSchedule(), "connectorSchedule")
}

// testAssistantScheduleDataSource covers the app connector and service edge
// schedule data sources: the tenant has one schedule, read from the endpoint.
func testAssistantScheduleDataSource(t *testing.T, r *schema.Resource, endpoint string) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path(endpoint), common.SuccessResponse(map[string]any{
		"id": "s1", "customerId": common.TestCustomerID, "enabled": true, "deleteDisabled": true,
		"frequency": "days", "frequencyInterval": "5",
	}))

	d := fx.resourceData(r, "", map[string]any{})
	fx.run(r.ReadContext, d)
	if d.Id() != "s1" {
		t.Errorf("id = %q, want s1", d.Id())
	}
	assertState(t, d, map[string]any{"frequency_interval": "5", "delete_disabled": true})
}
