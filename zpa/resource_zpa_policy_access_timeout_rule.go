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

func resourcePolicyTimeoutRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	var policySetID string
	var err error

	// Check if policy_set_id is provided by the user
	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		// Fetch policy_set_id based on the policy_type
		policySetID, err = fetchPolicySetIDByType(client, "TIMEOUT_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}

	req, err := expandCreatePolicyTimeoutRule(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy timeout rule with request\n%+v\n", req)
	if err := ValidateConditions(req.Conditions, client, GetString(d.Get("microtenant_id"))); err != nil {
		return err
	}

	policysetcontroller, _, err := client.policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id"))).CreateRule(req)
	if err != nil {
		return err
	}

	d.SetId(policysetcontroller.ID)

	return resourcePolicyTimeoutRuleRead(d, m)
}

func resourcePolicyTimeoutRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(client, "TIMEOUT_POLICY", microTenantID)
	if err != nil {
		return err
	}

	service := client.policysetcontroller.WithMicroTenant(microTenantID)
	log.Printf("[INFO] Getting Policy Set Rule: policySetID:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := service.GetPolicyRule(policySetID, d.Id())
	if err != nil {
		// Adjust this error handling to match how your client library exposes HTTP response details
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

func resourcePolicyTimeoutRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	var policySetID string
	var err error

	// Check if policy_set_id is provided by the user, otherwise fetch it
	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(client, "TIMEOUT_POLICY", GetString(d.Get("microtenant_id")))
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
	// Replace ValidatePolicyRuleConditions with ValidateConditions
	if err := ValidateConditions(req.Conditions, client, GetString(d.Get("microtenant_id"))); err != nil {
		return err
	}

	if _, err := client.policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id"))).UpdateRule(policySetID, ruleID, req); err != nil {
		return err
	}

	return resourcePolicyTimeoutRuleRead(d, m)
}

func resourcePolicyTimeoutRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	var policySetID string
	var err error

	// Check if policy_set_id is provided by the user, otherwise fetch it based on policy_type
	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		// Assuming "TIMEOUT_POLICY" as policy type for demonstration
		policySetID, err = fetchPolicySetIDByType(client, "TIMEOUT_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}

	log.Printf("[INFO] Deleting policy timeout rule with id %v\n", d.Id())

	if _, err := client.policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id"))).Delete(policySetID, d.Id()); err != nil {
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
