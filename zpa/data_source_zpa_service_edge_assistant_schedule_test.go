package zpa

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceServiceEdgeAssistantSchedule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceServiceEdgeAssistantScheduleConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceServiceEdgeAssistantScheduleCheck("data.zpa_service_edge_assistant_schedule.this"),
					testAccDataSourceServiceEdgeAssistantScheduleCheck("data.zpa_service_edge_assistant_schedule.by_id"),
					testAccDataSourceServiceEdgeAssistantScheduleCheck("data.zpa_service_edge_assistant_schedule.customer_id"),
				),
			},
		},
	})
}

func testAccDataSourceServiceEdgeAssistantScheduleCheck(id string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(id, "id"),
	)
}

var testAccCheckDataSourceServiceEdgeAssistantScheduleConfig_basic = `
data "zpa_service_edge_assistant_schedule" "this" {
}

data "zpa_service_edge_assistant_schedule" "by_id" {
	id = data.zpa_service_edge_assistant_schedule.this.id
}

data "zpa_service_edge_assistant_schedule" "customer_id" {
	customer_id = data.zpa_service_edge_assistant_schedule.this.customer_id
}
`
*/
