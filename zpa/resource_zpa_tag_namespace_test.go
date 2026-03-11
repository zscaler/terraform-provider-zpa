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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_namespace"
)

func TestAccResourceTagNamespace_Basic(t *testing.T) {
	var ns tag_namespace.Namespace
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPATagNamespace)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTagNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTagNamespaceConfigure(resourceTypeAndName, initialName, variable.TagNamespaceDescription, variable.TagNamespaceEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagNamespaceExists(resourceTypeAndName, &ns),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.TagNamespaceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.TagNamespaceEnabled)),
				),
			},
			// Update test
			{
				Config: testAccCheckTagNamespaceConfigure(resourceTypeAndName, updatedName, variable.TagNamespaceDescriptionUpdate, variable.TagNamespaceEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagNamespaceExists(resourceTypeAndName, &ns),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.TagNamespaceDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.TagNamespaceEnabledUpdate)),
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

func testAccCheckTagNamespaceDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPATagNamespace {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		ns, _, err := tag_namespace.Get(context.Background(), service, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if ns != nil {
			return fmt.Errorf("tag namespace with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTagNamespaceExists(resource string, ns *tag_namespace.Namespace) resource.TestCheckFunc {
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

		receivedNS, _, err := tag_namespace.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*ns = *receivedNS

		return nil
	}
}

func testAccCheckTagNamespaceConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]

	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "%s"
	enabled     = %s
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		resourcetype.ZPATagNamespace,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),

		resourcetype.ZPATagNamespace, resourceName,
		resourcetype.ZPATagNamespace, resourceName,
	)
}
