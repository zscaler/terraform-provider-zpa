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

func TestAccResourcePolicyInspectionRuleV2_Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyInspectionRuleV2)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	// updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(20)

	// inspectionProfileTypeAndName, _, inspectionProfileGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionProfile)
	// inspectionProfileHCL := testAccCheckInspectionProfileConfigure(inspectionProfileTypeAndName, "tf-acc-test-"+inspectionProfileGeneratedName, variable.InspectionProfileDescription, variable.InspectionProfileParanoia)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyInspectionRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyInspectionRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyInspectionRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INSPECT"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},
			// Update test
			{
				Config: testAccCheckPolicyInspectionRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyInspectionRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INSPECT"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
				// ExpectNonEmptyPlan: true,
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

func testAccCheckPolicyInspectionRuleV2Destroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(apiClient.PolicySetControllerV2, "INSPECTION_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyInspectionRuleV2 {
			continue
		}

		rule, _, err := policysetcontrollerv2.GetPolicyRule(apiClient.PolicySetControllerV2, accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy inspection rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyInspectionRuleV2Exists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := policysetcontrollerv2.GetByPolicyType(apiClient.PolicySetControllerV2, "INSPECTION_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
		}
		_, _, err = policysetcontrollerv2.GetPolicyRule(apiClient.PolicySetControllerV2, resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyInspectionRuleV2Configure(resourceTypeAndName, rName, generatedName, desc string) string {
	return fmt.Sprintf(`

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getPolicyInspectionRuleV2HCL(rName, generatedName, desc),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyInspectionRuleV2HCL(rName, generatedName, desc string) string {
	return fmt.Sprintf(`

data "zpa_inspection_profile" "this" {
	name = "BD_AppProtection_Profile1"
}

resource "%s" "%s" {
	name          				= "%s"
	description   				= "%s"
	action              		= "INSPECT"
	zpn_inspection_profile_id 	= data.zpa_inspection_profile.this.id
	conditions {
		operator = "OR"
		operands {
		  object_type = "CLIENT_TYPE"
		  values      = ["zpn_client_type_exporter"]
		}
	}
}
`,
		// resource variables
		resourcetype.ZPAPolicyInspectionRuleV2,
		rName,
		generatedName,
		desc,
	)
}
