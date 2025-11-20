package datasources_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
)

func TestAccPRACredentialPoolDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandString(8)
	rPassword := sdkacctest.RandString(10)
	resourceName := "zpa_pra_credential_pool.test"
	dataSourceName := "data.zpa_pra_credential_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPRACredentialPoolDataSourceConfig(rName, rPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "credential_type", resourceName, "credential_type"),
					resource.TestCheckResourceAttr(dataSourceName, "credentials.#", "1"),
				),
			},
		},
	})
}

func testAccPRACredentialPoolDataSourceConfig(rName, rPassword string) string {
	return fmt.Sprintf(`
resource "zpa_pra_credential_controller" "test" {
	name = "%s"
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
