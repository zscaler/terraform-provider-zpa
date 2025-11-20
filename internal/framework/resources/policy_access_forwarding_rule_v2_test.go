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

func TestAccPolicyAccessForwardingRuleV2_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	randDesc := sdkacctest.RandString(20)
	resourceName := "zpa_policy_access_forwarding_rule_v2.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolicyAccessForwardingRuleV2Destroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAccessForwardingRuleV2Config(rName, randDesc),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPolicyAccessForwardingRuleV2Exists(zClient, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "BYPASS"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
				),
			},
			// Update test
			{
				Config: testAccPolicyAccessForwardingRuleV2Config(rName, randDesc),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPolicyAccessForwardingRuleV2Exists(zClient, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "BYPASS"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
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

func testAccCheckPolicyAccessForwardingRuleV2Exists(zClient *client.Client, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy access forwarding rule v2 not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("policy access forwarding rule v2 ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(ctx, zClient.Service, "CLIENT_FORWARDING_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource CLIENT_FORWARDING_POLICY. Received error: %s", err)
		}

		_, _, err = policysetcontrollerv2.GetPolicyRule(ctx, service, accessPolicy.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resourceName, err)
		}

		return nil
	}
}

func testAccCheckPolicyAccessForwardingRuleV2Destroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		accessPolicy, _, err := policysetcontrollerv2.GetByPolicyType(ctx, zClient.Service, "CLIENT_FORWARDING_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource CLIENT_FORWARDING_POLICY. Received error: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_policy_access_forwarding_rule_v2" || rs.Primary.ID == "" {
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
				return fmt.Errorf("policy forwarding rule with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccPolicyAccessForwardingRuleV2Config(rName, desc string) string {
	return fmt.Sprintf(`
data "zpa_posture_profile" "crwd_zta_score_80" {
	name = "CrowdStrike_ZPA_ZTA_80"
}

data "zpa_idp_controller" "bd_user_okta" {
    name = "BD_Okta_Users"
}

data "zpa_scim_groups" "a000" {
	name     = "A000"
	idp_name = "BD_Okta_Users"
}

data "zpa_scim_groups" "b000" {
	name     = "B000"
	idp_name = "BD_Okta_Users"
}

resource "zpa_policy_access_forwarding_rule_v2" "test" {
	name          		= "%s"
	description   		= "%s"
	action              = "BYPASS"
	conditions {
		operator = "OR"
		operands {
		  object_type = "POSTURE"
		  entry_values {
			lhs = data.zpa_posture_profile.crwd_zta_score_80.posture_udid
			rhs = false
		  }
		}
	}
	conditions {
		operator = "OR"
		operands {
			object_type = "SCIM_GROUP"
			entry_values {
				lhs = data.zpa_idp_controller.bd_user_okta.id
				rhs = data.zpa_scim_groups.a000.id
			}
			entry_values {
				lhs = data.zpa_idp_controller.bd_user_okta.id
				rhs = data.zpa_scim_groups.b000.id
			}
		}
	}
}
`, rName, desc)
}
