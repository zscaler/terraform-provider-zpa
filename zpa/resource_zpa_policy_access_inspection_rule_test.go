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

func TestAccResourcePolicyInspectionRule_Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyInspectionRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(20)

	inspectionProfileTypeAndName, _, inspectionProfileGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionProfile)
	inspectionProfileHCL := testAccCheckInspectionProfileConfigure(inspectionProfileTypeAndName, "tf-acc-test-"+inspectionProfileGeneratedName, variable.InspectionProfileDescription, variable.InspectionProfileParanoia)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyInspectionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyInspectionRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, inspectionProfileHCL, inspectionProfileTypeAndName),
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
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update test
			{
				Config: testAccCheckPolicyInspectionRuleConfigure(resourceTypeAndName, generatedName, updatedRName, randDesc, inspectionProfileHCL, inspectionProfileTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyInspectionRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedRName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INSPECT"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
				ExpectNonEmptyPlan: true,
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

func testAccCheckPolicyInspectionRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := policysetcontroller.GetByPolicyType(apiClient.PolicySetController, "INSPECTION_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyInspectionRule {
			continue
		}

		rule, _, err := policysetcontroller.GetPolicyRule(apiClient.PolicySetController, accessPolicy.ID, rs.Primary.ID)

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
		resp, _, err := policysetcontroller.GetByPolicyType(apiClient.PolicySetController, "INSPECTION_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
		}
		_, _, err = policysetcontroller.GetPolicyRule(apiClient.PolicySetController, resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyInspectionRuleConfigure(resourceTypeAndName, rName, generatedName, desc, inspectionProfileHCL, inspectionProfileTypeAndName string) string {
	return fmt.Sprintf(`

// Inspection profile resource
%s

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		inspectionProfileHCL,
		getPolicyInspectionRuleHCL(rName, generatedName, desc, inspectionProfileTypeAndName),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyInspectionRuleHCL(rName, generatedName, desc, inspectionProfileTypeAndName string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name          				= "%s"
	description   				= "%s"
	action              		= "INSPECT"
	operator      				= "AND"
	zpn_inspection_profile_id 	= "${%s.id}"
	conditions {
		operator = "OR"
		operands {
			object_type = "CLIENT_TYPE"
			lhs         = "id"
			rhs         = "zpn_client_type_exporter"
			}
		}
	}
`,
		// resource variables
		resourcetype.ZPAPolicyInspectionRule,
		rName,
		generatedName,
		desc,
		inspectionProfileTypeAndName,
	)
}
