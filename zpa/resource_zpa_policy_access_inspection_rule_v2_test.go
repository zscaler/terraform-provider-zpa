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
)

func TestAccResourcePolicyInspectionRuleV2Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyInspectionRuleV2)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(20)

	inspectionProfileTypeAndName, _, inspectionProfileGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAInspectionProfile)
	inspectionProfileHCL := testAccCheckInspectionProfileConfigure(inspectionProfileTypeAndName, "tf-acc-test-"+inspectionProfileGeneratedName, variable.InspectionProfileDescription, variable.InspectionProfileParanoia)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyInspectionRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyInspectionRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc, inspectionProfileHCL, inspectionProfileTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyInspectionRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INSPECT"),
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
				Config: testAccCheckPolicyInspectionRuleV2Configure(resourceTypeAndName, generatedName, updatedRName, randDesc, inspectionProfileHCL, inspectionProfileTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyInspectionRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedRName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INSPECT"),
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

func testAccCheckPolicyInspectionRuleV2Destroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := apiClient.policysetcontrollerv2.GetByPolicyType("INSPECTION_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyInspectionRuleV2 {
			continue
		}

		rule, _, err := apiClient.policysetcontrollerv2.GetPolicyRule(accessPolicy.ID, rs.Primary.ID)

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
		resp, _, err := apiClient.policysetcontrollerv2.GetByPolicyType("INSPECTION_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource INSPECTION_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontrollerv2.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyInspectionRuleV2Configure(resourceTypeAndName, rName, generatedName, desc, inspectionProfileHCL, inspectionProfileTypeAndName string) string {
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
		getPolicyInspectionRuleV2HCL(rName, generatedName, desc, inspectionProfileTypeAndName),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyInspectionRuleV2HCL(rName, generatedName, desc, inspectionProfileTypeAndName string) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	name          				= "%s"
	description   				= "%s"
	action              		= "INSPECT"
	zpn_inspection_profile_id 	= "${%s.id}"
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
		inspectionProfileTypeAndName,
	)
}
