package zpa

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/policysetcontroller"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccPolicyAccessRuleBasic(t *testing.T) {
	var rules policysetcontroller.PolicySet
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyAccessRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, variable.AccessRuleDescription, variable.AccessRuleAction),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", fmt.Sprintf(rName)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceState),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, variable.AccessRuleDescription, variable.AccessRuleAction),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", fmt.Sprintf(rName)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceState),
				),
			},
		},
	})
}

func testAccCheckPolicyAccessRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyAccessRule {
			continue
		}

		rule, _, err := apiClient.policysetcontroller.GetPolicyRule(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %d already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy access rule with id %d exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyAccessRuleExists(resource string, rule *policysetcontroller.PolicySet) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedRule, err := apiClient.policysetcontroller.GetPolicyRule(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, description, action, state string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action        		= "ALLOW"
	rule_order    		= 4
	operator      		= "AND"
	policy_set_id 		= data.zpa_policy_type.access_policy.id
`,
		// resource variables
		resourcetype.ZPAPolicyAccessRule,
		generatedName,
		generatedName,
		description,
		action,
		state,

		// data source variables
		// resourcetype.FirewallFilteringRules,
		// generatedName,
		// resourceTypeAndName,
	)
}
*/
