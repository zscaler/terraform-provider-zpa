package zpa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
)

func TestAccDataSourcePRAPrivilegedApproval_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRAApprovalController)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPRAPrivilegedApprovalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPRAPrivilegedApprovalConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email_ids.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "domain", resourceTypeAndName, "domain"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "status", resourceTypeAndName, "status"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "applications.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "working_hours.#", "1"),
				),
			},
		},
	})
}
