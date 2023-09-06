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

func TestAccPolicyInspectionRuleBasic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyInspectionRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	randDesc := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyInspectionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyInspectionRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyInspectionRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INSPECT"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
				ExpectNonEmptyPlan: true,
			},

			// Update test
			{
				Config: testAccCheckPolicyInspectionRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyInspectionRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INSPECT"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckPolicyInspectionRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := apiClient.policysetcontroller.GetByPolicyType("INSPECTION_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyInspectionRule {
			continue
		}

		rule, _, err := apiClient.policysetcontroller.GetPolicyRule(accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy inspection rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyInspectionRuleExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := apiClient.policysetcontroller.GetByPolicyType("INSPECTION_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontroller.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyInspectionRuleConfigure(resourceTypeAndName, rName, generatedName, desc string) string {
	return fmt.Sprintf(`

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getPolicyInspectionRuleHCL(rName, generatedName, desc),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyInspectionRuleHCL(rName, generatedName, desc string) string {
	return fmt.Sprintf(`

data "zpa_policy_type" "inspection_policy" {
	policy_type = "INSPECTION_POLICY"
}

resource "%s" "%s" {
	name          				= "%s"
	description   				= "%s"
	action              		= "INSPECT"
	operator      				= "AND"
	policy_set_id 				= data.zpa_policy_type.inspection_policy.id
	zpn_inspection_profile_id 	= zpa_inspection_profile.this.id
	conditions {
		negated  = false
		operator = "OR"
		operands {
			object_type = "CLIENT_TYPE"
			lhs         = "id"
			rhs         = "zpn_client_type_exporter"
			}
		}
	depends_on = [zpa_inspection_profile.this]
}

data "zpa_inspection_predefined_controls" "this" {
	name = "Failed to parse request body"
	version    = "OWASP_CRS/3.3.0"
  }

  data "zpa_inspection_all_predefined_controls" "default_predefined_controls" {
	version    = "OWASP_CRS/3.3.0"
	group_name = "Preprocessors"
  }

  resource "zpa_inspection_profile" "this" {
	name                        = "tf-acc-test"
	description                 = "tf-acc-test"
	paranoia_level              = "2"
	dynamic "predefined_controls" {
	  for_each = data.zpa_inspection_all_predefined_controls.default_predefined_controls.list
	  content {
		id           = predefined_controls.value.id
		action       = predefined_controls.value.action == "" ? predefined_controls.value.default_action : predefined_controls.value.action
		action_value = predefined_controls.value.action_value
	  }
	}
	predefined_controls {
	  id     = data.zpa_inspection_predefined_controls.this.id
	  action = "BLOCK"
	}
}
`,
		// resource variables
		resourcetype.ZPAPolicyInspectionRule,
		rName,
		generatedName,
		desc,
	)
}
