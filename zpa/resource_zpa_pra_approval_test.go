package zpa

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
)

func TestAccResourcePRAPrivilegedApprovalController_Basic(t *testing.T) {
	var praApproval praapproval.PrivilegedApproval
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAApprovalController)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAPrivilegedApprovalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAPrivilegedApprovalConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRAPrivilegedApprovalExists(resourceTypeAndName, &praApproval),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "applications.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "working_hours.#", "1"),
				),
			},
			// Update test
			{
				Config: testAccCheckPRAPrivilegedApprovalConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRAPrivilegedApprovalExists(resourceTypeAndName, &praApproval),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "applications.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "working_hours.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPRAPrivilegedApprovalDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPRAApprovalController {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		approval, _, err := praapproval.Get(context.Background(), service, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if approval != nil {
			return fmt.Errorf("pra approval with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPRAPrivilegedApprovalExists(resource string, approval *praapproval.PrivilegedApproval) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		receivedApproval, _, err := praapproval.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*approval = *receivedApproval

		return nil
	}
}

func testAccCheckPRAPrivilegedApprovalConfigure(resourceTypeAndName, generatedName string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	// Generate start_time and end_time dynamically
	now := time.Now()
	// Ensure start_time is not more than 1 hour in the past
	startTime := now.Add(-30 * time.Minute).Format(time.RFC1123)
	// Ensure end_time is no more than 365 days from start_time
	endTime := now.AddDate(0, 3, 0).Format(time.RFC1123) // Example: 3 months from now

	userEmails := []string{
		"kathy.kavanagh@bd-hashicorp.com",
	}
	emailIDs := fmt.Sprintf(`["%s"]`, strings.Join(userEmails, `", "`))
	return fmt.Sprintf(`

resource "zpa_application_segment_pra" "this" {
	name             = "tf-acc-test-%s"
	description      = "tf-acc-test-%s"
	enabled          = true
	health_reporting = "ON_ACCESS"
	bypass_type      = "NEVER"
	is_cname_enabled = true
	tcp_port_ranges  = ["3222", "3222", "3391", "3391"]
	domain_names     = ["ssh_pra3222.example.com", "rdp_pra3391.example.com"]
	segment_group_id = zpa_segment_group.this.id
	common_apps_dto {
		apps_config {
		domain               = "rdp_pra3391.example.com"
		application_protocol = "RDP"
		connection_security  = "ANY"
		application_port     = "3391"
		enabled              = true
		app_types            = ["SECURE_REMOTE_ACCESS"]
		}
		apps_config {
		domain               = "ssh_pra3222.example.com"
		application_protocol = "SSH"
		application_port     = "3222"
		enabled              = true
		app_types            = ["SECURE_REMOTE_ACCESS"]
		}
	}
}

resource "zpa_segment_group" "this" {
	name        = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled     = true
}

resource "%s" "%s" {
	email_ids = %s
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

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource name for pra application segment and segment group
		generatedName,
		generatedName,
		generatedName,
		generatedName,

		// Resource type and name for privileged approval
		resourcetype.ZPAPRAApprovalController,
		resourceName,
		emailIDs,
		startTime,
		endTime,

		// Data source type and name
		resourcetype.ZPAPRAApprovalController, resourceName,

		// Reference to the resource
		resourcetype.ZPAPRAApprovalController, resourceName,
	)
}
