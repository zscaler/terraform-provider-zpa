package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

func TestAccResourcePolicyRedictionRule_Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyRedirectionRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(20)

	serviceEdgeGroupTypeAndName, _, serviceEdgeGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAServiceEdgeGroup)
	serviceEdgeGroupHCL := testAccCheckServiceEdgeGroupConfigure(serviceEdgeGroupTypeAndName, "tf-acc-test-"+serviceEdgeGroupGeneratedName, variable.ServiceEdgeDescription, variable.ServiceEdgeEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyRedictionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyRedictionRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, serviceEdgeGroupHCL, serviceEdgeGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyRedictionRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "REDIRECT_PREFERRED"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "service_edge_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyRedictionRuleConfigure(resourceTypeAndName, generatedName, updatedRName, randDesc, serviceEdgeGroupHCL, serviceEdgeGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyRedictionRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedRName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "REDIRECT_PREFERRED"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "service_edge_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
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

func testAccCheckPolicyRedictionRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := policysetcontroller.GetByPolicyType(apiClient.PolicySetController, "REDIRECTION_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource REDIRECTION_POLICY. Received error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyRedirectionRule {
			continue
		}

		rule, _, err := policysetcontroller.GetPolicyRule(apiClient.PolicySetController, accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("redirection access policy with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyRedictionRuleExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := policysetcontroller.GetByPolicyType(apiClient.PolicySetController, "REDIRECTION_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource REDIRECTION_POLICY. Recevied error: %s", err)
		}
		_, _, err = policysetcontroller.GetPolicyRule(apiClient.PolicySetController, resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyRedictionRuleConfigure(resourceTypeAndName, rName, generatedName, desc, serviceEdgeGroupHCL, serviceEdgeGroupTypeAndName string) string {
	return fmt.Sprintf(`

// app connector group resource
%s

// redirection policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		serviceEdgeGroupHCL,
		getPolicyRedictionRuleHCL(rName, generatedName, desc, serviceEdgeGroupTypeAndName),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyRedictionRuleHCL(rName, generatedName, desc, serviceEdgeGroupTypeAndName string) string {
	return fmt.Sprintf(`

data "zpa_policy_type" "this" {
	policy_type = "REDIRECTION_POLICY"
}

resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action        		= "REDIRECT_PREFERRED"
	operator      		= "AND"
	policy_set_id 		= data.zpa_policy_type.this.id
	service_edge_groups {
		id = ["${%s.id}"]
	}
	conditions {
		operator = "OR"
			operands {
				object_type = "CLIENT_TYPE"
				lhs         = "id"
				rhs         = "zpn_client_type_branch_connector"
		}
			operands {
				object_type = "CLIENT_TYPE"
				lhs         = "id"
				rhs         = "zpn_client_type_edge_connector"
		}
	  }
	depends_on = [ %s ]
}
`,
		// resource variables
		resourcetype.ZPAPolicyRedirectionRule,
		rName,
		generatedName,
		desc,
		serviceEdgeGroupTypeAndName,
		serviceEdgeGroupTypeAndName,
	)
}
