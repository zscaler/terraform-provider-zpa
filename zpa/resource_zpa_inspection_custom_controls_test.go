package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_custom_controls"
)

func TestAccResourceInspectionCustomControls_Basic(t *testing.T) {
	var control inspection_custom_controls.InspectionCustomControl
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionCustomControl)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInspectionCustomControlsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInspectionCustomControlsConfigure(resourceTypeAndName, generatedName, variable.CustomControlDescription, variable.CustomControlSeverity, variable.CustomControlControlType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionCustomControlsExists(resourceTypeAndName, &control),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CustomControlDescription),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "action", "BLOCK"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "default_action", variable.InspectionCustomControlDefaultAction),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", variable.CustomControlParanoia),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocol_type", "HTTP"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "severity", variable.CustomControlSeverity),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.CustomControlControlType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "rules.#", "2"),
				),
				ExpectNonEmptyPlan: true,
			},

			// Update test
			{
				Config: testAccCheckInspectionCustomControlsConfigure(resourceTypeAndName, generatedName, variable.CustomControlDescriptionUpdate, variable.CustomControlSeverityUpdate, variable.CustomControlControlType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionCustomControlsExists(resourceTypeAndName, &control),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CustomControlDescriptionUpdate),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "action", "PASS"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "default_action", variable.CustomControlDefaultAction),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", variable.CustomControlParanoiaUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocol_type", "HTTP"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "severity", variable.CustomControlSeverityUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.CustomControlControlType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "rules.#", "2"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckInspectionCustomControlsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAInspectionCustomControl {
			continue
		}

		rule, _, err := inspection_custom_controls.Get(apiClient.InspectionCustomControls, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("inspection custom control with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckInspectionCustomControlsExists(resource string, rule *inspection_custom_controls.InspectionCustomControl) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedControl, _, err := inspection_custom_controls.Get(apiClient.InspectionCustomControls, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedControl

		return nil
	}
}

func testAccCheckInspectionCustomControlsConfigure(resourceTypeAndName, generatedName, description, severity, controlType string) string {
	return fmt.Sprintf(`
// inspection custom control resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getInspectionCustomControlsResourceHCL(generatedName, description, severity, controlType),

		// data source variables
		resourcetype.ZPAInspectionCustomControl,
		generatedName,
		resourceTypeAndName,
	)
}

func getInspectionCustomControlsResourceHCL(generatedName, description, severity, controlType string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name           = "tf-acc-test-%s"
	description    = "%s"
	action         = "PASS"
	default_action = "PASS"
	paranoia_level = "1"
	protocol_type  = "HTTP"
	severity       = "%s"
	type           = "%s"
	rules {
	  names = ["test"]
	  type  = "RESPONSE_HEADERS"
	  conditions {
		lhs = "SIZE"
		op  = "GE"
		rhs = "1000"
	  }
	}
	rules {
		type  = "RESPONSE_BODY"
		conditions {
		  lhs = "SIZE"
		  op  = "GE"
		  rhs = "1000"
		}
	}
  }
`,
		// resource variables
		resourcetype.ZPAInspectionCustomControl,
		generatedName,
		generatedName,

		description,
		// action,
		severity,
		controlType,
	)
}
