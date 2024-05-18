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

func TestAccResourcePolicyAccessRuleV2Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyAccessRuleV2)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	// updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(10)

	// appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	// appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, "tf-acc-test-"+segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyAccessRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyAccessRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc, segmentGroupHCL, segmentGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "4"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyAccessRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc, segmentGroupHCL, segmentGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "4"),
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

func testAccCheckPolicyAccessRuleV2Destroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := apiClient.policysetcontrollerv2.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyAccessRuleV2 {
			continue
		}

		rule, _, err := apiClient.policysetcontrollerv2.GetPolicyRule(accessPolicy.ID, rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("policy access rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyAccessRuleV2Exists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := apiClient.policysetcontrollerv2.GetByPolicyType("ACCESS_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
		}
		_, _, err = apiClient.policysetcontrollerv2.GetPolicyRule(resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyAccessRuleV2Configure(resourceTypeAndName, rName, generatedName, desc, segmentGroupHCL, segmentGroupTypeAndName string) string {
	return fmt.Sprintf(`

// app connector group resource

// segment group resource
%s

// policy access rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		// appConnectorGroupHCL,
		segmentGroupHCL,
		getPolicyAccessRuleV2HCL(rName, generatedName, desc, segmentGroupTypeAndName),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyAccessRuleV2HCL(rName, generatedName, desc, segmentGroupTypeAndName string) string {
	return fmt.Sprintf(`

data "zpa_idp_controller" "this" {
	name = "BD_Okta_Users"
}

data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_BD_Okta_Users"
    idp_name = "BD_Okta_Users"
}

data "zpa_saml_attribute" "group_user" {
    name = "GroupName_BD_Okta_Users"
    idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "a000" {
    name = "A000"
    idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "b000" {
    name = "B000"
    idp_name = "BD_Okta_Users"
}

resource "%s" "%s" {
	name          		= "%s"
	description   		= "%s"
	action        		= "ALLOW"
	operator      		= "AND"
	conditions {
		operator = "OR"
		operands {
			object_type = "APP_GROUP"
			values      = ["${%s.id}"]
		}
	}
	conditions {
		operator = "OR"
		operands {
		  object_type = "SCIM_GROUP"
		  entry_values {
			rhs = data.zpa_scim_groups.a000.id
			lhs = data.zpa_idp_controller.this.id // This is the IdP ID
		  }
		  entry_values {
			rhs = data.zpa_scim_groups.b000.id
			lhs = data.zpa_idp_controller.this.id // This is the IdP ID
		  }
		}
	}
	conditions {
		operator = "OR"
		operands {
		  object_type = "PLATFORM"
		  entry_values {
			rhs = "true"
			lhs = "linux"
		  }
		  entry_values {
			rhs = "true"
			lhs = "android"
		  }
		}
	}
	conditions {
		operator = "OR"
		operands {
		  object_type = "COUNTRY_CODE"
		  entry_values {
			lhs = "CA"
			rhs = "true"
		  }
		  entry_values {
			lhs = "US"
			rhs = "true"
		  }
		}
	}
	depends_on = [ %s ]
}
`,
		// resource variables
		resourcetype.ZPAPolicyAccessRuleV2,
		rName,
		generatedName,
		desc,
		// appConnectorGroupTypeAndName,
		segmentGroupTypeAndName,
		// appConnectorGroupTypeAndName,
		segmentGroupTypeAndName,
	)
}
