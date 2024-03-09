package zpa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/pracredential"
)

func TestAccResourcePRACredentialControllerBasic(t *testing.T) {
	var praCredential pracredential.Credential
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRACredentialController)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName
	rPassword := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRACredentialControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRACredentialControllerConfigure(resourceTypeAndName, initialName, variable.CredentialDescription, rPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRACredentialControllerExists(resourceTypeAndName, &praCredential),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CredentialDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "credential_type", "USERNAME_PASSWORD"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_domain", "acme.com"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", "jcarrow"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", rPassword),
				),
			},

			// Update test
			{
				Config: testAccCheckPRACredentialControllerConfigure(resourceTypeAndName, updatedName, variable.CredentialDescriptionUpdate, rPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRACredentialControllerExists(resourceTypeAndName, &praCredential),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CredentialDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "credential_type", "USERNAME_PASSWORD"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "user_domain", "acme.com"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", "jcarrow"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", rPassword),
				),
			},
			// Import test with ImportStateVerifyIgnore for password and user_domain
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"passphrase",
					"private_key",
					"user_domain",
				},
			},
		},
	})
}

func testAccCheckPRACredentialControllerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPRACredentialController {
			continue
		}

		group, _, err := apiClient.pracredential.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if group != nil {
			return fmt.Errorf("pra credential with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPRACredentialControllerExists(resource string, credential *pracredential.Credential) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedCredential, _, err := apiClient.pracredential.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*credential = *receivedCredential

		return nil
	}
}

func testAccCheckPRACredentialControllerConfigure(resourceTypeAndName, generatedName, description, rPassword string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name = "%s"
	description = "%s"
	credential_type = "USERNAME_PASSWORD"
    user_domain = "acme.com"
    username = "jcarrow"
    password = "%s"
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// Resource type and name for the certificate
		resourcetype.ZPAPRACredentialController,
		resourceName,
		generatedName,
		description,
		rPassword,

		// Data source type and name
		resourcetype.ZPAPRACredentialController, resourceName,

		// Reference to the resource
		resourcetype.ZPAPRACredentialController, resourceName,
	)
}
