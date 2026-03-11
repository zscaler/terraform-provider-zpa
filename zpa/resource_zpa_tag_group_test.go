package zpa

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_group"
)

func TestAccResourceTagGroup_Basic(t *testing.T) {
	var tg tag_group.TagGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPATagGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTagGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTagGroupConfigure(resourceTypeAndName, initialName, variable.TagGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagGroupExists(resourceTypeAndName, &tg),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.TagGroupDescription),
				),
			},
			// Update test
			{
				Config: testAccCheckTagGroupConfigure(resourceTypeAndName, updatedName, variable.TagGroupDescriptionUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagGroupExists(resourceTypeAndName, &tg),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.TagGroupDescriptionUpdate),
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

func testAccCheckTagGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPATagGroup {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		tg, _, err := tag_group.Get(context.Background(), service, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if tg != nil {
			return fmt.Errorf("tag group with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTagGroupExists(resource string, tg *tag_group.TagGroup) resource.TestCheckFunc {
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

		receivedTG, _, err := tag_group.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*tg = *receivedTG

		return nil
	}
}

func testAccCheckTagGroupConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]

	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		resourcetype.ZPATagGroup,
		resourceName,
		generatedName,
		description,

		resourcetype.ZPATagGroup, resourceName,
		resourcetype.ZPATagGroup, resourceName,
	)
}
