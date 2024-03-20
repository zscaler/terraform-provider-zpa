package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
)

func TestAccResourcePolicyTimeoutRuleV2Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyTimeOutRuleV2)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyTimeoutRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyTimeoutRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyTimeoutRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "RE_AUTH"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "reauth_idle_timeout", "10 Days"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "reauth_timeout", "10 Days"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyTimeoutRuleV2Configure(resourceTypeAndName, generatedName, updatedRName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyTimeoutRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedRName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "RE_AUTH"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "reauth_idle_timeout", "10 Days"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "reauth_timeout", "10 Days"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
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

func testAccCheckPolicyTimeoutRuleV2Destroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := apiClient.policysetcontrollerv2.GetByPolicyType("TIMEOUT_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource TIMEOUT_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyTimeOutRuleV2 {
			continue
		}

		rule, _, err := apiClient.policysetcontrollerv2.GetPolicyRule(accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy timeout rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyTimeoutRuleV2Exists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := apiClient.policysetcontrollerv2.GetByPolicyType("TIMEOUT_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource TIMEOUT_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontrollerv2.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyTimeoutRuleV2Configure(resourceTypeAndName, rName, generatedName, desc string) string {
	return fmt.Sprintf(`

// policy timeout rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getPolicyTimeoutRuleV2HCL(rName, generatedName, desc),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyTimeoutRuleV2HCL(rName, generatedName, desc string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action              = "RE_AUTH"
	reauth_idle_timeout = "10 Days"
	reauth_timeout      = "10 Days"
	conditions {
		operator = "OR"
		operands {
		  object_type = "CLIENT_TYPE"
		  values      = ["zpn_client_type_exporter"]
		}
	}
}
`,
		// resource variables
		resourcetype.ZPAPolicyTimeOutRuleV2,
		rName,
		generatedName,
		desc,
	)
}
