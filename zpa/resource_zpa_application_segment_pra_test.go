package zpa

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
)

func TestAccResourceApplicationSegmentPRA_Basic(t *testing.T) {
	var appSegment applicationsegmentpra.AppSegmentPRA
	appSegmentTypeAndName, _, appSegmentGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationSegmentPRA)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, "tf-acc-test-"+segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationSegmentPRADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationSegmentPRAConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, appSegmentGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, variable.AppSegmentEnabled, variable.AppSegmentCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentPRAExists(appSegmentTypeAndName, &appSegment),
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
				Config: testAccCheckApplicationSegmentPRAConfigure(appSegmentTypeAndName, appSegmentGeneratedName, appSegmentGeneratedName, appSegmentGeneratedName, segmentGroupHCL, segmentGroupTypeAndName, variable.AppSegmentEnabled, variable.AppSegmentCnameEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationSegmentPRAExists(appSegmentTypeAndName, &appSegment),
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
			{
				ResourceName:      appSegmentTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckApplicationSegmentPRADestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAApplicationSegmentPRA {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		appSegment, _, err := applicationsegmentpra.Get(context.Background(), service, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if appSegment != nil {
			return fmt.Errorf("pra application segment with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckApplicationSegmentPRAExists(resource string, segment *applicationsegmentpra.AppSegmentPRA) resource.TestCheckFunc {
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

		receivedSegment, _, err := applicationsegmentpra.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*segment = *receivedSegment

		return nil
	}
}

func testAccCheckApplicationSegmentPRAConfigure(resourceTypeAndName, generatedName, name, description, segmentGroupHCL, segmentGroupTypeAndName string, enabled bool, cnameEnabled bool) string {
	return fmt.Sprintf(`

// segment group resource
%s

// application segment pra resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		segmentGroupHCL,
		// serverGroupHCL,
		getApplicationSegmentPRAResourceHCL(generatedName, name, description, segmentGroupTypeAndName, enabled, cnameEnabled),

		// data source variables
		resourcetype.ZPAApplicationSegmentPRA,
		generatedName,
		resourceTypeAndName,
	)
}

func getApplicationSegmentPRAResourceHCL(generatedName, name, description, segmentGroupTypeAndName string, enabled bool, cnameEnabled bool) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled = "%s"
	is_cname_enabled = "%s"
	select_connector_close_to_app = true
	health_reporting = "ON_ACCESS"
	bypass_type = "NEVER"
	tcp_port_ranges = ["3389", "3389"]
	domain_names = ["rdp_pra.example.com"]
	segment_group_id = "${%s.id}"
	tcp_keep_alive = "1"
	common_apps_dto {
		  apps_config {
		  domain               = "rdp_pra.example.com"
		  application_protocol = "RDP"
		  connection_security  = "ANY"
		  application_port     = "3389"
		  enabled = true
		  app_types = ["SECURE_REMOTE_ACCESS"]
		}
	}
	depends_on = [ %s ]
}
`,
		// resource variables
		resourcetype.ZPAApplicationSegmentPRA,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		strconv.FormatBool(cnameEnabled),
		segmentGroupTypeAndName,
		segmentGroupTypeAndName,
	)
}
