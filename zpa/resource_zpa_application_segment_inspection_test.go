package zpa

/*
import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/applicationsegmentinspection"
)

func TestAccResourceApplicationSegmentInspectionBasic(t *testing.T) {
	var appSegment applicationsegmentinspection.AppSegmentInspection
	appSegmentTypeAndName, _, appSegmentGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationSegmentInspection)

	serverGroupTypeAndName, _, serverGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)
	serverGroupHCL := testAccCheckServerGroupConfigure(serverGroupTypeAndName, serverGroupGeneratedName, "", "", "", "", variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)
	rDomain := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationSegmentInspectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationSegmentInspectionConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, appSegmentGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName, variable.AppSegmentEnabled, rDomain, variable.AppSegmentCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentInspectionExists(appSegmentTypeAndName, &appSegment),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "name", "tf-acc-test-"+appSegmentGeneratedName),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "description", "tf-acc-test-"+appSegmentGeneratedName),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "enabled", strconv.FormatBool(variable.AppSegmentEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.AppSegmentCnameEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttrSet(appSegmentTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "common_apps_dto.#", "1"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "tcp_port_ranges.#", "2"),
				),
				ExpectNonEmptyPlan: true,
			},

			// Update test
			{
				Config: testAccCheckApplicationSegmentInspectionConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, appSegmentGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName, variable.AppSegmentEnabled, rDomain, variable.AppSegmentCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentInspectionExists(appSegmentTypeAndName, &appSegment),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "name", "tf-acc-test-"+appSegmentGeneratedName),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "description", "tf-acc-test-"+appSegmentGeneratedName),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "enabled", strconv.FormatBool(variable.AppSegmentEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.AppSegmentCnameEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttrSet(appSegmentTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "common_apps_dto.#", "1"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "tcp_port_ranges.#", "2"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckApplicationSegmentInspectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAApplicationSegmentInspection {
			continue
		}

		_, _, err := client.applicationsegmentinspection.GetByName(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Inspection Application Segment still exists")
		}

		return nil
	}
	return nil
}

func testAccCheckApplicationSegmentInspectionExists(resource string, segment *applicationsegmentinspection.AppSegmentInspection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Inspection Application Segment Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Inspection Application Segment ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.applicationsegmentinspection.GetByName(rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}
		if resp.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("name Not found in created attributes")
		}
		if resp.Description != rs.Primary.Attributes["description"] {
			return fmt.Errorf("description Not found in created attributes")
		}
		return nil
	}
}

func testAccCheckApplicationSegmentInspectionConfigure(resourceTypeAndName, generatedName, name, description, segmentGroupHCL, segmentGroupTypeAndName, serverGroupHCL, serverGroupTypeAndName string, enabled bool,rDomain string, cnameEnabled bool) string {
	return fmt.Sprintf(`

// segment group resource
%s

// application segment resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		segmentGroupHCL,
		// serverGroupHCL,
		getApplicationSegmentInspectionResourceHCL(generatedName, name, description, segmentGroupTypeAndName, serverGroupTypeAndName, enabled, rDomain, cnameEnabled),

		// data source variables
		resourcetype.ZPAApplicationSegmentInspection,
		generatedName,
		resourceTypeAndName,
	)
}

func getApplicationSegmentInspectionResourceHCL(generatedName, name, description, segmentGroupTypeAndName, serverGroupTypeAndName string, enabled bool, rDomain string, cnameEnabled bool) string {
	return fmt.Sprintf(`

data "zpa_ba_certificate" "sales" {
	name = "sales.bd-hashicorp.com"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled = "%s"
	is_cname_enabled = "%s"
	select_connector_close_to_app = true
	health_reporting = "ON_ACCESS"
	bypass_type = "NEVER"
	tcp_keep_alive = "1"
	tcp_port_range {
		from = "4443"
		to = "4443"
	}
	domain_names = ["sales.bd-hashicorp.com"]
	segment_group_id = "${%s.id}"
	common_apps_dto {
		apps_config {
		  name                 = "%s.bd-hashicorp.com"
		  domain               = "%s.bd-hashicorp.com"
		  application_protocol = "HTTPS"
		  application_port     = "4443"
		  certificate_id       = data.zpa_ba_certificate.sales.id
		  enabled = true
		  app_types = ["INSPECT"]
		}
	}
	server_groups {
		id = []
	}
	depends_on = [ %s ]
}
`,

		// resource variables
		resourcetype.ZPAApplicationSegmentInspection,
		generatedName,
		generatedName,
		generatedName,
		strconv.FormatBool(enabled),
		strconv.FormatBool(cnameEnabled),
		// generatedName,
		segmentGroupTypeAndName,
		generatedName,
		generatedName,
		// serverGroupTypeAndName,
		segmentGroupTypeAndName,
		// serverGroupTypeAndName,
	)
}
*/