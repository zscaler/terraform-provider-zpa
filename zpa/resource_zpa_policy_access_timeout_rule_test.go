package zpa

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePolicyTimeoutRule(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	resourceName := "zpa_policy_timeout_rule.timeout_rule_test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyTimeoutRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZPAPolicyTimeoutRuleConfigBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckPolicyTimeoutRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "RE_AUTH"),
					resource.TestCheckResourceAttr(resourceName, "rule_order", "1"),
					resource.TestCheckResourceAttr(resourceName, "operator", "AND"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceZPAPolicyTimeoutRuleConfigBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

resource "zpa_policy_timeout_rule" "timeout_rule_test" {
	name                          = "%s"
	description                   = "%s"
	action                        = "RE_AUTH"
	rule_order                    = 1
	reauth_idle_timeout 		  = "600"
	reauth_timeout 				  = "172800"
	operator 					  = "AND"
	policy_set_id 				  = data.zpa_policy_type.timeout_policy.id
}
	`, rName, rDesc)
}

func tesAccCheckPolicyTimeoutRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Policy Timeout Rule Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Policy Timeout Rule ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		resp, _, err := client.policysetrule.Get(rs.Primary.Attributes["policy_set_id"], rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		if resp.Name != rs.Primary.Attributes["policy_set_id"] {
			return fmt.Errorf("name Not found in created attributes")
		}
		if resp.Description != rs.Primary.Attributes["id"] {
			return fmt.Errorf("description Not found in created attributes")
		}
		return nil
	}
}

func testAccCheckPolicyTimeoutRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_policy_timeout_rule" {
			continue
		}

		_, _, err := client.policysetrule.Get(rs.Primary.Attributes["policy_set_id"], rs.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("Policy Timeout Rule still exists")
		}

		return nil
	}
	return nil
}
*/
