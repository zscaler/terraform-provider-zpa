package zpa

/*
import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAppConnectorAssistantSchedule_basic(t *testing.T) {
	customerID := os.Getenv("ZPA_CUSTOMER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConnectorAssistantScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppConnectorAssistantScheduleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorAssistantScheduleExists("zpa_app_connector_assistant_schedule.this"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "enabled", "true"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "delete_disabled", "true"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "frequency", "days"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "frequency_interval", "7"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "customer_id", customerID),
				),
			},
			// Test resource update
			{
				Config: testAccAppConnectorAssistantScheduleConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorAssistantScheduleExists("zpa_app_connector_assistant_schedule.this"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "enabled", "false"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "delete_disabled", "false"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "frequency_interval", "14"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "customer_id", customerID),
				),
			},
			// Test resource update
			{
				Config: testAccAppConnectorAssistantScheduleConfigEnabled(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorAssistantScheduleExists("zpa_app_connector_assistant_schedule.this"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "enabled", "true"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "delete_disabled", "true"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "frequency_interval", "30"),
					resource.TestCheckResourceAttr("zpa_app_connector_assistant_schedule.this", "customer_id", customerID),
				),
			},
		},
	})
}

func testAccAppConnectorAssistantScheduleConfig() string {
	return `
resource "zpa_app_connector_assistant_schedule" "this" {
  enabled = true
  delete_disabled = true
  frequency = "days"
  frequency_interval = "7"
}
`
}

func testAccAppConnectorAssistantScheduleConfigUpdated() string {
	return `
resource "zpa_app_connector_assistant_schedule" "this" {
  enabled = false
  delete_disabled = false
  frequency = "days"
  frequency_interval = "14"
}
`
}

func testAccAppConnectorAssistantScheduleConfigEnabled() string {
	return `
resource "zpa_app_connector_assistant_schedule" "this" {
  enabled = true
  delete_disabled = true
  frequency = "days"
  frequency_interval = "30"
}
`
}

func testAccCheckAppConnectorAssistantScheduleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Connector Assistant Schedule ID is set")
		}

		return nil
	}
}

func testAccCheckAppConnectorAssistantScheduleDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}
*/
