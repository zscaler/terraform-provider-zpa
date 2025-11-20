package datasources_test

import (
	"fmt"
	"testing"
	"time"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccPRAPrivilegedApprovalDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	resourceName := "zpa_pra_approval_controller.test"
	dataSourceName := "data.zpa_pra_privileged_approval.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPRAPrivilegedApprovalDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "email_ids.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "status", resourceName, "status"),
					resource.TestCheckResourceAttr(dataSourceName, "applications.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "working_hours.#", "1"),
				),
			},
		},
	})
}

func testAccPRAPrivilegedApprovalDataSourceConfig(rName string) string {
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
