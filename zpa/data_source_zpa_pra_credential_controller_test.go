package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
)

func TestAccDataSourcePRACredentialController_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRACredentialController)
	rPassword := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRACredentialControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRACredentialControllerConfigure(resourceTypeAndName, generatedName, variable.CredentialDescription, rPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "credential_type", resourceTypeAndName, "credential_type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "user_domain", resourceTypeAndName, "user_domain"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "username", resourceTypeAndName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "password", resourceTypeAndName, rPassword),
				),
			},
		},
	})
}
