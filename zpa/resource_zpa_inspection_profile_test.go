package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/inspectioncontrol/inspection_profile"
)

func TestAccResourceInspectionProfileBasic(t *testing.T) {
	var profile inspection_profile.InspectionProfile
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionProfile)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInspectionProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInspectionProfileConfigure(resourceTypeAndName, generatedName, variable.InspectionProfileDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.InspectionProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "predefined_controls.#", "7"),
				),
			},

			// Update test
			{
				Config: testAccCheckInspectionProfileConfigure(resourceTypeAndName, generatedName, variable.InspectionProfileDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInspectionProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.InspectionProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "paranoia_level", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "predefined_controls.#", "7"),
				),
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

func testAccCheckInspectionProfileConfigure(resourceTypeAndName, generatedName, description string) string {
	return fmt.Sprintf(`
// inspection profile resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getInspectionProfileResourceHCL(generatedName, description),

		// data source variables
		resourcetype.ZPAInspectionProfile,
		generatedName,
		resourceTypeAndName,
	)
}

func getInspectionProfileResourceHCL(generatedName, description string) string {
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
	paranoia_level              = "1"

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
	}
}
`,
		// resource variables
		resourcetype.ZPAInspectionProfile,
		generatedName,
		generatedName,
		description,
	)
}
