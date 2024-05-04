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
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgegroup"
)

func TestAccResourceServiceEdgeGroupBasic(t *testing.T) {
	var groups serviceedgegroup.ServiceEdgeGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServiceEdgeGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceEdgeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServiceEdgeGroupConfigure(resourceTypeAndName, initialName, variable.ServiceEdgeDescription, variable.ServiceEdgeEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEdgeGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ServiceEdgeDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServiceEdgeEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_public", strconv.FormatBool(variable.ServiceEdgeIsPublic)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "latitude", variable.ServiceEdgeLatitude),
					resource.TestCheckResourceAttr(resourceTypeAndName, "longitude", variable.ServiceEdgeLongitude),
					resource.TestCheckResourceAttr(resourceTypeAndName, "location", variable.ServiceEdgeLocation),
					resource.TestCheckResourceAttr(resourceTypeAndName, "version_profile_name", variable.ServiceEdgeVersionProfileName),
				),
			},

			// Update test
			{
				Config: testAccCheckServiceEdgeGroupConfigure(resourceTypeAndName, updatedName, variable.ServiceEdgeDescription, variable.ServiceEdgeEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceEdgeGroupExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ServiceEdgeDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServiceEdgeEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_public", strconv.FormatBool(variable.ServiceEdgeIsPublic)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "latitude", variable.ServiceEdgeLatitude),
					resource.TestCheckResourceAttr(resourceTypeAndName, "longitude", variable.ServiceEdgeLongitude),
					resource.TestCheckResourceAttr(resourceTypeAndName, "location", variable.ServiceEdgeLocation),
					resource.TestCheckResourceAttr(resourceTypeAndName, "version_profile_name", variable.ServiceEdgeVersionProfileName),
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
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name                      = "%s"
	description               = "%s"
	enabled				      = "%s"
	is_public			      = "%s"
	upgrade_day               = "SUNDAY"
	upgrade_time_in_secs      = "66600"
	country_code              = "US"
	city_country              = "San Jose, US"
	latitude                  = "37.33874"
	longitude                 = "-121.8852525"
	location                  = "San Jose, CA, USA"
	version_profile_id        = 0
	grace_distance_enabled    = true
	grace_distance_value      = "10"
	grace_distance_value_unit = "KMS"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the certificate
		// resource variables
		resourcetype.ZPAServiceEdgeGroup,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		strconv.FormatBool(enabled),

		// Data source type and name
		resourcetype.ZPAServiceEdgeGroup,
		resourceName,

		// Reference to the resource
		resourcetype.ZPAServiceEdgeGroup, resourceName,
	)
}
