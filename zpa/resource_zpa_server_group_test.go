package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
)

func TestAccResourceServerGroupBasic(t *testing.T) {
	var serverGroup servergroup.ServerGroup
	serverGroupTypeAndName, _, serverGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServerGroup)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerGroupConfigure(serverGroupTypeAndName, serverGroupGeneratedName, serverGroupGeneratedName, serverGroupGeneratedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerGroupExists(serverGroupTypeAndName, &serverGroup),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "name", "tf-acc-test-"+serverGroupGeneratedName),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "description", "tf-acc-test-"+serverGroupGeneratedName),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "enabled", strconv.FormatBool(variable.ServerGroupEnabled)),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "dynamic_discovery", strconv.FormatBool(variable.ServerGroupDynamicDiscovery)),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "app_connector_groups.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckServerGroupConfigure(serverGroupTypeAndName, serverGroupGeneratedName, serverGroupGeneratedName, serverGroupGeneratedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, variable.ServerGroupEnabled, variable.ServerGroupDynamicDiscovery),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerGroupExists(serverGroupTypeAndName, &serverGroup),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "name", "tf-acc-test-"+serverGroupGeneratedName),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "description", "tf-acc-test-"+serverGroupGeneratedName),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "enabled", strconv.FormatBool(variable.ServerGroupEnabled)),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "dynamic_discovery", strconv.FormatBool(variable.ServerGroupDynamicDiscovery)),
					resource.TestCheckResourceAttr(serverGroupTypeAndName, "app_connector_groups.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      serverGroupTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccCheckServerGroupConfigure(resourceTypeAndName, generatedName, name, description, appConnectorGroupHCL, appConnectorGroupTypeAndName string, enabled, dynDiscovery bool) string {
	return fmt.Sprintf(`

// app connector group resource
%s

// server group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		appConnectorGroupHCL,
		getServerGroupResourceHCL(generatedName, name, description, appConnectorGroupTypeAndName, enabled, dynDiscovery),

		// data source variables
		resourcetype.ZPAServerGroup,
		generatedName,
		resourceTypeAndName,
	)
}

func getServerGroupResourceHCL(generatedName, name, description, appConnectorGroupTypeAndName string, enabled, dynDiscovery bool) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled = "%s"
	dynamic_discovery = "%s"
	app_connector_groups {
		id = ["${%s.id}"]
	}
	depends_on = [ %s ]
}
`,
		// resource variables
		resourcetype.ZPAServerGroup,
		generatedName,
		name,
		description,
		strconv.FormatBool(enabled),
		strconv.FormatBool(dynDiscovery),
		appConnectorGroupTypeAndName,
		appConnectorGroupTypeAndName,
	)
}
