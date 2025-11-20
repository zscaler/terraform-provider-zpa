package resources_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"

	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/acctest"
	"github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
)

func TestAccPolicyAccessTimeoutRule_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	updatedRName := sdkacctest.RandomWithPrefix("tf-updated")
	randDesc := sdkacctest.RandString(20)
	resourceName := "zpa_policy_access_timeout_rule.test"
	zClient := acctest.TestClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolicyAccessTimeoutRuleDestroy(zClient),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAccessTimeoutRuleConfig(rName, randDesc),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPolicyAccessTimeoutRuleExists(zClient, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "RE_AUTH"),
					resource.TestCheckResourceAttr(resourceName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceName, "reauth_idle_timeout", "600"),
					resource.TestCheckResourceAttr(resourceName, "reauth_timeout", "172800"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
				),
			},
			// Update test
			{
				Config: testAccPolicyAccessTimeoutRuleConfig(updatedRName, randDesc),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPolicyAccessTimeoutRuleExists(zClient, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedRName),
					resource.TestCheckResourceAttr(resourceName, "description", randDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "RE_AUTH"),
					resource.TestCheckResourceAttr(resourceName, "operator", "AND"),
					resource.TestCheckResourceAttr(resourceName, "reauth_idle_timeout", "600"),
					resource.TestCheckResourceAttr(resourceName, "reauth_timeout", "172800"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
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

func testAccCheckPolicyAccessTimeoutRuleExists(zClient *client.Client, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy access timeout rule not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("policy access timeout rule ID not set")
		}

		service := zClient.Service
		if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
			service = service.WithMicroTenant(microtenantID)
		}

		ctx := context.Background()
		accessPolicy, _, err := policysetcontroller.GetByPolicyType(ctx, zClient.Service, "TIMEOUT_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource TIMEOUT_POLICY. Received error: %s", err)
		}

		_, _, err = policysetcontroller.GetPolicyRule(ctx, service, accessPolicy.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resourceName, err)
		}

		return nil
	}
}

func testAccCheckPolicyAccessTimeoutRuleDestroy(zClient *client.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ctx := context.Background()
		accessPolicy, _, err := policysetcontroller.GetByPolicyType(ctx, zClient.Service, "TIMEOUT_POLICY")
		if err != nil {
			return fmt.Errorf("failed fetching resource TIMEOUT_POLICY. Received error: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "zpa_policy_access_timeout_rule" || rs.Primary.ID == "" {
				continue
			}

			service := zClient.Service
			if microtenantID := rs.Primary.Attributes["microtenant_id"]; microtenantID != "" {
				service = service.WithMicroTenant(microtenantID)
			}

			rule, _, err := policysetcontroller.GetPolicyRule(ctx, service, accessPolicy.ID, rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("id %s already exists", rs.Primary.ID)
			}

			if rule != nil {
				return fmt.Errorf("policy timeout rule with id %s exists and wasn't destroyed", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccPolicyAccessTimeoutRuleConfig(rName, desc string) string {
	return fmt.Sprintf(`
resource "zpa_policy_access_timeout_rule" "test" {
	name          		= "%s"
	description   		= "%s"
	action              = "RE_AUTH"
	reauth_idle_timeout = "600"
	reauth_timeout      = "172800"
	operator      		= "AND"
	conditions {
		operator = "OR"
		operands {
		  object_type = "CLIENT_TYPE"
		  lhs         = "id"
		  rhs         = "zpn_client_type_exporter"
		}
	}
}
`, rName, desc)
}
