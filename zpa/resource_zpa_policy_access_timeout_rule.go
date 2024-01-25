package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

func resourcePolicyTimeoutRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyTimeoutRuleCreate,
		Read:   resourcePolicyTimeoutRuleRead,
		Update: resourcePolicyTimeoutRuleUpdate,
		Delete: resourcePolicyTimeoutRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"TIMEOUT_POLICY", "REAUTH_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(),
			map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "  This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"RE_AUTH",
					}, false),
				},
				"conditions": GetPolicyConditionsSchema([]string{
					"APP",
					"APP_GROUP",
					"CLIENT_TYPE",
					"IDP",
					"POSTURE",
					"PLATFORM",
					"SAML",
					"SCIM",
					"SCIM_GROUP",
				}),
			},
		),
	}
}

func resourcePolicyTimeoutRuleCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	req, err := expandCreatePolicyTimeoutRule(d)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy timeout rule with request\n%+v\n", req)
	if err := ValidateConditions(req.Conditions, zClient, req.MicroTenantID); err == nil {
		policysetcontroller, _, err := service.CreateRuleV1(req)
		if err != nil {
			return err
		}
		d.SetId(policysetcontroller.ID)

		return resourcePolicyTimeoutRuleRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy timeout (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}
}

func resourcePolicyTimeoutRuleRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("TIMEOUT_POLICY")
	if err != nil {
		return err
	}
	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", globalPolicySet.ID, d.Id())
	resp, _, err := service.GetPolicyRule(globalPolicySet.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[INFO] Got Policy Set Timeout Rule:\n%+v\n", resp)
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
	_ = d.Set("microtenant_id", resp.MicroTenantID)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))

	return nil
}

func resourcePolicyTimeoutRuleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("TIMEOUT_POLICY")
	if err != nil {
		return err
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy timeout rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyTimeoutRule(d)
	if err != nil {
		return err
	}
	if err := ValidateConditions(req.Conditions, zClient, req.MicroTenantID); err == nil {
		if _, _, err := service.GetPolicyRule(globalPolicySet.ID, ruleID); err != nil {
			if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
				d.SetId("")
				return nil
			}
		}

		if _, err := service.UpdateRuleV1(globalPolicySet.ID, ruleID, req); err != nil {
			return err
		}

		return resourcePolicyTimeoutRuleRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy timeout (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}
}

func resourcePolicyTimeoutRuleDelete(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("TIMEOUT_POLICY")
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy timeout rule with id %v\n", d.Id())

	if _, err := service.Delete(globalPolicySet.ID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyTimeoutRule(d *schema.ResourceData) (*policysetcontroller.PolicyRule, error) {
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
		MicroTenantID:     GetString(d.Get("microtenant_id")),
		ReauthDefaultRule: d.Get("reauth_default_rule").(bool),
		ReauthIdleTimeout: d.Get("reauth_idle_timeout").(string),
		ReauthTimeout:     d.Get("reauth_timeout").(string),
		Conditions:        conditions,
	}, nil
}
