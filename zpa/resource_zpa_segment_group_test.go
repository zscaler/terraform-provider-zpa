package zpa

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func TestAccResourceSegmentGroup_Basic(t *testing.T) {
	var segmentGroup segmentgroup.SegmentGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSegmentGroupConfigure(resourceTypeAndName, initialName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentGroupExists(resourceTypeAndName, &segmentGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SegmentGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.SegmentGroupEnabled)),
				),
			},

			// Update test
			{
				Config: testAccCheckSegmentGroupConfigure(resourceTypeAndName, updatedName, variable.SegmentGroupDescriptionUpdate, variable.SegmentGroupEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentGroupExists(resourceTypeAndName, &segmentGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SegmentGroupDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.SegmentGroupEnabledUpdate)),
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

func testAccCheckSegmentGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPASegmentGroup {
			continue
		}

		group, _, err := segmentgroup.Get(apiClient.SegmentGroup, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if group != nil {
			return fmt.Errorf("segment group with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSegmentGroupExists(resource string, group *segmentgroup.SegmentGroup) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedGroup, _, err := segmentgroup.Get(apiClient.SegmentGroup, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*group = *receivedGroup

		return nil
	}
}

func testAccCheckSegmentGroupConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the certificate
		resourcetype.ZPASegmentGroup,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),

		// Data source type and name
		resourcetype.ZPASegmentGroup, resourceName,

		// Reference to the resource
		resourcetype.ZPASegmentGroup, resourceName,
	)
}
