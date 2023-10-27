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

func TestAccResourceInspectionCustomControlsBasic(t *testing.T) {
	var control inspection_custom_controls.InspectionCustomControl
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionCustomControl)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInspectionCustomControlsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInspectionCustomControlsConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionCustomControlsExists(resourceTypeAndName, &control),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.InspectionCustomControlAction),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "default_action", variable.InspectionCustomControlDefaultAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocol_type", "HTTP"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "severity", variable.InspectionCustomControlSeverity),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.InspectionCustomControlType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "rules.#", "2"),
				),
				ExpectNonEmptyPlan: true,
			},

			// Update test
			{
				Config: testAccCheckInspectionCustomControlsConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionCustomControlsExists(resourceTypeAndName, &control),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.InspectionCustomControlAction),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "default_action", variable.InspectionCustomControlDefaultAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocol_type", "HTTP"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "severity", variable.InspectionCustomControlSeverity),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.InspectionCustomControlType),
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

		rule, _, err := apiClient.inspection_custom_controls.Get(rs.Primary.ID)

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
		receivedControl, _, err := apiClient.inspection_custom_controls.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedControl

		return nil
	}
}

func testAccCheckInspectionCustomControlsConfigure(resourceTypeAndName, generatedName string) string {
	return fmt.Sprintf(`
// inspection custom control resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getInspectionCustomControlsResourceHCL(generatedName),

		// data source variables
		resourcetype.ZPAInspectionCustomControl,
		generatedName,
		resourceTypeAndName,
	)
}

func getInspectionCustomControlsResourceHCL(generatedName string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name           = "tf-acc-test-%s"
	description    = "tf-acc-test-%s"
	action         = "%s"
	default_action = "%s"
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
		generatedName,

		variable.InspectionCustomControlAction,
		variable.InspectionCustomControlDefaultAction,
		variable.InspectionCustomControlSeverity,
		variable.InspectionCustomControlType,
	)
}
