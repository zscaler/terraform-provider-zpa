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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
)

func TestAccResourceApplicationServer_Basic(t *testing.T) {
	var servers appservercontroller.ApplicationServer
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAApplicationServer)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationServerConfigure(resourceTypeAndName, initialName, variable.AppServerDescription, variable.AppServerAddress, variable.AppServerEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationServerExists(resourceTypeAndName, &servers),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AppServerDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "address", variable.AppServerAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.AppServerEnabled)),
				),
			},

			// Update test
			{
				Config: testAccCheckApplicationServerConfigure(resourceTypeAndName, updatedName, variable.AppServerDescription, variable.AppServerAddress, variable.AppServerEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationServerExists(resourceTypeAndName, &servers),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AppServerDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "address", variable.AppServerAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.AppServerEnabled)),
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

func testAccCheckApplicationServerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAApplicationServer {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		server, _, err := appservercontroller.Get(context.Background(), service, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if server != nil {
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
		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		receivedServer, _, err := appservercontroller.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*server = *receivedServer

		return nil
	}
}

func testAccCheckApplicationServerConfigure(resourceTypeAndName, generatedName, description, address string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`
	resource "%s" "%s" {
	name            = "%s"
	description     = "%s"
	address         = "%s"
	enabled         = "%s"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.ZPAApplicationServer,
		resourceName,
		generatedName,
		description,
		address,
		strconv.FormatBool(enabled),

		// Data source type and name
		resourcetype.ZPAApplicationServer, resourceName,

		// Reference to the resource
		resourcetype.ZPAApplicationServer, resourceName,
	)
}
