package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccPRACredentialPool_basic(t *testing.T) {
	var credentialPool pracredentialpool.CredentialPool
	rName := sdkacctest.RandString(8)
	rPassword := sdkacctest.RandString(10)
	resourceName := "zpa_pra_credential_pool.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPRACredentialPoolDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPRACredentialPoolConfig(rName, rPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPRACredentialPoolExists(zClient, resourceName, &credentialPool),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "credential_type", "USERNAME_PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "credentials.#", "1"),
				),
			},
			// Update test
			{
				Config: testAccPRACredentialPoolConfig(rName, rPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPRACredentialPoolExists(zClient, resourceName, &credentialPool),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "credential_type", "USERNAME_PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "credentials.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPRACredentialPoolExists(zClient *client.Client, resourceName string, pool *pracredentialpool.CredentialPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("pra credential pool not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("pra credential pool ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		received, _, err := pracredentialpool.Get(ctx, service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch pra credential pool %s: %w", rs.Primary.ID, err)
		}

		*pool = *received
		return nil
	}
}

func testAccCheckPRACredentialPoolDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_pra_credential_pool" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			pool, _, err := pracredentialpool.Get(ctx, service, rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("id %s already exists", rs.Primary.ID)
			}

			if pool != nil {
				return fmt.Errorf("pra credential pool with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccPRACredentialPoolConfig(rName, rPassword string) string {
	return fmt.Sprintf(`
resource "zpa_pra_credential_controller" "test" {
	name = "tf-acc-test-%s"
	description = "Credential Controller Test"
	credential_type = "USERNAME_PASSWORD"
	user_domain = "acme.com"
	username = "jcarrow"
	password = "%s"
}

resource "zpa_pra_credential_pool" "test" {
	name = "tf-acc-test-%s"
	credential_type = "USERNAME_PASSWORD"
	credentials {
		id = [zpa_pra_credential_controller.test.id]
	}
	depends_on = [zpa_pra_credential_controller.test]
}

data "zpa_pra_credential_pool" "test" {
	id = zpa_pra_credential_pool.test.id
}
`, rName, rPassword, rName)
}
