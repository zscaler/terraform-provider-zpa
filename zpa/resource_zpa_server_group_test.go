package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/servergroup"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccResourceServerGroupBasic(t *testing.T) {
	var serverGroup servergroup.ServerGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerGroupConfigure(resourceTypeAndName, generatedName, variable.ServerGroupDescription, variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerGroupExists(resourceTypeAndName, &serverGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ServerGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServerGroupEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServerGroupDynamicDiscovery)),
				),
			},

			// Update test
			{
				Config: testAccCheckServerGroupConfigure(resourceTypeAndName, generatedName, variable.ServerGroupDescription, variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerGroupExists(resourceTypeAndName, &serverGroup),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ServerGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.ServerGroupEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dynamic_discovery", strconv.FormatBool(variable.ServerGroupDynamicDiscovery)),
				),
			},
		},
	})
}

func testAccCheckServerGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAServerGroup {
			continue
		}

		rule, _, err := apiClient.servergroup.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("server group with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckServerGroupExists(resource string, rule *servergroup.ServerGroup) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedGroup, _, err := apiClient.servergroup.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedGroup

		return nil
	}
}

func testAccCheckServerGroupConfigure(resourceTypeAndName, generatedName, description string, enabled, dynamic_discovery bool) string {
	return fmt.Sprintf(`
// server group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		ServerGroupResourceHCL(generatedName, description, enabled, dynamic_discovery),

		// data source variables
		resourcetype.ZPAServerGroup,
		generatedName,
		resourceTypeAndName,
	)
}

func ServerGroupResourceHCL(generatedName, description string, enabled, dynamic_discovery bool) string {
	return fmt.Sprintf(`
resource "zpa_app_connector_group" "testAcc" {
	name                          = "testAcc"
	description                   = "testAcc"
	enabled                       = true
	country_code                  = "US"
	latitude                      = "37.3382082"
	longitude                     = "-121.8863286"
	location                      = "San Jose, CA, USA"
	upgrade_day                   = "SUNDAY"
	upgrade_time_in_secs          = "66600"
	override_version_profile      = true
	version_profile_id            = 0
	dns_query_type                = "IPV4"
}

resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = "%s"
	dynamic_discovery = "%s"
	app_connector_groups {
		id = [zpa_app_connector_group.testAcc.id]
	}
	depends_on = [zpa_app_connector_group.testAcc]
}

`,
		// resource variables
		resourcetype.ZPAServerGroup,
		generatedName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		strconv.FormatBool(dynamic_discovery),
	)
}
