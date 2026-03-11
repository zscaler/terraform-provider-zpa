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
	tag_key_controller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_key"
)

func TestAccResourceTagKey_Basic(t *testing.T) {
	var tagKey tag_key_controller.TagKey
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPATagKey)
	nsResourceTypeAndName, _, nsGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPATagNamespace)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName
	nsName := "tf-acc-test-" + nsGeneratedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTagKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTagKeyConfigure(nsResourceTypeAndName, nsName, resourceTypeAndName, initialName, variable.TagKeyDescription, variable.TagKeyEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagKeyExists(resourceTypeAndName, &tagKey),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.TagKeyDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.TagKeyEnabled)),
				),
			},
			// Update test
			{
				Config: testAccCheckTagKeyConfigure(nsResourceTypeAndName, nsName, resourceTypeAndName, updatedName, variable.TagKeyDescriptionUpdate, variable.TagKeyEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagKeyExists(resourceTypeAndName, &tagKey),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.TagKeyDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.TagKeyEnabledUpdate)),
				),
			},
		},
	})
}

func testAccCheckTagKeyDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPATagKey {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		namespaceID := rs.Primary.Attributes["namespace_id"]
		tagKey, _, err := tag_key_controller.Get(context.Background(), service, namespaceID, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if tagKey != nil {
			return fmt.Errorf("tag key with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTagKeyExists(resource string, tagKey *tag_key_controller.TagKey) resource.TestCheckFunc {
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

		namespaceID := rs.Primary.Attributes["namespace_id"]
		receivedTagKey, _, err := tag_key_controller.Get(context.Background(), service, namespaceID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*tagKey = *receivedTagKey

		return nil
	}
}

func testAccCheckTagKeyConfigure(nsResourceTypeAndName, nsName, resourceTypeAndName, generatedName, description string, enabled bool) string {
	nsResourceName := strings.Split(nsResourceTypeAndName, ".")[1]
	resourceName := strings.Split(resourceTypeAndName, ".")[1]

	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "Namespace for tag key test"
	enabled     = true
}

resource "%s" "%s" {
	name         = "%s"
	description  = "%s"
	enabled      = %s
	namespace_id = "${%s.%s.id}"
}

data "%s" "%s" {
  id           = "${%s.%s.id}"
  namespace_id = "${%s.%s.id}"
}
`,
		resourcetype.ZPATagNamespace,
		nsResourceName,
		nsName,

		resourcetype.ZPATagKey,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		resourcetype.ZPATagNamespace, nsResourceName,

		resourcetype.ZPATagKey, resourceName,
		resourcetype.ZPATagKey, resourceName,
		resourcetype.ZPATagNamespace, nsResourceName,
	)
}
