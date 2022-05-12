package zpa

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/lssconfigcontroller"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/resourcetype"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/method"
	"github.com/zscaler/terraform-provider-zpa/zpa/common/testing/variable"
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
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.name", "test-lss-config-"+lssControllerGeneratedName),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.description", "test-lss-config-"+lssControllerGeneratedName),
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
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.name", "test-lss-config-"+lssControllerGeneratedName),
					resource.TestCheckResourceAttr(lssControllerTypeAndName, "config.0.description", "test-lss-config-"+lssControllerGeneratedName),
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

resource "%s" "%s" {
	config {
		name            = "test-lss-config-%s"
		description     = "test-lss-config-%s"
		enabled         = "%s"
		use_tls         = "%s"
		lss_host        = "%s"
		lss_port        = "%d"
		format          = data.zpa_lss_config_log_type_formats.zpn_trans_log.json
		source_log_type = "zpn_trans_log"
	}
	policy_rule_resource {
		name   = "policy_rule_resource-lss_auth_logs"
		action = "ALLOW"
		conditions {
		  negated  = false
		  operator = "OR"
		  operands {
			object_type = "CLIENT_TYPE"
			values      = ["zpn_client_type_exporter"]
		  }
		  operands {
			object_type = "CLIENT_TYPE"
			values      = ["zpn_client_type_ip_anchoring"]
		  }
		  operands {
			object_type = "CLIENT_TYPE"
			values      = ["zpn_client_type_zapp"]
		  }
		  operands {
			object_type = "CLIENT_TYPE"
			values      = ["zpn_client_type_edge_connector"]
		  }
		  operands {
			object_type = "CLIENT_TYPE"
			values      = ["zpn_client_type_machine_tunnel"]
		  }
		  operands {
			object_type = "CLIENT_TYPE"
			values      = ["zpn_client_type_browser_isolation"]
		  }
		  operands {
			object_type = "CLIENT_TYPE"
			values      = ["zpn_client_type_slogger"]
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
