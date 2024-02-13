package zpa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_profile"
)

func TestAccResourceInspectionProfile_Basic(t *testing.T) {
	var profile inspection_profile.InspectionProfile
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionProfile)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInspectionProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInspectionProfileConfigure(resourceTypeAndName, initialName, variable.InspectionProfileDescription, variable.InspectionProfileParanoia),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.InspectionProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", variable.InspectionProfileParanoia),
					resource.TestCheckResourceAttr(resourceTypeAndName, "predefined_controls.#", "7"),
				),
				ExpectNonEmptyPlan: true,
			},

			// Update test
			{
				Config: testAccCheckInspectionProfileConfigure(resourceTypeAndName, updatedName, variable.InspectionProfileDescriptionUpdate, variable.InspectionProfileParanoiaUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.InspectionProfileDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", variable.InspectionProfileParanoiaUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "predefined_controls.#", "7"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckInspectionProfileDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAInspectionProfile {
			continue
		}

		rule, _, err := apiClient.inspection_profile.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("inspection profile with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckInspectionProfileExists(resource string, rule *inspection_profile.InspectionProfile) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedProfile, _, err := apiClient.inspection_profile.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedProfile

		return nil
	}
}

func testAccCheckInspectionProfileConfigure(resourceTypeAndName, generatedName, description, paranoia string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

data "zpa_inspection_all_predefined_controls" "default_predefined_controls" {
	version    = "OWASP_CRS/3.3.0"
	group_name = "preprocessors"
}

data "zpa_inspection_predefined_controls" "this" {
    name = "Failed to parse request body"
    version = "OWASP_CRS/3.3.0"
}

resource "%s" "%s" {
	name                          = "%s"
	description                   = "%s"
	paranoia_level                = "%s"

	dynamic "predefined_controls" {
		for_each = data.zpa_inspection_all_predefined_controls.default_predefined_controls.list
		content {
		id           = predefined_controls.value.id
		action       = predefined_controls.value.action == "" ? predefined_controls.value.default_action : predefined_controls.value.action
		action_value = predefined_controls.value.action_value
		}
	}

	predefined_controls {
		id     = data.zpa_inspection_predefined_controls.this.id
		action = "BLOCK"
		protocol_type = "HTTP"
	}
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.ZPAInspectionProfile,
		resourceName,
		generatedName,
		description,
		paranoia,

		// Data source type and name
		resourcetype.ZPAInspectionProfile, resourceName,

		// Reference to the resource
		resourcetype.ZPAInspectionProfile, resourceName,
	)
}
