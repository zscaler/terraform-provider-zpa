package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
)

func TestAccDataSourcePRACredentialPool_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRACredentialPool)

	rPassword := acctest.RandString(10)

	credentialControllerTypeAndName, _, credentialControllerGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRACredentialController)
	credentialControllerHCL := testAccCheckPRACredentialControllerConfigure(credentialControllerTypeAndName, credentialControllerGeneratedName, variable.CredentialDescription, rPassword)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRACredentialPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRACredentialPoolConfigure(resourceTypeAndName, generatedName, generatedName, credentialControllerHCL, credentialControllerTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "credential_type", resourceTypeAndName, "credential_type"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "credentials.#", "1"),
				),
			},
		},
	})
}
