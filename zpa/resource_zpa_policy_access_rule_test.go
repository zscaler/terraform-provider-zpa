package zpa

<<<<<<< HEAD
=======
/*
>>>>>>> master
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
<<<<<<< HEAD
=======
	"github.com/willguibr/terraform-provider-zpa/gozscaler/policysetcontroller"
>>>>>>> master
	"github.com/willguibr/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/willguibr/terraform-provider-zpa/zpa/common/testing/variable"
)

func TestAccPolicyAccessRuleBasic(t *testing.T) {
<<<<<<< HEAD
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyAccessRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	randDesc := acctest.RandString(20)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)
=======
	var rules policysetcontroller.PolicySet
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyAccessRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
>>>>>>> master

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
<<<<<<< HEAD
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, appConnectorGroupHCL, appConnectorGroupTypeAndName, segmentGroupHCL, segmentGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
=======
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, variable.AccessRuleDescription, variable.AccessRuleAction),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", fmt.Sprintf(rName)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceState),
>>>>>>> master
				),
			},

			// Update test
			{
<<<<<<< HEAD
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, appConnectorGroupHCL, appConnectorGroupTypeAndName, segmentGroupHCL, segmentGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
=======
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, variable.AccessRuleDescription, variable.AccessRuleAction),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", fmt.Sprintf(rName)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceState),
>>>>>>> master
				),
			},
		},
	})
}

func testAccCheckPolicyAccessRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
<<<<<<< HEAD
	accessPolicy, _, err := apiClient.policysetcontroller.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
	}
=======

>>>>>>> master
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyAccessRule {
			continue
		}

<<<<<<< HEAD
		rule, _, err := apiClient.policysetcontroller.GetPolicyRule(accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy access rule with id %s exists and wasn't destroyed", rs.Primary.ID)
=======
		rule, _, err := apiClient.policysetcontroller.GetPolicyRule(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %d already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy access rule with id %d exists and wasn't destroyed", rs.Primary.ID)
>>>>>>> master
		}
	}

	return nil
}

<<<<<<< HEAD
func testAccCheckPolicyAccessRuleExists(resource string) resource.TestCheckFunc {
=======
func testAccCheckPolicyAccessRuleExists(resource string, rule *policysetcontroller.PolicySet) resource.TestCheckFunc {
>>>>>>> master
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
<<<<<<< HEAD
		resp, _, err := apiClient.policysetcontroller.GetByPolicyType("ACCESS_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontroller.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
=======
		receivedRule, err := apiClient.policysetcontroller.GetPolicyRule(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

>>>>>>> master
		return nil
	}
}

<<<<<<< HEAD
func testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, rName, generatedName, desc, appConnectorGroupHCL, appConnectorGroupTypeAndName, segmentGroupHCL, segmentGroupTypeAndName string) string {
	return fmt.Sprintf(`

// app connector group resource
%s

// segment group resource
%s

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		appConnectorGroupHCL,
		segmentGroupHCL,
		getPolicyAccessRuleHCL(rName, generatedName, desc, appConnectorGroupTypeAndName, segmentGroupTypeAndName),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyAccessRuleHCL(rName, generatedName, desc, appConnectorGroupTypeAndName, segmentGroupTypeAndName string) string {
	return fmt.Sprintf(`

data "zpa_policy_type" "access_policy" {
	policy_type = "ACCESS_POLICY"
}

=======
func testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, description, action, state string) string {
	return fmt.Sprintf(`

>>>>>>> master
resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action        		= "ALLOW"
<<<<<<< HEAD
	operator      		= "AND"
	policy_set_id 		= data.zpa_policy_type.access_policy.id
	app_connector_groups {
		id = ["${%s.id}"]
	}
	conditions {
		negated  = false
		operator = "OR"
		operands {
		  object_type = "APP_GROUP"
		  lhs         = "id"
		  rhs         = "${%s.id}"
		}
	  }
	depends_on = [ %s, %s ]
}
`,
		// resource variables
		resourcetype.ZPAPolicyAccessRule,
		rName,
		generatedName,
		desc,
		appConnectorGroupTypeAndName,
		segmentGroupTypeAndName,
		appConnectorGroupTypeAndName,
		segmentGroupTypeAndName,
	)
}

/*
func testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, rName, generatedName, desc, appConnectorGroupHCL, appConnectorGroupTypeAndName string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action        		= "ALLOW"
	operator      		= "AND"
	policy_set_id 		= data.zpa_policy_type.access_policy.id
}
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}
`,
		// resource variables
		resourcetype.ZPAPolicyAccessRule,
		rName,
		generatedName,
		desc,
=======
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
>>>>>>> master
	)
}
*/
