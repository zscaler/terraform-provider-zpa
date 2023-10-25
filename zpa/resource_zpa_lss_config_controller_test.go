package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/v3/zpa/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/lssconfigcontroller"
)

func TestAccResourceLSSConfigControllerBasic(t *testing.T) {
	var lssConfig lssconfigcontroller.LSSConfig
	lssControllerTypeAndName, _, lssControllerGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPALSSController)
	rPort := acctest.RandIntRange(1000, 9999)
	rIP, _ := acctest.RandIpAddress("192.168.100.0/25")

	appConnectorGroupTypeAndName, _, appConnectorGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ZPAAppConnectorGroup)
	appConnectorGroupHCL := testAccCheckAppConnectorGroupConfigure(appConnectorGroupTypeAndName, appConnectorGroupGeneratedName, variable.AppConnectorDescription, variable.AppConnectorEnabled)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLSSConfigControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLSSConfigControllerConfigure(lssControllerTypeAndName, lssControllerGeneratedName, lssControllerGeneratedName, lssControllerGeneratedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, rIP, rPort, variable.LSSControllerEnabled, variable.LSSControllerTLSEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLSSConfigControllerExists(lssControllerTypeAndName, &lssConfig),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.name", "tf-acc-test-"+lssControllerGeneratedName),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.description", "tf-acc-test-"+lssControllerGeneratedName),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.enabled", strconv.FormatBool(variable.LSSControllerEnabled)),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.use_tls", strconv.FormatBool(variable.LSSControllerTLSEnabled)),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "policy_rule_resource.#", "1"),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "connector_groups.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckLSSConfigControllerConfigure(lssControllerTypeAndName, lssControllerGeneratedName, lssControllerGeneratedName, lssControllerGeneratedName, appConnectorGroupHCL, appConnectorGroupTypeAndName, rIP, rPort, variable.LSSControllerEnabled, variable.LSSControllerTLSEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLSSConfigControllerExists(lssControllerTypeAndName, &lssConfig),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.name", "tf-acc-test-"+lssControllerGeneratedName),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.description", "tf-acc-test-"+lssControllerGeneratedName),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.enabled", strconv.FormatBool(variable.LSSControllerEnabled)),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.use_tls", strconv.FormatBool(variable.LSSControllerTLSEnabled)),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "policy_rule_resource.#", "1"),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "connector_groups.#", "1"),
				),
			},
		},
	})
}

func testAccCheckLSSConfigControllerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ZPALSSController {
			continue
		}

		lss, _, err := client.lssconfigcontroller.Get(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("id %s still exists", rs.Primary.ID)
		}

		if lss != nil {
			return fmt.Errorf("lss config controller with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckLSSConfigControllerExists(resource string, lss *lssconfigcontroller.LSSConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("lss config controller Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no lss config controller ID is set")
		}
		apiClient := testAccProvider.Meta().(*Client)
		resp, _, err := apiClient.lssconfigcontroller.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp.LSSConfig.Name != rs.Primary.Attributes["config.0.name"] {
			return fmt.Errorf("name Not found in created attributes")
		}
		if resp.LSSConfig.Description != rs.Primary.Attributes["config.0.description"] {
			return fmt.Errorf("description Not found in created attributes")
		}
		return nil
	}
}

func testAccCheckLSSConfigControllerConfigure(resourceTypeAndName, generatedName, name, description, appConnectorGroupHCL, appConnectorGroupTypeAndName, lssHost string, rPort int, enabled, tlsEnabled bool) string {
	return fmt.Sprintf(`

// app connector group resource
%s

// lss controller resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		appConnectorGroupHCL,
		getLSSControllerResourceHCL(generatedName, name, description, appConnectorGroupTypeAndName, lssHost, rPort, enabled, tlsEnabled),

		// data source variables
		resourcetype.ZPALSSController,
		generatedName,
		resourceTypeAndName,
	)
}

func getLSSControllerResourceHCL(generatedName, name, description, appConnectorGroupTypeAndName, lssHost string, rPort int, enabled, tlsEnabled bool) string {
	return fmt.Sprintf(`

// Retrieve LSS Config Format
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
	log_type="zpn_trans_log"
}

// Retrieve the Policy Set ID from Policy Type SIEM_POLICY
data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

data "zpa_idp_controller" "this" {
	name = "BD_Okta_Users"
   }

# Retrieve the SCIM_GROUP ID(s)
data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "finance" {
  name     = "Finance"
  idp_name = "BD_Okta_Users"
}
resource "%s" "%s" {
	config {
		name            = "tf-acc-test-%s"
		description     = "tf-acc-test-%s"
		enabled         = "%s"
		use_tls         = "%s"
		lss_host        = "%s"
		lss_port        = "%d"
		format          = data.zpa_lss_config_log_type_formats.zpn_trans_log.json
		source_log_type = "zpn_trans_log"
	}
	policy_rule_resource {
		name   = "policy_rule_resource-lss_auth_logs"
		action = "LOG"
		policy_set_id = data.zpa_policy_type.lss_siem_policy.id
		conditions {
			negated  = false
			operator = "OR"
			operands {
			  object_type = "CLIENT_TYPE"
			  values      = ["zpn_client_type_exporter", "zpn_client_type_ip_anchoring", "zpn_client_type_zapp", "zpn_client_type_edge_connector", "zpn_client_type_machine_tunnel", "zpn_client_type_browser_isolation", "zpn_client_type_slogger", "zpn_client_type_zapp_partner", "zpn_client_type_branch_connector"]
			}
		  }
		conditions {
		negated  = false
		operator = "OR"
		operands {
			object_type = "SCIM_GROUP"
			entry_values {
			rhs = data.zpa_scim_groups.engineering.id
			lhs = data.zpa_idp_controller.this.id
			}
			entry_values {
			rhs = data.zpa_scim_groups.finance.id
			lhs = data.zpa_idp_controller.this.id
			}
		}
		}
	}
	connector_groups {
		id = [ "${%s.id}" ]
	}
	depends_on = [ %s ]
}
`,

		// resource variables
		resourcetype.ZPALSSController,
		generatedName,
		generatedName,
		generatedName,
		strconv.FormatBool(enabled),
		strconv.FormatBool(tlsEnabled),
		lssHost,
		rPort,
		appConnectorGroupTypeAndName,
		appConnectorGroupTypeAndName,
	)
}
