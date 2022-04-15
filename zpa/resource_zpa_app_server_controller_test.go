package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/appservercontroller"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccResourceApplicationServerBasic(t *testing.T) {
	var servers appservercontroller.ApplicationServer
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationServer)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationServerConfigure(resourceTypeAndName, generatedName, variable.AppServerDescription, variable.AppServerAddress, variable.AppServerEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationServerExists(resourceTypeAndName, &servers),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AppServerDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "address", variable.AppServerAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.AppServerEnabled)),
				),
			},

			// Update test
			{
				Config: testAccCheckApplicationServerConfigure(resourceTypeAndName, generatedName, variable.AppServerDescription, variable.AppServerAddress, variable.AppServerEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationServerExists(resourceTypeAndName, &servers),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AppServerDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "address", variable.AppServerAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.AppServerEnabled)),
				),
			},
		},
	})
}

func testAccCheckApplicationServerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAApplicationServer {
			continue
		}

		rule, _, err := apiClient.appservercontroller.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("application server with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckApplicationServerExists(resource string, server *appservercontroller.ApplicationServer) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedServer, _, err := apiClient.appservercontroller.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*server = *receivedServer

		return nil
	}
}

func testAccCheckApplicationServerConfigure(resourceTypeAndName, generatedName, description, address string, enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name            = "%s"
	description     = "%s"
	address         = "%s"
	enabled         = "%s"
}

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		resourcetype.ZPAApplicationServer,
		generatedName,
		generatedName,
		description,
		address,
		strconv.FormatBool(enabled),

		// data source variables
		resourcetype.ZPAApplicationServer,
		generatedName,
		resourceTypeAndName,
	)
}
