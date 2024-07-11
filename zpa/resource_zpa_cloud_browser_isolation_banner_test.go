package zpa

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func TestAccResourceCBIBanners_Basic(t *testing.T) {
	var cbiBanner cbibannercontroller.CBIBannerController
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPACBIBannerController)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBIBannerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCBIBannerConfigure(resourceTypeAndName, initialName, variable.PrimaryColor, variable.TextColor, variable.NotificationTitle, variable.NotificationText, variable.Banner, variable.Persist, variable.Logo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBIBannerExists(resourceTypeAndName, &cbiBanner),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "primary_color", variable.PrimaryColor),
					resource.TestCheckResourceAttr(resourceTypeAndName, "text_color", variable.TextColor),
					resource.TestCheckResourceAttr(resourceTypeAndName, "notification_title", variable.NotificationTitle),
					resource.TestCheckResourceAttr(resourceTypeAndName, "notification_text", variable.NotificationText),
					resource.TestCheckResourceAttr(resourceTypeAndName, "logo", variable.Logo),
					resource.TestCheckResourceAttr(resourceTypeAndName, "banner", strconv.FormatBool(variable.Banner)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "persist", strconv.FormatBool(variable.Persist)),
				),
			},

			// Update test
			{
				Config: testAccCheckCBIBannerConfigure(resourceTypeAndName, updatedName, variable.PrimaryColor, variable.TextColor, variable.NotificationTitle, variable.NotificationText, variable.Banner, variable.Persist, variable.Logo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBIBannerExists(resourceTypeAndName, &cbiBanner),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "primary_color", variable.PrimaryColor),
					resource.TestCheckResourceAttr(resourceTypeAndName, "text_color", variable.TextColor),
					resource.TestCheckResourceAttr(resourceTypeAndName, "notification_title", variable.NotificationTitle),
					resource.TestCheckResourceAttr(resourceTypeAndName, "notification_text", variable.NotificationText),
					resource.TestCheckResourceAttr(resourceTypeAndName, "logo", variable.Logo),
					resource.TestCheckResourceAttr(resourceTypeAndName, "banner", strconv.FormatBool(variable.Banner)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "persist", strconv.FormatBool(variable.Persist)),
				),
			},
			// Import test by ID
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return cbiBanner.ID, nil
				},
			},
		},
	})
}

func testAccCheckCBIBannerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPACBIBannerController {
			continue
		}

		banner, _, err := cbibannercontroller.Get(apiClient.CBIBannerController, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if banner != nil {
			return fmt.Errorf("cbi banner with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCBIBannerExists(resource string, banner *cbibannercontroller.CBIBannerController) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedBanner, _, err := cbibannercontroller.Get(apiClient.CBIBannerController, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*banner = *receivedBanner

		return nil
	}
}

func testAccCheckCBIBannerConfigure(resourceTypeAndName, generatedName, primaryColor, textColor, notificationTitle, NotificationText string, banner, persist bool, logo string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	primary_color = "%s"
	text_color = "%s"
	notification_title = "%s"
	notification_text = "%s"
	banner = "%s"
	persist = "%s"
	logo = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the cbi banner
		resourcetype.ZPACBIBannerController,
		resourceName,
		generatedName,
		primaryColor,
		textColor,
		notificationTitle,
		NotificationText,
		strconv.FormatBool(banner),
		strconv.FormatBool(persist),
		logo,

		// Data source type and name
		resourcetype.ZPACBIBannerController, resourceName,

		// Reference to the resource
		resourcetype.ZPACBIBannerController, resourceName,
	)
}
