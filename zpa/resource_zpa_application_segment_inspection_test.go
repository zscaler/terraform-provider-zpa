package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentinspection"
)

func TestAccResourceApplicationSegmentInspection_Basic(t *testing.T) {
	var appSegment applicationsegmentinspection.AppSegmentInspection
	appSegmentTypeAndName, _, appSegmentGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationSegmentInspection)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, "tf-acc-test-"+segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationSegmentInspectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationSegmentInspectionConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, appSegmentGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, variable.AppSegmentEnabled, variable.AppSegmentCnameEnabled),
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
			},

			// Update test
			{
				Config: testAccCheckApplicationSegmentInspectionConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, appSegmentGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, variable.AppSegmentEnabled, variable.AppSegmentCnameEnabled),
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
			},
			// Import test
			// {
			// 	ResourceName:      appSegmentTypeAndName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckApplicationSegmentInspectionDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAApplicationSegmentInspection {
			continue
		}

		_, _, err := applicationsegmentinspection.GetByName(apiClient.ApplicationSegmentInspection, rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Inspection Application Segment Inspection still exists")
		}

		return nil
	}
	return nil
}

func testAccCheckApplicationSegmentInspectionExists(resource string, segment *applicationsegmentinspection.AppSegmentInspection) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedSegment, _, err := applicationsegmentinspection.Get(apiClient.ApplicationSegmentInspection, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*segment = *receivedSegment

		return nil
	}
}

func testAccCheckApplicationSegmentInspectionConfigure(resourceTypeAndName, generatedName, name, description, segmentGroupHCL, segmentGroupTypeAndName string, enabled bool, cnameEnabled bool) string {
	return fmt.Sprintf(`

// segment group resource
%s

// application segment inspection resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		segmentGroupHCL,
		// serverGroupHCL,
		getApplicationSegmentInspectionResourceHCL(generatedName, name, description, segmentGroupTypeAndName, enabled, cnameEnabled),

		// data source variables
		resourcetype.ZPAApplicationSegmentInspection,
		generatedName,
		resourceTypeAndName,
	)
}

func getApplicationSegmentInspectionResourceHCL(generatedName, name, description, segmentGroupTypeAndName string, enabled bool, cnameEnabled bool) string {
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
	tcp_port_ranges  = ["443", "443"]
	domain_names = ["sales.bd-hashicorp.com"]
	segment_group_id = "${%s.id}"
	common_apps_dto {
		apps_config {
		  name                 = "sales.bd-hashicorp.com"
		  domain               = "sales.bd-hashicorp.com"
		  application_protocol = "HTTPS"
		  application_port     = "443"
		  certificate_id       = data.zpa_ba_certificate.sales.id
		  enabled = true
		  app_types = ["INSPECT"]
		}
	}
	depends_on = [ %s ]
}
`,

		// resource variables
		resourcetype.ZPAApplicationSegmentInspection,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		strconv.FormatBool(cnameEnabled),
		segmentGroupTypeAndName,
		segmentGroupTypeAndName,
	)
}
