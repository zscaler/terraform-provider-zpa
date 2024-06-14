package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontrollerv2"
)

func TestAccResourcePolicyForwardingRuleV2Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyForwardingRuleV2)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	// updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyForwardingRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyForwardingRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyForwardingRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "BYPASS"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyForwardingRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyForwardingRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "BYPASS"),
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

func testAccCheckPolicyForwardingRuleV2Destroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(apiClient.PolicySetControllerV2, "CLIENT_FORWARDING_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource CLIENT_FORWARDING_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyForwardingRuleV2 {
			continue
		}

		rule, _, err := policysetcontrollerv2.GetPolicyRule(apiClient.PolicySetControllerV2, accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy forwarding rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyForwardingRuleV2Exists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := policysetcontrollerv2.GetByPolicyType(apiClient.PolicySetControllerV2, "CLIENT_FORWARDING_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource CLIENT_FORWARDING_POLICY. Recevied error: %s", err)
		}
		_, _, err = policysetcontrollerv2.GetPolicyRule(apiClient.PolicySetControllerV2, resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyForwardingRuleV2Configure(resourceTypeAndName, rName, generatedName, desc string) string {
	return fmt.Sprintf(`

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getPolicyForwardingRuleV2HCL(rName, generatedName, desc),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyForwardingRuleV2HCL(rName, generatedName, desc string) string {
	return fmt.Sprintf(`

data "zpa_posture_profile" "crwd_zta_score_80" {
	name = "CrowdStrike_ZPA_ZTA_80"
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
	conditions {
		operator = "OR"
		operands {
		  object_type = "POSTURE"
		  entry_values {
			lhs = data.zpa_posture_profile.crwd_zta_score_80.posture_udid
			rhs = false
		  }
		}
	  }
	conditions {
	operator = "OR"
	operands {
		object_type = "SCIM_GROUP"
		entry_values {
		lhs = data.zpa_idp_controller.bd_user_okta.id
		rhs = data.zpa_scim_groups.contractors.id
		}
	}
	}
}
`,
		// resource variables
		resourcetype.ZPAPolicyForwardingRuleV2,
		rName,
		generatedName,
		desc,
	)
}
