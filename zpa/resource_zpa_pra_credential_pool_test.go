package zpa

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

func TestAccTesourcePRACredentialPool_Basic(t *testing.T) {
	var credentialPool pracredentialpool.CredentialPool
	rPassword := acctest.RandString(10)

	credentialPoolTypeAndName, _, credentialPoolGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRACredentialPool)

	credentialControllerTypeAndName, _, credentialControllerGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRACredentialController)
	credentialControllerHCL := testAccCheckPRACredentialControllerConfigure(credentialControllerTypeAndName, "tf-acc-test-"+credentialControllerGeneratedName, variable.CredentialDescription, rPassword)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRACredentialPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRACredentialPoolConfigure(credentialPoolTypeAndName, credentialPoolGeneratedName, credentialPoolGeneratedName, credentialControllerHCL, credentialControllerTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRACredentialPoolExists(credentialPoolTypeAndName, &credentialPool),
					resource.TestCheckResourceAttr(credentialPoolTypeAndName, "name", "tf-acc-test-"+credentialPoolGeneratedName),
					resource.TestCheckResourceAttr(credentialPoolTypeAndName, "credential_type", variable.CredentialType),
					resource.TestCheckResourceAttr(credentialPoolTypeAndName, "credentials.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckPRACredentialPoolConfigure(credentialPoolTypeAndName, credentialPoolGeneratedName, credentialPoolGeneratedName, credentialControllerHCL, credentialControllerTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPRACredentialPoolExists(credentialPoolTypeAndName, &credentialPool),
					resource.TestCheckResourceAttr(credentialPoolTypeAndName, "name", "tf-acc-test-"+credentialPoolGeneratedName),
					resource.TestCheckResourceAttr(credentialPoolTypeAndName, "credential_type", variable.CredentialType),
					resource.TestCheckResourceAttr(credentialPoolTypeAndName, "credentials.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      credentialPoolTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPRACredentialPoolDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPRACredentialPool {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		pool, _, err := pracredentialpool.Get(context.Background(), service, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if pool != nil {
			return fmt.Errorf("pra credential pool with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPRACredentialPoolExists(resource string, pool *pracredentialpool.CredentialPool) resource.TestCheckFunc {
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

		receivedPool, _, err := pracredentialpool.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*pool = *receivedPool

		return nil
	}
}

func testAccCheckPRACredentialPoolConfigure(resourceTypeAndName, generatedName, name, credentialControllerHCL, credentialControllerTypeAndName string) string {
	return fmt.Sprintf(`

// credential controller resource
%s

// credential pool resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// credential controller block
		credentialControllerHCL,

		// credential pool block (fix: use only generated name for `name`)
		getPRACredentialPoolResourceHCL(generatedName, generatedName, credentialControllerTypeAndName),

		// data block
		resourcetype.ZPAPRACredentialPool,
		generatedName,
		resourceTypeAndName,
	)
}

func getPRACredentialPoolResourceHCL(generatedName, name, credentialControllerTypeAndName string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	credential_type = "%s"
	credentials {
		id = ["${%s.id}"]
	}
	depends_on = [ %s ]
}
`,
		// resource variables
		resourcetype.ZPAPRACredentialPool,
		generatedName,
		name,
		variable.CredentialType,
		credentialControllerTypeAndName,
		credentialControllerTypeAndName,
	)
}
