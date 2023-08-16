package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccZPAResourcePolicyAccessRuleReorder_basic(t *testing.T) {
	randName := acctest.RandString(10) // This will generate a random string of 10 characters
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyAccessReorderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAccessRuleReorderConfig(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zpa_policy_access_rule_reorder.this", "policy_type", "ACCESS_POLICY"),
					resource.TestCheckResourceAttrPair("zpa_policy_access_rule_reorder.this", "rules.0.id", "zpa_policy_access_rule.rule1", "id"),
					resource.TestCheckResourceAttr("zpa_policy_access_rule_reorder.this", "rules.0.order", "1"),
					resource.TestCheckResourceAttrPair("zpa_policy_access_rule_reorder.this", "rules.1.id", "zpa_policy_access_rule.rule2", "id"),
					resource.TestCheckResourceAttr("zpa_policy_access_rule_reorder.this", "rules.1.order", "2"),
				),
			},
		},
	})

}

func testAccCheckPolicyAccessReorderDestroy(s *terraform.State) error {
	// Here you can add checks to verify if the resource has been destroyed if applicable.
	return nil
}

func testAccPolicyAccessRuleReorderConfig(randName string) string {
	return fmt.Sprintf(`
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

resource "zpa_policy_access_rule" "rule1" {
    name          = "%s-rule1"
    description   = "%s-desc1"
    action        = "ALLOW"
    operator      = "AND"
    policy_set_id = data.zpa_policy_type.access_policy.id

    lifecycle {
        create_before_destroy = true
    }
}

resource "zpa_policy_access_rule" "rule2" {
    name          = "%s-rule2"
    description   = "%s-desc2"
    action        = "ALLOW"
    operator      = "AND"
    policy_set_id = data.zpa_policy_type.access_policy.id
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
