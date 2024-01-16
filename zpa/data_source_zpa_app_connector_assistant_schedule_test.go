package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAppConnectorAssistantSchedule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAppConnectorAssistantScheduleConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAppConnectorAssistantScheduleCheck("data.zpa_app_connector_assistant_schedule.this"),
					testAccDataSourceAppConnectorAssistantScheduleCheck("data.zpa_app_connector_assistant_schedule.by_id"),
					testAccDataSourceAppConnectorAssistantScheduleCheck("data.zpa_app_connector_assistant_schedule.customer_id"),
				),
			},
		},
	})
}

func testAccDataSourceAppConnectorAssistantScheduleCheck(id string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(id, "id"),
	)
}

var testAccCheckDataSourceAppConnectorAssistantScheduleConfig_basic = `
data "zpa_app_connector_assistant_schedule" "this" {
}

data "zpa_app_connector_assistant_schedule" "by_id" {
	id = data.zpa_app_connector_assistant_schedule.this.id
}

data "zpa_app_connector_assistant_schedule" "customer_id" {
	customer_id = data.zpa_app_connector_assistant_schedule.this.customer_id
}
`*/
