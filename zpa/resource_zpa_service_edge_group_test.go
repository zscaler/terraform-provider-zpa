package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/serviceedgegroup"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccResourceServiceEdgeGroupBasic(t *testing.T) {
	var groups serviceedgegroup.ServiceEdgeGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServiceEdgeGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceEdgeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, variable.ServiceEdgeDescription, variable.ServiceEdgeEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEdgeGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ServiceEdgeDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServiceEdgeEnabled)),
				),
			},

			// Update test
			{
				Config: testAccCheckServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, variable.ServiceEdgeDescription, variable.ServiceEdgeEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEdgeGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ServiceEdgeDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServiceEdgeEnabled)),
				),
			},
		},
	})
}

func testAccCheckServiceEdgeGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAServiceEdgeGroup {
			continue
		}

		group, _, err := apiClient.serviceedgegroup.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if group != nil {
			return fmt.Errorf("service edge group with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckServiceEdgeGroupExists(resource string, group *serviceedgegroup.ServiceEdgeGroup) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedGroup, _, err := apiClient.serviceedgegroup.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*group = *receivedGroup

		return nil
	}
}

func testAccCheckServiceEdgeGroupConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	return fmt.Sprintf(`
// service edge group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		ServiceEdgeGroupResourceHCL(generatedName, description, enabled),

		// data source variables
		resourcetype.ZPAServiceEdgeGroup,
		generatedName,
		resourceTypeAndName,
	)
}

func ServiceEdgeGroupResourceHCL(generatedName, description string, enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name                 = "%s"
	description          = "%s"
	enabled				 = "%s"
	upgrade_day          = "SUNDAY"
	upgrade_time_in_secs = "66600"
	latitude             = "37.3382082"
	longitude            = "-121.8863286"
	location             = "San Jose, CA, USA"
	version_profile_id   = "0"
}
`,
		// resource variables
		resourcetype.ZPAServiceEdgeGroup,
		generatedName,
		variable.ServiceEdgeResourceName,
		description,
		strconv.FormatBool(enabled),
	)
}
