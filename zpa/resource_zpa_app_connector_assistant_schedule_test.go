package zpa

/*
import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAppConnectorAssistantSchedule_Basic(t *testing.T) {
	customerID := os.Getenv("ZPA_CUSTOMER_ID")
	if customerID == "" {
		t.Fatal("ZPA_CUSTOMER_ID must be set for acceptance tests")
	}

	resourceTypeAndName := "zpa_app_connector_assistant_schedule.this"
	initialConfig := testAccAppConnectorAssistantScheduleConfig(customerID, "true")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConnectorAssistantScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppConnectorAssistantScheduleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "delete_disabled", "true"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "frequency", "days"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "customer_id", customerID),
				),
			},
		},
	})
}

func testAccAppConnectorAssistantScheduleConfig(customerID, deleteDisabled string) string {
	return fmt.Sprintf(`
resource "zpa_app_connector_assistant_schedule" "this" {
  enabled = true
  delete_disabled = %s
  frequency = "days"
  frequency_interval = "5"
  customer_id = "%s"
}
`, deleteDisabled, customerID)
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
