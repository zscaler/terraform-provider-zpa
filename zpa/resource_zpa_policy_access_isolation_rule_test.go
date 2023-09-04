package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/testing/method"
)

func TestAccPolicyIsolationRuleBasic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyIsolationRule)
	rName := acctest.RandomWithPrefix("tf-acc-test-")
	randDesc := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyIsolationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyIsolationRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyIsolationRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ISOLATE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyIsolationRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyIsolationRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ISOLATE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},
		},
	})
}

func testAccCheckPolicyIsolationRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := apiClient.policysetcontroller.GetByPolicyType("ISOLATION_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource ISOLATION_POLICY. Received error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyIsolationRule {
			continue
		}

		rule, _, err := apiClient.policysetcontroller.GetPolicyRule(accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy isolation rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyIsolationRuleExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := apiClient.policysetcontroller.GetByPolicyType("ISOLATION_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource ISOLATION_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontroller.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyIsolationRuleConfigure(resourceTypeAndName, rName, generatedName, desc string) string {
	return fmt.Sprintf(`

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getPolicyIsolationRuleHCL(rName, generatedName, desc),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyIsolationRuleHCL(rName, generatedName, desc string) string {
	return fmt.Sprintf(`

data "zpa_policy_type" "isolation_policy" {
	policy_type = "ISOLATION_POLICY"
}

data "zpa_isolation_profile" "bd_sa_profile1" {
	name = "BD_SA_Profile1"
}
resource "%s" "%s" {
	name          				= "%s"
	description   				= "%s"
	action              		= "ISOLATE"
	operator      				= "AND"
	policy_set_id 				= data.zpa_policy_type.isolation_policy.id
	zpn_isolation_profile_id 	= data.zpa_isolation_profile.bd_sa_profile1.id
	conditions {
		negated  = false
		operator = "OR"
		operands {
			object_type = "CLIENT_TYPE"
			lhs         = "id"
			rhs         = "zpn_client_type_exporter"
			}
		}
}

`,
		// resource variables
		resourcetype.ZPAPolicyIsolationRule,
		rName,
		generatedName,
		desc,
	)
}
