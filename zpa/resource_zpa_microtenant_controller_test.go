package zpa

/*
import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
)

func TestAccResourceMicroTenant_Basic(t *testing.T) {
	var microTenant microtenants.MicroTenant
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAMicrotenant)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMicroTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMicrotenantConfigure(resourceTypeAndName, generatedName, variable.MicrotenantDescription, variable.MicrotenantCriteriaAttribute, variable.MicrotenantEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicrotenantExists(resourceTypeAndName, &microTenant),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.MicrotenantDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.MicrotenantEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "criteria_attribute_values.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckMicrotenantConfigure(resourceTypeAndName, generatedName, variable.MicrotenantDescription, variable.MicrotenantCriteriaAttribute, variable.MicrotenantEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicrotenantExists(resourceTypeAndName, &microTenant),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.MicrotenantDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.MicrotenantEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "criteria_attribute_values.#", "1"),
				),
			},
		},
	})
}

func testAccCheckMicroTenantDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPASegmentGroup {
			continue
		}

		group, _, err := apiClient.segmentgroup.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if group != nil {
			return fmt.Errorf("segment group with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMicrotenantExists(resource string, microtenant *microtenants.MicroTenant) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedMicrotenant, _, err := apiClient.microtenants.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*microtenant = *receivedMicrotenant

		return nil
	}
}

func testAccCheckMicrotenantConfigure(resourceTypeAndName, generatedName, description, criteria_attribute string, enabled bool) string {
	return fmt.Sprintf(`
// microtenant resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		MicroTenantHCL(generatedName, description, criteria_attribute, enabled),

		// data source variables
		resourcetype.ZPAMicrotenant,
		generatedName,
		resourceTypeAndName,
	)
}

func MicroTenantHCL(generatedName, description, criteria_attribute string, enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	enabled = %s
	criteria_attribute = "%s"
	criteria_attribute_values = ["securitygeek.io"]
}
`,
		// resource variables
		resourcetype.ZPAMicrotenant,
		generatedName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		criteria_attribute,
	)
}
*/
