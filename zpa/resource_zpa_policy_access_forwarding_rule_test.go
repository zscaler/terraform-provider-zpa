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

func TestAccPolicyForwardingRuleBasic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyForwardingRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	updatedRName := acctest.RandomWithPrefix("tf-acc-test-updated") // New name for update test
	randDesc := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyForwardingRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyForwardingRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "BYPASS"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyForwardingRuleConfigure(resourceTypeAndName, generatedName, updatedRName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyForwardingRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedRName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "BYPASS"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "2"),
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

func testAccCheckPolicyForwardingRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := apiClient.policysetcontroller.GetByPolicyType("CLIENT_FORWARDING_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource CLIENT_FORWARDING_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyAccessRule {
			continue
		}

		rule, _, err := apiClient.policysetcontroller.GetPolicyRule(accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy forwarding rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyForwardingRuleExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := apiClient.policysetcontroller.GetByPolicyType("CLIENT_FORWARDING_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource CLIENT_FORWARDING_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontroller.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyForwardingRuleConfigure(resourceTypeAndName, rName, generatedName, desc string) string {
	return fmt.Sprintf(`

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getPolicyForwardingRuleHCL(rName, generatedName, desc),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyForwardingRuleHCL(rName, generatedName, desc string) string {
	return fmt.Sprintf(`

data "zpa_policy_type" "forwarding_policy" {
	policy_type = "CLIENT_FORWARDING_POLICY"
}

data "zpa_posture_profile" "crwd_zta_score_80" {
	name = "CrowdStrike_ZPA_ZTA_80 (zscalertwo.net)"
}

data "zpa_idp_controller" "bd_user_okta" {
    name = "BD_Okta_Users"
}

data "zpa_scim_groups" "contractors" {
	name     = "Contractors"
	idp_name = "BD_Okta_Users"
}

resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action              = "BYPASS"
	operator      		= "AND"
	policy_set_id 		= data.zpa_policy_type.forwarding_policy.id
	conditions {
		negated  = false
		operator = "OR"
		operands {
		  object_type = "POSTURE"
		  lhs         = data.zpa_posture_profile.crwd_zta_score_80.posture_udid
		  rhs         = false
		}
	  }
	  conditions {
		negated  = false
		operator = "OR"
		operands {
		  object_type = "SCIM_GROUP"
		  lhs         = data.zpa_idp_controller.bd_user_okta.id
		  rhs         = data.zpa_scim_groups.contractors.id
		  idp_id      = data.zpa_idp_controller.bd_user_okta.id
		}
	  }
}
`,
		// resource variables
		resourcetype.ZPAPolicyForwardingRule,
		rName,
		generatedName,
		desc,
	)
}
