package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

func resourcePolicyForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyForwardingRuleCreate,
		Read:   resourcePolicyForwardingRuleRead,
		Update: resourcePolicyForwardingRuleUpdate,
		Delete: resourcePolicyForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"CLIENT_FORWARDING_POLICY", "BYPASS_POLICY"}),
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
					"BRANCH_CONNECTOR_GROUP",
					"EDGE_CONNECTOR_GROUP",
					"POSTURE",
					"MACHINE_GRP",
					"TRUSTED_NETWORK",
					"PLATFORM",
					"IDP",
					"SAML",
					"SCIM",
					"SCIM_GROUP",
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
	log.Printf("[INFO] Creating zpa policy rule with request\n%+v\n", req)

	if err := ValidateConditions(req.Conditions, zClient, req.MicroTenantID); err != nil {
		return err
	}
	policysetcontroller, _, err := zClient.policysetcontroller.Create(req)
	if err != nil {
		return err
	}
	d.SetId(policysetcontroller.ID)

	return resourcePolicyForwardingRuleRead(d, m)
}

func resourcePolicyForwardingRuleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	globalPolicySet, _, err := zClient.policysetcontroller.GetByPolicyType("CLIENT_FORWARDING_POLICY")
	if err != nil {
		return err
	}
	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", globalPolicySet.ID, d.Id())
	resp, _, err := zClient.policysetcontroller.GetPolicyRule(globalPolicySet.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Got Policy Set Forwarding Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("action", resp.Action)
	_ = d.Set("action_id", resp.ActionID)
	_ = d.Set("custom_msg", resp.CustomMsg)
	_ = d.Set("bypass_default_rule", resp.BypassDefaultRule)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("priority", resp.Priority)
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))

	return nil
}

func resourcePolicyForwardingRuleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicySet, _, err := zClient.policysetcontroller.GetByPolicyType("CLIENT_FORWARDING_POLICY")
	if err != nil {
		return err
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy forwarding rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyForwardingRule(d)
	if err != nil {
		return err
	}
	if err := ValidateConditions(req.Conditions, zClient, req.MicroTenantID); err == nil {
		if _, _, err := zClient.policysetcontroller.GetPolicyRule(globalPolicySet.ID, ruleID); err != nil {
			if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
				d.SetId("")
				return nil
			}
		}

		if _, err := zClient.policysetcontroller.Update(globalPolicySet.ID, ruleID, req); err != nil {
			return err
		}

		return resourcePolicyForwardingRuleRead(d, m)
	} else {
		return err
	}

}

func resourcePolicyForwardingRuleDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicySet, _, err := zClient.policysetcontroller.GetByPolicyType("CLIENT_FORWARDING_POLICY")
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy forwarding rule with id %v\n", d.Id())

	if _, err := zClient.policysetcontroller.Delete(globalPolicySet.ID, d.Id()); err != nil {
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
		ID:                d.Get("id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Action:            d.Get("action").(string),
		ActionID:          d.Get("action_id").(string),
		CustomMsg:         d.Get("custom_msg").(string),
		BypassDefaultRule: d.Get("bypass_default_rule").(bool),
		DefaultRule:       d.Get("default_rule").(bool),
		Operator:          d.Get("operator").(string),
		PolicySetID:       policySetID,
		PolicyType:        d.Get("policy_type").(string),
		MicroTenantID:     GetString(d.Get("microtenant_id")),
		Priority:          d.Get("priority").(string),
		Conditions:        conditions,
	}, nil
}
