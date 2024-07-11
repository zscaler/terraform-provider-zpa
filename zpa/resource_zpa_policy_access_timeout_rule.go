package zpa

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func resourcePolicyTimeoutRuleCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(zClient, "TIMEOUT_POLICY", microTenantID)
		if err != nil {
			return err
		}
	}

	req, err := expandCreatePolicyTimeoutRule(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy timeout rule with request\n%+v\n", req)
	if err := ValidateConditions(req.Conditions, zClient, microTenantID); err != nil {
		return err
	}

	resp, _, err := policysetcontroller.CreateRule(service, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePolicyTimeoutRuleRead(d, meta)
}

func resourcePolicyTimeoutRuleRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(zClient, "TIMEOUT_POLICY", microTenantID)
	if err != nil {
		return err
	}

	service := zClient.PolicySetController
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Getting Policy Set Rule: policySetID:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontroller.GetPolicyRule(service, policySetID, d.Id())
	if err != nil {
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
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

func resourcePolicyTimeoutRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	service := client.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(client, "TIMEOUT_POLICY", microTenantID)
		if err != nil {
			return err
		}
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy timeout rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyTimeoutRule(d, policySetID)
	if err != nil {
		return err
	}

	if err := ValidateConditions(req.Conditions, client, microTenantID); err != nil {
		return err
	}

	if _, err := policysetcontroller.UpdateRule(service, policySetID, ruleID, req); err != nil {
		return err
	}

	return resourcePolicyTimeoutRuleRead(d, meta)
}

func resourcePolicyTimeoutRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	service := client.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(client, "TIMEOUT_POLICY", microTenantID)
		if err != nil {
			return err
		}
	}

	log.Printf("[INFO] Deleting policy timeout rule with id %v\n", d.Id())

	if _, err := policysetcontroller.Delete(service, policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyTimeoutRule(d *schema.ResourceData, policySetID string) (*policysetcontroller.PolicyRule, error) {
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
		PolicyType:        d.Get("policy_type").(string),
		Priority:          d.Get("priority").(string),
		MicroTenantID:     GetString(d.Get("microtenant_id")),
		ReauthDefaultRule: d.Get("reauth_default_rule").(bool),
		ReauthIdleTimeout: d.Get("reauth_idle_timeout").(string),
		ReauthTimeout:     d.Get("reauth_timeout").(string),
		PolicySetID:       policySetID,
		Conditions:        conditions,
	}, nil
}
