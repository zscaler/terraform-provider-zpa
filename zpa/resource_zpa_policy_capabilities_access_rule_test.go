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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func TestAccResourcePolicyCapabilitiesAccessRule_Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyCapabilitiesRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	randDesc := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyCapabilitiesAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyCapabilitiesAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyCapabilitiesAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "CHECK_CAPABILITIES"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "privileged_capabilities.#", "1"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyCapabilitiesAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyCapabilitiesAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "CHECK_CAPABILITIES"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "privileged_capabilities.#", "1"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPolicyCapabilitiesAccessRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(context.Background(), apiClient.Service, "CAPABILITIES_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource CAPABILITIES_POLICY. Received error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyCapabilitiesRule {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		rule, _, err := policysetcontrollerv2.GetPolicyRule(context.Background(), service, accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy credential rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyCapabilitiesAccessRuleExists(resource string) resource.TestCheckFunc {
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

		resp, _, err := policysetcontrollerv2.GetByPolicyType(context.Background(), apiClient.Service, "CAPABILITIES_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource CAPABILITIES_POLICY. Recevied error: %s", err)
		}
		_, _, err = policysetcontrollerv2.GetPolicyRule(context.Background(), service, resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyCapabilitiesAccessRuleConfigure(resourceTypeAndName, rName, generatedName, desc string) string {
	return fmt.Sprintf(`

// pra policy capabilities access rule
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getPolicyCapabilitiesAccessRuleHCL(rName, generatedName, desc),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyCapabilitiesAccessRuleHCL(rName, generatedName, desc string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name          				= "%s"
	description   				= "%s"
	action              		= "CHECK_CAPABILITIES"
	privileged_capabilities {
		clipboard_copy        = true
		clipboard_paste       = true
		file_upload           = true
		file_download         = true
		inspect_file_upload   = true
		inspect_file_download = true
	}
}
`,
		// resource variables
		resourcetype.ZPAPolicyCapabilitiesRule,
		rName,
		generatedName,
		desc,
	)
}
