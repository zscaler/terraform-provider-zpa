package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegment"
)

func TestAccResourceApplicationSegment_Basic(t *testing.T) {
	var appSegment applicationsegment.ApplicationSegmentResource
	appSegmentTypeAndName, _, appSegmentGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationSegment)
	rPort := acctest.RandIntRange(1000, 9999)
	rDescription := acctest.RandomWithPrefix("tf-acc-test-")
	updatedDescription := acctest.RandomWithPrefix("tf-updated-") // New name for update test

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, "tf-acc-test-"+segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationSegmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationSegmentConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, rDescription, segmentGroupHCL, segmentGroupTypeAndName, rPort, variable.AppSegmentEnabled, variable.AppSegmentCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentExists(appSegmentTypeAndName, &appSegment),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "name", "tf-acc-test-"+appSegmentGeneratedName),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "description", rDescription),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "enabled", strconv.FormatBool(variable.AppSegmentEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.AppSegmentCnameEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttrSet(appSegmentTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "tcp_port_range.#", "1"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "udp_port_range.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckApplicationSegmentConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, updatedDescription, segmentGroupHCL, segmentGroupTypeAndName, rPort, variable.AppSegmentEnabled, variable.AppSegmentCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentExists(appSegmentTypeAndName, &appSegment),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "name", "tf-acc-test-"+appSegmentGeneratedName),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "description", updatedDescription),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "enabled", strconv.FormatBool(variable.AppSegmentEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "is_cname_enabled", strconv.FormatBool(variable.AppSegmentCnameEnabled)),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "bypass_type", "NEVER"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "health_reporting", "ON_ACCESS"),
					resource.TestCheckResourceAttrSet(appSegmentTypeAndName, "segment_group_id"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "tcp_port_range.#", "1"),
					resource.TestCheckResourceAttr(appSegmentTypeAndName, "udp_port_range.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      appSegmentTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckApplicationSegmentDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAApplicationSegment {
			continue
		}

		_, _, err := applicationsegment.GetByName(apiClient.ApplicationSegment, rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Application Segment still exists")
		}

		return nil
	}
	return nil
}

func testAccCheckApplicationSegmentExists(resource string, segment *applicationsegment.ApplicationSegmentResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Application Segment Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Application Segment ID is set")
		}
		apiClient := testAccProvider.Meta().(*Client)
		receivedApp, _, err := applicationsegment.Get(apiClient.ApplicationSegment, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*segment = *receivedApp

		return nil
	}
}

func testAccCheckApplicationSegmentConfigure(resourceTypeAndName, generatedName, name, description, segmentGroupHCL, segmentGroupTypeAndName string, rPort int, enabled, cnameEnabled bool) string {
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
		getApplicationSegmentResourceHCL(generatedName, name, description, segmentGroupTypeAndName, rPort, enabled, cnameEnabled),

		// data source variables
		resourcetype.ZPAApplicationSegment,
		generatedName,
		resourceTypeAndName,
	)
}

func getApplicationSegmentResourceHCL(generatedName, name, description, segmentGroupTypeAndName string, rPort int, enabled, cnameEnabled bool) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	enabled = "%s"
	is_cname_enabled = "%s"
	health_reporting = "ON_ACCESS"
	bypass_type = "NEVER"
	tcp_port_range = [
	  {
		from = "%d"
		to   = "%d"
	  }
	]
	udp_port_range = [
	  {
		from = "%d"
		to   = "%d"
	  }
	]
	domain_names = ["test.example.com"]
	segment_group_id = "${%s.id}"
	tcp_keep_alive = "1"
	server_groups {
		id = []
	}
	depends_on = [ %s ]
}
`,

		// resource variables
		resourcetype.ZPAApplicationSegment,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		strconv.FormatBool(cnameEnabled),
		rPort,
		rPort,
		rPort,
		rPort,
		segmentGroupTypeAndName,
		// serverGroupTypeAndName,
		segmentGroupTypeAndName,
		// serverGroupTypeAndName,
	)
}
