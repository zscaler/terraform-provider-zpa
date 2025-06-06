package zpa

/*
// TODO: Testing disabled as QA environments have limited region access
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

func TestAccResourceCBIExternalProfile_Basic(t *testing.T) {
	var cbiIsolationProfile cbiprofilecontroller.IsolationProfile
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPACBIExternalIsolationProfile)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCBIExternalProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCBIExternalProfileConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBIExternalProfileExists(resourceTypeAndName, &cbiIsolationProfile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_experience.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "security_controls.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckCBIExternalProfileConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBIExternalProfileExists(resourceTypeAndName, &cbiIsolationProfile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_experience.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "security_controls.#", "1"),
				),
			},
			// Import test by ID
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return cbiIsolationProfile.ID, nil
				},
			},
		},
	})
}

func testAccCheckCBIExternalProfileDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPACBIExternalIsolationProfile {
			continue
		}

		profile, _, err := cbiprofilecontroller.Get(apiClient.CBIProfileController, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if profile != nil {
			return fmt.Errorf("cbi external profile with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCBIExternalProfileExists(resource string, profile *cbiprofilecontroller.IsolationProfile) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedProfile, _, err := cbiprofilecontroller.Get(apiClient.CBIProfileController, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*profile = *receivedProfile

		return nil
	}
}

func testAccCheckCBIExternalProfileConfigure(resourceTypeAndName, generatedName string) string {
	return fmt.Sprintf(`
// cbi external profile resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		CBIExternalProfileResourceHCL(generatedName),

		// data source variables
		resourcetype.ZPACBIExternalIsolationProfile,
		generatedName,
		resourceTypeAndName,
	)
}

func CBIExternalProfileResourceHCL(generatedName string) string {
	return fmt.Sprintf(`

data "zpa_cloud_browser_isolation_banner" "this" {
	name = "Default"
	}

data "zpa_cloud_browser_isolation_region" "singapore" {
name = "Singapore"
}

data "zpa_cloud_browser_isolation_region" "frankfurt" {
name = "Frankfurt"
}

data "zpa_cloud_browser_isolation_certificate" "this" {
  name = "Zscaler Root Certificate"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
    banner_id = data.zpa_cloud_browser_isolation_banner.this.id
    region_ids = [data.zpa_cloud_browser_isolation_region.singapore.id, data.zpa_cloud_browser_isolation_region.frankfurt.id]
    certificate_ids = [data.zpa_cloud_browser_isolation_certificate.this.id]

    user_experience {
		session_persistence    = true
		browser_in_browser 	   = true
		persist_isolation_bar  = true
		translate              = true
	}
	  security_controls {
		copy_paste 			= "all"
		upload_download 	= "upstream"
		document_viewer     = true
		local_render        = true
		allow_printing      = true
		restrict_keystrokes = true
		flattened_pdf       = true
		deep_link {
			enabled           	= true
			applications      	= ["test_app1", "test_app2"]
		}
		watermark {
			enabled          = true
			show_user_id     = true
			show_timestamp   = true
			show_message     = true
			message          = "CBI Via ZPA Terraform Provider"
		}
	}
}
`,
		// resource variables
		resourcetype.ZPACBIExternalIsolationProfile,
		generatedName,
		generatedName,
		generatedName,
	)
}
*/
