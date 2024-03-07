package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
)

func TestAccDataSourcePRAPrivilegedApproval_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAApprovalController)
	domainName := "pra_" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAPrivilegedApprovalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAPrivilegedApprovalConfigure(resourceTypeAndName, domainName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "email_ids.0", resourceTypeAndName, "email_ids.0"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "domain", resourceTypeAndName, "domain"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "status", resourceTypeAndName, "status"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "applications.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "working_hours.#", "1"),
				),
			},
		},
	})
}
