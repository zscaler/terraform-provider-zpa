package zpa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v2/zpa/common/testing/variable"
)

func TestAccPolicyAccessRuleBasic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyAccessRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	randDesc := acctest.RandString(20)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, appConnectorGroupHCL, appConnectorGroupTypeAndName, segmentGroupHCL, segmentGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, appConnectorGroupHCL, appConnectorGroupTypeAndName, segmentGroupHCL, segmentGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},
		},
	})
}

func testAccCheckPolicyAccessRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := apiClient.policysetcontroller.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
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
			return fmt.Errorf("policy access rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyAccessRuleExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := apiClient.policysetcontroller.GetByPolicyType("ACCESS_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontroller.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

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

data "zpa_scim_attribute_header" "givenName" {
    name = "name.givenName"
    idp_name = "BD_Okta_Users"
}

data "zpa_scim_attribute_header" "familyName" {
    name = "name.familyName"
    idp_name = "BD_Okta_Users"
}

data "zpa_scim_attribute_header" "username" {
    name = "userName"
    idp_name = "BD_Okta_Users"
}

resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action        		= "ALLOW"
	rule_order    		= 1
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
		operands {
			object_type = "SCIM"
			lhs =  data.zpa_scim_attribute_header.givenName.id
			rhs = "Charles"
			idp_id = data.zpa_scim_attribute_header.givenName.idp_id
		  }
		  operands {
			object_type = "SCIM"
			lhs =  data.zpa_scim_attribute_header.familyName.id
			rhs = "Keenan"
			idp_id = data.zpa_scim_attribute_header.familyName.idp_id
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
