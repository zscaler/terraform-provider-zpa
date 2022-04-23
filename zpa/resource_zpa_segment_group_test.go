package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/segmentgroup"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccResourceSegmentGroupBasic(t *testing.T) {
	var segmentGroup segmentgroup.SegmentGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSegmentGroupConfigure(resourceTypeAndName, generatedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentGroupExists(resourceTypeAndName, &segmentGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SegmentGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.SegmentGroupEnabled)),
				),
			},

			// Update test
			{
				Config: testAccCheckSegmentGroupConfigure(resourceTypeAndName, generatedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentGroupExists(resourceTypeAndName, &segmentGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SegmentGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.SegmentGroupEnabled)),
				),
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

		group, _, err := apiClient.segmentgroup.Get(rs.Primary.ID)

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
		receivedGroup, _, err := apiClient.segmentgroup.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*group = *receivedGroup

		return nil
	}
}

func testAccCheckSegmentGroupConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	return fmt.Sprintf(`
// segment group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		SegmentGroupResourceHCL(generatedName, description, enabled),

		// data source variables
		resourcetype.ZPASegmentGroup,
		generatedName,
		resourceTypeAndName,
	)
}

func SegmentGroupResourceHCL(generatedName, description string, enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = "%s"
}
`,
		// resource variables
		resourcetype.ZPASegmentGroup,
		generatedName,
		generatedName,
		// variable.SegmentGroupResourceName,
		description,
		strconv.FormatBool(enabled),
	)
}
