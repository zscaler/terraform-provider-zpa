package resources_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccPRAPrivilegedApproval_basic(t *testing.T) {
	var approval praapproval.PrivilegedApproval
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_pra_approval_controller.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPRAPrivilegedApprovalDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPRAPrivilegedApprovalConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPRAPrivilegedApprovalExists(zClient, resourceName, &approval),
					resource.TestCheckResourceAttr(resourceName, "email_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "applications.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "working_hours.#", "1"),
				),
			},
			// Update test
			{
				Config: testAccPRAPrivilegedApprovalConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPRAPrivilegedApprovalExists(zClient, resourceName, &approval),
					resource.TestCheckResourceAttr(resourceName, "email_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "applications.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "working_hours.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPRAPrivilegedApprovalExists(zClient *client.Client, resourceName string, approval *praapproval.PrivilegedApproval) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("pra privileged approval not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("pra privileged approval ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		received, _, err := praapproval.Get(ctx, service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch pra privileged approval %s: %w", rs.Primary.ID, err)
		}

		*approval = *received
		return nil
	}
}

func testAccCheckPRAPrivilegedApprovalDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_pra_approval_controller" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			approval, _, err := praapproval.Get(ctx, service, rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("id %s already exists", rs.Primary.ID)
			}

			if approval != nil {
				return fmt.Errorf("pra privileged approval with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccPRAPrivilegedApprovalConfig(rName string) string {
	// Generate start_time and end_time dynamically
	now := time.Now()
	// Ensure start_time is not more than 1 hour in the past
	startTime := now.Add(-30 * time.Minute).Format(time.RFC1123)
	// Ensure end_time is no more than 365 days from start_time
	endTime := now.AddDate(0, 3, 0).Format(time.RFC1123) // Example: 3 months from now

	return fmt.Sprintf(`
resource "zpa_application_segment_pra" "this" {
	name             = "tf-acc-test-%s"
	description      = "tf-acc-test-%s"
	enabled          = true
	health_reporting = "ON_ACCESS"
	bypass_type      = "NEVER"
	is_cname_enabled = true
	tcp_port_ranges  = ["3222", "3222"]
	domain_names     = ["rdp_pra3222.example.com"]
	segment_group_id = zpa_segment_group.this.id
}

resource "zpa_segment_group" "this" {
	name        = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled     = true
}

resource "zpa_pra_approval_controller" "test" {
	email_ids = ["kathy.kavanagh@securitygeek.io"]
	start_time = "%s"
	end_time = "%s"
	status = "ACTIVE"
	applications {
		id = [zpa_application_segment_pra.this.id]
	}
	working_hours {
		days = ["FRI", "MON", "SAT", "SUN", "THU", "TUE", "WED"]
		start_time = "00:10"
		start_time_cron = "0 0 8 ? * MON,TUE,WED,THU,FRI,SAT"
		end_time = "09:15"
		end_time_cron = "0 15 17 ? * MON,TUE,WED,THU,FRI,SAT"
		timezone = "America/Vancouver"
	}
}

data "zpa_pra_privileged_approval" "test" {
	id = zpa_pra_approval_controller.test.id
}
`, rName, rName, rName, rName, startTime, endTime)
}
