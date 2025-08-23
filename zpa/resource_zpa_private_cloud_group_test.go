package zpa

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

func TestAccResourcePrivateCloudGroup_Basic(t *testing.T) {
	var privateCloudGroup private_cloud_group.PrivateCloudGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPrivateCloudGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateCloudGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPrivateCloudGroupConfigure(resourceTypeAndName, initialName, variable.PrivateCloudGroupDescription, variable.PrivateCloudGroupEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateCloudGroupExists(resourceTypeAndName, &privateCloudGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.PrivateCloudGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PrivateCloudGroupEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "city_country", variable.PrivateCloudGroupCityCountry),
					resource.TestCheckResourceAttr(resourceTypeAndName, "latitude", variable.PrivateCloudGroupLatitude),
					resource.TestCheckResourceAttr(resourceTypeAndName, "longitude", variable.PrivateCloudGroupLongitude),
					resource.TestCheckResourceAttr(resourceTypeAndName, "location", variable.PrivateCloudGroupLocation),
					resource.TestCheckResourceAttr(resourceTypeAndName, "upgrade_day", variable.PrivateCloudGroupUpgradeDay),
					resource.TestCheckResourceAttr(resourceTypeAndName, "upgrade_time_in_secs", variable.PrivateCloudGroupUpgradeTimeInSecs),
					resource.TestCheckResourceAttr(resourceTypeAndName, "version_profile_id", variable.PrivateCloudGroupVersionProfileID),
					resource.TestCheckResourceAttr(resourceTypeAndName, "override_version_profile", strconv.FormatBool(variable.PrivateCloudGroupOverrideVersionProfile)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_public", variable.PrivateCloudGroupIsPublic),
				),
			},

			// Update test
			{
				Config: testAccCheckPrivateCloudGroupConfigureUpdate(resourceTypeAndName, updatedName, variable.PrivateCloudGroupDescriptionUpdate, variable.PrivateCloudGroupEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateCloudGroupExists(resourceTypeAndName, &privateCloudGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.PrivateCloudGroupDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.PrivateCloudGroupEnabledUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "upgrade_day", variable.PrivateCloudGroupUpgradeDayUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "upgrade_time_in_secs", variable.PrivateCloudGroupUpgradeTimeInSecsUpdate),
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

func testAccCheckPrivateCloudGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPrivateCloudGroup {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		group, _, err := private_cloud_group.Get(context.Background(), service, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if group != nil {
			return fmt.Errorf("private cloud group with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPrivateCloudGroupExists(resource string, group *private_cloud_group.PrivateCloudGroup) resource.TestCheckFunc {
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

		receivedGroup, _, err := private_cloud_group.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*group = *receivedGroup

		return nil
	}
}

func testAccCheckPrivateCloudGroupConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = %s
	city_country = "%s"
	latitude = "%s"
	longitude = "%s"
	location = "%s"
	upgrade_day = "%s"
	upgrade_time_in_secs = "%s"
	version_profile_id = "%s"
	override_version_profile = %s
	is_public = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the private cloud group
		resourcetype.ZPAPrivateCloudGroup,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		variable.PrivateCloudGroupCityCountry,
		variable.PrivateCloudGroupLatitude,
		variable.PrivateCloudGroupLongitude,
		variable.PrivateCloudGroupLocation,
		variable.PrivateCloudGroupUpgradeDay,
		variable.PrivateCloudGroupUpgradeTimeInSecs,
		variable.PrivateCloudGroupVersionProfileID,
		strconv.FormatBool(variable.PrivateCloudGroupOverrideVersionProfile),
		variable.PrivateCloudGroupIsPublic,

		// Data source type and name
		resourcetype.ZPAPrivateCloudGroup, resourceName,

		// Reference to the resource
		resourcetype.ZPAPrivateCloudGroup, resourceName,
	)
}

func testAccCheckPrivateCloudGroupConfigureUpdate(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = %s
	city_country = "%s"
	latitude = "%s"
	longitude = "%s"
	location = "%s"
	upgrade_day = "%s"
	upgrade_time_in_secs = "%s"
	version_profile_id = "%s"
	override_version_profile = %s
	is_public = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the private cloud group
		resourcetype.ZPAPrivateCloudGroup,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		variable.PrivateCloudGroupCityCountry,
		variable.PrivateCloudGroupLatitude,
		variable.PrivateCloudGroupLongitude,
		variable.PrivateCloudGroupLocation,
		variable.PrivateCloudGroupUpgradeDayUpdate,
		variable.PrivateCloudGroupUpgradeTimeInSecsUpdate,
		variable.PrivateCloudGroupVersionProfileID,
		strconv.FormatBool(variable.PrivateCloudGroupOverrideVersionProfile),
		variable.PrivateCloudGroupIsPublic,

		// Data source type and name
		resourcetype.ZPAPrivateCloudGroup, resourceName,

		// Reference to the resource
		resourcetype.ZPAPrivateCloudGroup, resourceName,
	)
}
