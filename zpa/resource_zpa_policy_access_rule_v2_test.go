package zpa

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v4/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func TestAccResourcePolicyAccessRuleV2_Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyAccessRuleV2)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	// updatedRName := acctest.RandomWithPrefix("tf-updated") // New name for update test
	randDesc := acctest.RandString(10)

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	segmentGroupTypeAndName, _, segmentGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPASegmentGroup)
	segmentGroupHCL := testAccCheckSegmentGroupConfigure(segmentGroupTypeAndName, "tf-acc-test-"+segmentGroupGeneratedName, variable.SegmentGroupDescription, variable.SegmentGroupEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyAccessRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyAccessRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc, segmentGroupHCL, segmentGroupTypeAndName, appConnectorGroupHCL, appConnectorGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "6"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyAccessRuleV2Configure(resourceTypeAndName, generatedName, rName, randDesc, segmentGroupHCL, segmentGroupTypeAndName, appConnectorGroupHCL, appConnectorGroupTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyAccessRuleV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "6"),
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
	accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(context.Background(), apiClient.Service, "ACCESS_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyAccessRuleV2 {
			continue
		}

		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		rule, _, err := policysetcontrollerv2.GetPolicyRule(context.Background(), service, accessPolicy.ID, rs.Primary.ID)

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
		microTenantID := rs.Primary.Attributes["microtenant_id"]
		service := apiClient.Service
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		resp, _, err := policysetcontrollerv2.GetByPolicyType(context.Background(), apiClient.Service, "ACCESS_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource ACCESS_POLICY. Recevied error: %s", err)
		}
		_, _, err = policysetcontrollerv2.GetPolicyRule(context.Background(), service, resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyAccessRuleV2Configure(resourceTypeAndName, rName, generatedName, desc, segmentGroupHCL, segmentGroupTypeAndName, appConnectorGroupHCL, appConnectorGroupTypeAndName string) string {
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
		segmentGroupHCL,
		appConnectorGroupHCL,
		getPolicyAccessRuleV2HCL(rName, generatedName, desc, segmentGroupTypeAndName, appConnectorGroupTypeAndName),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyAccessRuleV2HCL(rName, generatedName, desc, segmentGroupTypeAndName, appConnectorGroupTypeAndName string) string {
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
	app_connector_groups {
		id = ["${%s.id}"]
	}
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
	conditions {
		operator = "OR"
		operands {
		  object_type = "RISK_FACTOR_TYPE"
		  entry_values {
			lhs = "ZIA"
			rhs = "UNKNOWN"
		  }
		  entry_values {
			lhs = "ZIA"
			rhs = "LOW"
		  }
		  entry_values {
			lhs = "ZIA"
			rhs = "MEDIUM"
		  }
		}
	}
	conditions {
		operator = "OR"
		operands {
			object_type = "CHROME_ENTERPRISE"
			entry_values {
				lhs = "managed"
				rhs = "true"
			}
		}
	}
	depends_on = [ %s, %s ]
}
`,
		// resource variables
		resourcetype.ZPAPolicyAccessRuleV2,
		rName,
		generatedName,
		desc,
		appConnectorGroupTypeAndName,
		segmentGroupTypeAndName,
		appConnectorGroupTypeAndName,
		segmentGroupTypeAndName,
	)
}
