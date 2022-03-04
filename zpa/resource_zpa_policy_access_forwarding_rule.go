package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/policysetcontroller"
)

func resourcePolicyForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyForwardingRuleCreate,
		Read:   resourcePolicyForwardingRuleRead,
		Update: resourcePolicyForwardingRuleUpdate,
		Delete: resourcePolicyForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: importPolicyStateFunc([]string{"CLIENT_FORWARDING_POLICY", "BYPASS_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(),
			map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "  This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"BYPASS",
						"INTERCEPT",
						"INTERCEPT_ACCESSIBLE",
					}, false),
				},
				"conditions": GetPolicyConditionsSchema([]string{
					"APP",
					"APP_GROUP",
					"CLIENT_TYPE",
					"CLOUD_CONNECTOR_GROUP",
					"IDP",
					"POSTURE",
					"SAML",
					"SCIM",
					"SCIM_GROUP",
					"TRUSTED_NETWORK",
				}),
			},
		),
	}
}

func resourcePolicyForwardingRuleCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req, err := expandCreatePolicyForwardingRule(d)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy forwarding rule with request\n%+v\n", req)
	if ValidateConditions(req.Conditions, zClient) {
		policysetcontroller, _, err := zClient.policysetcontroller.Create(req)
		if err != nil {
			return err
		}
		d.SetId(policysetcontroller.ID)
		order, ok := d.GetOk("rule_order")
		if ok {
			reorder(order, policysetcontroller.PolicySetID, policysetcontroller.ID, zClient)
		}
		return resourcePolicyForwardingRuleRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy forwarding (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}

}

func resourcePolicyForwardingRuleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	globalPolicyForwarding, _, err := zClient.policytype.GetBypass()
	if err != nil {
		return err
	}
	log.Printf("[INFO] Getting Policy Set Forwarding Rule: globalPolicySet:%s id: %s\n", globalPolicyForwarding.ID, d.Id())
	resp, _, err := zClient.policysetcontroller.Get(globalPolicyForwarding.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy forwarding rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Got Policy Set Forwarding Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("action", resp.Action)
	_ = d.Set("action_id", resp.ActionID)
	_ = d.Set("custom_msg", resp.CustomMsg)
	_ = d.Set("description", resp.Description)
	_ = d.Set("name", resp.Name)
	_ = d.Set("bypass_default_rule", resp.BypassDefaultRule)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("priority", resp.Priority)
	_ = d.Set("reauth_default_rule", resp.ReauthDefaultRule)
	_ = d.Set("reauth_idle_timeout", resp.ReauthIdleTimeout)
	_ = d.Set("reauth_timeout", resp.ReauthTimeout)
	_ = d.Set("rule_order", resp.RuleOrder)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))

	return nil
}

func resourcePolicyForwardingRuleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicyForwarding, _, err := zClient.policytype.GetBypass()
	if err != nil {
		return err
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy forwarding rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyRule(d)
	if err != nil {
		return err
	}
	if ValidateConditions(req.Conditions, zClient) {
		if _, err := zClient.policysetcontroller.Update(globalPolicyForwarding.ID, ruleID, req); err != nil {
			return err
		}
		if d.HasChange("rule_order") {
			order, ok := d.GetOk("rule_order")
			if ok {
				reorder(order, globalPolicyForwarding.ID, ruleID, zClient)
			}
		}
		return resourcePolicyForwardingRuleRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy forwarding (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}

}

func resourcePolicyForwardingRuleDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicyForwarding, _, err := zClient.policytype.GetBypass()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy forwarding rule with id %v\n", d.Id())

	if _, err := zClient.policysetcontroller.Delete(globalPolicyForwarding.ID, d.Id()); err != nil {
		return err
	}

	return nil

}

func expandCreatePolicyForwardingRule(d *schema.ResourceData) (*policysetcontroller.PolicyRule, error) {
	policySetID, ok := d.Get("policy_set_id").(string)
	if !ok {
		return nil, fmt.Errorf("policy_set_id is not set")
	}
	log.Printf("[INFO] action_id:%v\n", d.Get("action_id"))
	conditions, err := ExpandPolicyConditions(d)
	if err != nil {
		return nil, err
	}
	return &policysetcontroller.PolicyRule{
		Action:            d.Get("action").(string),
		ActionID:          d.Get("action_id").(string),
		CustomMsg:         d.Get("custom_msg").(string),
		Description:       d.Get("description").(string),
		ID:                d.Get("id").(string),
		Name:              d.Get("name").(string),
		Operator:          d.Get("operator").(string),
		PolicySetID:       policySetID,
		PolicyType:        d.Get("policy_type").(string),
		Priority:          d.Get("priority").(string),
		ReauthDefaultRule: d.Get("reauth_default_rule").(bool),
		ReauthIdleTimeout: d.Get("reauth_idle_timeout").(string),
		ReauthTimeout:     d.Get("reauth_timeout").(string),
		RuleOrder:         d.Get("rule_order").(string),
		Conditions:        conditions,
	}, nil
}
