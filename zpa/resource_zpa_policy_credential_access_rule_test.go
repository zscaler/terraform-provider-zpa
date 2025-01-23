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

func TestAccResourcePolicyCredentialAccessRule_Basic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPolicyCredentialRule)
	rName := acctest.RandomWithPrefix("tf-acc-test")
	randDesc := acctest.RandString(20)
	rPassword := acctest.RandString(10)

	praCredentialTypeAndName, _, praCredentialGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAPRACredentialController)
	praCredentialHCL := testAccCheckPRACredentialControllerConfigure(praCredentialTypeAndName, "tf-acc-test-"+praCredentialGeneratedName, variable.PraConsoleDescription, rPassword)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyCredentialAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyCredentialAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, praCredentialHCL, praCredentialTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyCredentialAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INJECT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "credential.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "conditions.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckPolicyCredentialAccessRuleConfigure(resourceTypeAndName, generatedName, rName, randDesc, praCredentialHCL, praCredentialTypeAndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyCredentialAccessRuleExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", rName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", "INJECT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "credential.#", "1"),
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

func testAccCheckPolicyCredentialAccessRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(context.Background(), apiClient.Service, "CREDENTIAL_POLICY")
	if err != nil {
		return fmt.Errorf("failed fetching resource CREDENTIAL_POLICY. Received error: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPAPolicyCredentialRule {
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
			return fmt.Errorf("policy credential rule with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPolicyCredentialAccessRuleExists(resource string) resource.TestCheckFunc {
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

		resp, _, err := policysetcontrollerv2.GetByPolicyType(context.Background(), apiClient.Service, "CREDENTIAL_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource CREDENTIAL_POLICY. Recevied error: %s", err)
		}
		_, _, err = policysetcontrollerv2.GetPolicyRule(context.Background(), service, resp.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyCredentialAccessRuleConfigure(resourceTypeAndName, rName, generatedName, desc, praCredentialHCL, praCredentialTypeAndName string) string {
	return fmt.Sprintf(`

// pra credential
%s

// pra policy credential access rule
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		praCredentialHCL,
		getPolicyCredentialAccessRuleHCL(rName, generatedName, desc, praCredentialTypeAndName),

		// data source variables
		resourcetype.ZPAPolicyType,
		generatedName,
		resourceTypeAndName,
	)
}

func getPolicyCredentialAccessRuleHCL(rName, generatedName, desc, praCredentialTypeAndName string) string {
	return fmt.Sprintf(`

resource "zpa_application_segment_pra" "this" {
	name             = "tf-acc-test"
	description      = "tf-acc-test"
	enabled          = true
	health_reporting = "ON_ACCESS"
	bypass_type      = "NEVER"
	is_cname_enabled = true
	tcp_port_ranges  = ["3223", "3223", "3392", "3392"]
	domain_names     = ["ssh_pra3223.example.com", "rdp_pra3392.example.com"]
	segment_group_id = zpa_segment_group.this.id
	common_apps_dto {
		apps_config {
			domain               = "rdp_pra3392.example.com"
			application_protocol = "RDP"
			connection_security  = "ANY"
			application_port     = "3392"
			enabled              = true
			app_types            = ["SECURE_REMOTE_ACCESS"]
		}
		apps_config {
			domain               = "ssh_pra3223.example.com"
			application_protocol = "SSH"
			application_port     = "3223"
			enabled              = true
			app_types            = ["SECURE_REMOTE_ACCESS"]
		}
	}
}
	
resource "zpa_segment_group" "this" {
	name        = "tf-acc-test-%s"
	description = "tf-acc-test-%s"
	enabled     = true
}

data "zpa_application_segment_by_type" "rdp_pra3392" {
    application_type = "SECURE_REMOTE_ACCESS"
    name = "rdp_pra3392"
	depends_on = [zpa_application_segment_pra.this]
}

data "zpa_ba_certificate" "this1" {
	name = "pra01.bd-hashicorp.com"
}
  
resource "zpa_pra_portal_controller" "this" {
	name                      = "pra01.bd-hashicorp.com"
	description               = "pra01.bd-hashicorp.com"
	enabled                   = true
	domain                    = "pra01.bd-hashicorp.com"
	certificate_id            = data.zpa_ba_certificate.this1.id
	user_notification         = "Created with Terraform"
	user_notification_enabled = true
  }

resource "zpa_pra_console_controller" "rdp_pra" {
	name        = "rdp_console"
	description = "Created with Terraform"
	enabled     = true
	pra_application {
		id = data.zpa_application_segment_by_type.rdp_pra3392.id
	}
	pra_portals {
		id = [ zpa_pra_portal_controller.this.id ]
	}
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
	name          				= "%s"
	description   				= "%s"
	action              		= "INJECT_CREDENTIALS"
	credential {
		id = "${%s.id}"
	}
	conditions {
		operator = "OR"
		operands {
			object_type = "CONSOLE"
			values         = [zpa_pra_console_controller.rdp_pra.id]
			}
		}
	depends_on = [zpa_pra_console_controller.rdp_pra, zpa_pra_portal_controller.this]
}

`,
		// PRA Application Segment and Segment Group name generation
		generatedName,
		generatedName,

		// resource variables
		resourcetype.ZPAPolicyCredentialRule,
		rName,
		generatedName,
		desc,
		praCredentialTypeAndName,
	)
}
