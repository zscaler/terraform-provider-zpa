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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

func TestAccResourceC2CIPRanges_Basic(t *testing.T) {
	var c2cIPRanges c2c_ip_ranges.IPRanges
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAC2CIPRanges)

	initialName := "tf-acc-test-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckC2CIPRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckC2CIPRangesConfigure(resourceTypeAndName, initialName, variable.C2CIPRangesDescription, variable.C2CIPRangesEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckC2CIPRangesExists(resourceTypeAndName, &c2cIPRanges),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.C2CIPRangesDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.C2CIPRangesEnabled)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "location_hint", "Created_via_Terraform"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_range_begin", "192.168.1.1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_range_end", "192.168.1.254"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "location", "San Jose, CA, USA"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "sccm_flag", "true"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "latitude_in_db", "37.33874"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "longitude_in_db", "-121.8852525"),
				),
			},

			// Update test
			{
				Config: testAccCheckC2CIPRangesConfigureUpdate(resourceTypeAndName, initialName, variable.C2CIPRangesDescriptionUpdate, variable.C2CIPRangesEnabledUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckC2CIPRangesExists(resourceTypeAndName, &c2cIPRanges),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.C2CIPRangesDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", strconv.FormatBool(variable.C2CIPRangesEnabledUpdate)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_range_begin", variable.C2CIPRangesIPRangeBeginUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_range_end", variable.C2CIPRangesIPRangeEndUpdate),
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

func testAccCheckC2CIPRangesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAC2CIPRanges {
			continue
		}

		service := apiClient.Service

		ipRanges, _, err := c2c_ip_ranges.Get(context.Background(), service, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if ipRanges != nil {
			return fmt.Errorf("C2C IP ranges with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckC2CIPRangesExists(resource string, ipRanges *c2c_ip_ranges.IPRanges) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		receivedIPRanges, _, err := c2c_ip_ranges.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*ipRanges = *receivedIPRanges

		return nil
	}
}

func testAccCheckC2CIPRangesConfigure(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = %s
	location_hint = "Created_via_Terraform"
	ip_range_begin = "192.168.1.1"
	ip_range_end = "192.168.1.254"
	location = "San Jose, CA, USA"
	sccm_flag = true
	country_code = "US"
	latitude_in_db = "37.33874"
	longitude_in_db = "-121.8852525"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the C2C IP ranges
		resourcetype.ZPAC2CIPRanges,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),

		// Data source type and name
		resourcetype.ZPAC2CIPRanges, resourceName,

		// Reference to the resource
		resourcetype.ZPAC2CIPRanges, resourceName,
	)
}

func testAccCheckC2CIPRangesConfigureUpdate(resourceTypeAndName, generatedName, description string, enabled bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	enabled = %s
	location_hint = "Created_via_Terraform"
	ip_range_begin = "%s"
	ip_range_end = "%s"
	location = "San Jose, CA, USA"
	sccm_flag = true
	country_code = "US"
	latitude_in_db = "37.33874"
	longitude_in_db = "-121.8852525"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the C2C IP ranges
		resourcetype.ZPAC2CIPRanges,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(enabled),
		variable.C2CIPRangesIPRangeBeginUpdate,
		variable.C2CIPRangesIPRangeEndUpdate,

		// Data source type and name
		resourcetype.ZPAC2CIPRanges, resourceName,

		// Reference to the resource
		resourcetype.ZPAC2CIPRanges, resourceName,
	)
}
