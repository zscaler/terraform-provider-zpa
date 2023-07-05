package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func TestAccResourceCBIBannersBasic(t *testing.T) {
	var cbiBanner cbibannercontroller.CBIBannerController
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPACBIBannerController)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBIBannerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCBIBannerConfigure(resourceTypeAndName, generatedName, variable.PrimaryColor, variable.TextColor, variable.NotificationTitle, variable.NotificationText, variable.Banner, variable.Persist, variable.Logo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBIBannerExists(resourceTypeAndName, &cbiBanner),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
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
				Config: testAccCheckCBIBannerConfigure(resourceTypeAndName, generatedName, variable.PrimaryColor, variable.TextColor, variable.NotificationTitle, variable.NotificationText, variable.Banner, variable.Persist, variable.Logo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBIBannerExists(resourceTypeAndName, &cbiBanner),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "primary_color", variable.PrimaryColor),
					resource.TestCheckResourceAttr(resourceTypeAndName, "text_color", variable.TextColor),
					resource.TestCheckResourceAttr(resourceTypeAndName, "notification_title", variable.NotificationTitle),
					resource.TestCheckResourceAttr(resourceTypeAndName, "notification_text", variable.NotificationText),
					resource.TestCheckResourceAttr(resourceTypeAndName, "logo", variable.Logo),
					resource.TestCheckResourceAttr(resourceTypeAndName, "banner", strconv.FormatBool(variable.Banner)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "persist", strconv.FormatBool(variable.Persist)),
				),
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

		banner, _, err := apiClient.cbibannercontroller.Get(rs.Primary.ID)

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
		receivedBanner, _, err := apiClient.cbibannercontroller.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*banner = *receivedBanner

		return nil
	}
}

func testAccCheckCBIBannerConfigure(resourceTypeAndName, generatedName, primaryColor, textColor, notificationTitle, NotificationText string, banner, persist bool, logo string) string {
	return fmt.Sprintf(`
// cbi banner resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		CBIBannerResourceHCL(generatedName, primaryColor, textColor, notificationTitle, NotificationText, banner, persist, logo),

		// data source variables
		resourcetype.ZPACBIBannerController,
		generatedName,
		resourceTypeAndName,
	)
}

func CBIBannerResourceHCL(generatedName, primaryColor, textColor, notificationTitle, NotificationText string, banner, persist bool, logo string) string {
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
`,
		// resource variables
		resourcetype.ZPACBIBannerController,
		generatedName,
		generatedName,
		primaryColor,
		textColor,
		notificationTitle,
		NotificationText,
		strconv.FormatBool(banner),
		strconv.FormatBool(persist),
		logo,
	)
}
