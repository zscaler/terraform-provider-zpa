package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccPolicyAccessRuleV2_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	randDesc := sdkacctest.RandString(10)
	appConnectorGroupName := sdkacctest.RandString(8)
	segmentGroupName := sdkacctest.RandString(8)
	resourceName := "zpa_policy_access_rule_v2.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolicyAccessRuleV2Destroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAccessRuleV2Config(rName, randDesc, appConnectorGroupName, segmentGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPolicyAccessRuleV2Exists(zClient, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "6"),
				),
			},
			// Update test
			{
				Config: testAccPolicyAccessRuleV2Config(rName, randDesc, appConnectorGroupName, segmentGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPolicyAccessRuleV2Exists(zClient, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceName, "app_connector_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "6"),
				),
			},
			// Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPolicyAccessRuleV2Exists(zClient *client.Client, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy access rule v2 not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("policy access rule v2 ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(ctx, zClient.Service, "ACCESS_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource ACCESS_POLICY. Received error: %s", err)
		}

		_, _, err = policysetcontrollerv2.GetPolicyRule(ctx, service, accessPolicy.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resourceName, err)
		}

		return nil
	}
}

func testAccCheckPolicyAccessRuleV2Destroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(ctx, zClient.Service, "ACCESS_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource ACCESS_POLICY. Received error: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_policy_access_rule_v2" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			rule, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, accessPolicy.ID, rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("id %s already exists", rs.Primary.ID)
			}

			if rule != nil {
				return fmt.Errorf("policy access rule v2 with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccPolicyAccessRuleV2Config(rName, desc, appConnectorGroupName, segmentGroupName string) string {
	return fmt.Sprintf(`
resource "zpa_app_connector_group" "test" {
	name                          = "%s"
	description                   = "testAcc_app_connector_group"
	enabled                       = "true"
	country_code                  = "US"
	city_country                  = "San Jose, US"
	latitude                      = "37.33874"
	longitude                     = "-121.8852525"
	location                      = "San Jose, CA, USA"
	upgrade_day                   = "SUNDAY"
	upgrade_time_in_secs          = "66600"
	dns_query_type                = "IPV4_IPV6"
	tcp_quick_ack_app 			  = true
	tcp_quick_ack_assistant 	  = true
	tcp_quick_ack_read_assistant  = true
}

resource "zpa_segment_group" "test" {
	name = "tf-acc-test-%s"
	description = "testAcc_segment_group"
	enabled = "true"
}

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

resource "zpa_policy_access_rule_v2" "test" {
	name          		= "%s"
	description   		= "%s"
	action        		= "ALLOW"
	operator      		= "AND"
	app_connector_groups {
		id = [zpa_app_connector_group.test.id]
	}
	conditions {
		operator = "OR"
		operands {
			object_type = "APP_GROUP"
			values      = [zpa_segment_group.test.id]
		}
	}
	conditions {
		operator = "OR"
		operands {
		  object_type = "SCIM_GROUP"
		  entry_values {
			rhs = data.zpa_scim_groups.a000.id
			lhs = data.zpa_idp_controller.this.id
		  }
		  entry_values {
			rhs = data.zpa_scim_groups.b000.id
			lhs = data.zpa_idp_controller.this.id
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
	depends_on = [zpa_app_connector_group.test, zpa_segment_group.test]
}
`, appConnectorGroupName, segmentGroupName, rName, desc)
}
