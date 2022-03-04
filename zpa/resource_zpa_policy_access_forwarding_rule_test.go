package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePolicyForwardingRule(t *testing.T) {

	rName := acctest.RandString(10)
	rDesc := acctest.RandString(20)
	resourceName := "zpa_policy_forwarding_rule.policy_forwarding_test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZPAPolicyForwardingRuleConfigBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckPolicyForwardingRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "BYPASS"),
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

func testAccResourceZPAPolicyForwardingRuleConfigBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

data "zpa_policy_type" "forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

resource "zpa_policy_forwarding_rule" "policy_forwarding_test" {
	name                          = "%s"
	description                   = "%s"
	action                        = "BYPASS"
	rule_order                    = 1
	operator = "AND"
	policy_set_id = data.zpa_policy_type.forwarding_policy.id
}
	`, rName, rDesc)
}

func tesAccCheckPolicyForwardingRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Policy Forwarding Rule Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Policy Forwarding Rule ID is set")
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

func testAccCheckPolicyForwardingRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zpa_policy_access_rule" {
			continue
		}

		_, _, err := client.policysetrule.Get(rs.Primary.Attributes["policy_set_id"], rs.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("Policy Access Rule still exists")
		}

		return nil
	}
	return nil
}
