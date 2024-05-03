package zpa

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccZPAResourcePolicyAccessRuleReorder_basic(t *testing.T) {
	randName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyAccessReorderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAccessRuleReorderConfig(randName),
				Check: resource.ComposeTestCheckFunc(
					// Wrapper function for the CheckPolicyType test
					func(s *terraform.State) error {
						success := t.Run("CheckPolicyType", func(_ *testing.T) {
							resource.TestCheckResourceAttr("zpa_policy_access_rule_reorder.this", "policy_type", "ACCESS_POLICY")(s)
						})
						if !success {
							return fmt.Errorf("CheckPolicyType failed")
						}
						return nil
					},
					// Wrapper function for the CheckFirstRuleID test
					func(s *terraform.State) error {
						success := t.Run("CheckFirstRuleID", func(_ *testing.T) {
							resource.TestCheckResourceAttrPair("zpa_policy_access_rule_reorder.this", "rules.0.id", "zpa_policy_access_rule.rule1", "id")(s)
						})
						if !success {
							return fmt.Errorf("CheckFirstRuleID failed")
						}
						return nil
					},
					// ... add similar wrappers for the other checks ...
					func(s *terraform.State) error {
						success := t.Run("CheckZscalerDeceptionOrder", func(_ *testing.T) {
							testCheckZscalerDeceptionOrder(s)
						})
						if !success {
							return fmt.Errorf("CheckZscalerDeceptionOrder failed")
						}
						return nil
					},
					func(s *terraform.State) error {
						success := t.Run("CheckRuleOrderValidity", func(_ *testing.T) {
							testCheckRuleOrderValidity(s)
						})
						if !success {
							return fmt.Errorf("CheckRuleOrderValidity failed")
						}
						return nil
					},
				),
			},
		},
	})
}

func testCheckZscalerDeceptionOrder(s *terraform.State) error {
	rs, ok := s.RootModule().Resources["zpa_policy_access_rule_reorder.this"]
	if !ok {
		return fmt.Errorf("Not found: zpa_policy_access_rule_reorder.this")
	}
	for attrKey, attrValue := range rs.Primary.Attributes {
		if attrKey == "name" && attrValue == "Zscaler Deception" {
			if rs.Primary.Attributes["order"] != "2" {
				return fmt.Errorf("Zscaler Deception rule order is not 2")
			}
		}
	}
	return nil
}

func testCheckRuleOrderValidity(s *terraform.State) error {
	rs, ok := s.RootModule().Resources["zpa_policy_access_rule_reorder.this"]
	if !ok {
		return fmt.Errorf("Not found: zpa_policy_access_rule_reorder.this")
	}

	totalRulesStr, exists := rs.Primary.Attributes["rules.#"]
	if !exists {
		return fmt.Errorf("Failed to get the count of rules from the state")
	}

	totalRules, err := strconv.Atoi(totalRulesStr)
	if err != nil {
		return fmt.Errorf("Failed to convert total rules count to integer: %s", err)
	}

	encounteredOrders := make(map[int]bool)

	for attrKey, attrValue := range rs.Primary.Attributes {
		if strings.HasSuffix(attrKey, ".order") {
			order, err := strconv.Atoi(attrValue)
			if err != nil {
				return fmt.Errorf("Failed to convert order to integer: %s", err)
			}

			// Use Case 2: Rule Order <= 0 is not allowed
			if order <= 0 {
				return fmt.Errorf("Rule order <= 0 found")
			}

			// Use Case 5: Rule order numbers cannot be duplicated
			if _, exists := encounteredOrders[order]; exists {
				return fmt.Errorf("Duplicated rule order found: %d", order)
			}
			encounteredOrders[order] = true
		}
	}

	// Use Case 3: Gaps in between rules is not allowed
	for i := 1; i <= totalRules; i++ {
		if _, exists := encounteredOrders[i]; !exists {
			return fmt.Errorf("Gap found in rule orders at position: %d", i)
		}
	}

	// Use Case 4: Rule order number cannot be higher than the total number of rules
	if len(encounteredOrders) != totalRules {
		return fmt.Errorf("Rule order number exceeds total number of rules")
	}

	return nil
}

func testAccCheckPolicyAccessReorderDestroy(s *terraform.State) error {
	// Here you can add checks to verify if the resource has been destroyed if applicable.
	return nil
}

func testAccPolicyAccessRuleReorderConfig(randName string) string {
	return fmt.Sprintf(`

resource "zpa_policy_access_rule" "rule1" {
    name          = "%s-rule1"
    description   = "%s-desc1"
    action        = "ALLOW"
    operator      = "AND"

    lifecycle {
        create_before_destroy = true
    }
}

resource "zpa_policy_access_rule" "rule2" {
    name          = "%s-rule2"
    description   = "%s-desc2"
    action        = "ALLOW"
    operator      = "AND"
    depends_on    = [zpa_policy_access_rule.rule1]

    lifecycle {
        create_before_destroy = true
    }
}

resource "zpa_policy_access_rule_reorder" "this" {
    policy_type   = "ACCESS_POLICY"

    dynamic "rules" {
        for_each = [zpa_policy_access_rule.rule1.id, zpa_policy_access_rule.rule2.id]
        content {
            id    = rules.value
            order = tostring(rules.key + 1)  // Since the key starts at 0, we add 1 to get the order
        }
    }

    lifecycle {
        create_before_destroy = true
    }
}
`, randName, randName, randName, randName)
}
