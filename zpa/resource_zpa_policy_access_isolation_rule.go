package zpa

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

func resourcePolicyIsolationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyIsolationRuleCreate,
		Read:   resourcePolicyIsolationRuleRead,
		Update: resourcePolicyIsolationRuleUpdate,
		Delete: resourcePolicyIsolationRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"ISOLATION_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(),
			map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "  This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"ISOLATE",
						"BYPASS_ISOLATE",
					}, false),
				},
				"conditions": GetPolicyConditionsSchema([]string{
					"APP",
					"APP_GROUP",
					"CLIENT_TYPE",
					"EDGE_CONNECTOR_GROUP",
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

func resourcePolicyIsolationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	// Check if policy_set_id is provided by the user
	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		// Fetch policy_set_id based on the policy_type
		policySetID, err = fetchPolicySetIDByType(zClient, "ISOLATION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}
	req, err := expandCreatePolicyIsolationRule(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy isolation rule with request\n%+v\n", req)
	if err := ValidateConditions(req.Conditions, zClient, GetString(d.Get("microtenant_id"))); err != nil {
		return err
	}

	resp, _, err := policysetcontroller.CreateRule(service, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePolicyIsolationRuleRead(d, meta)
}

func resourcePolicyIsolationRuleRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(zClient, "ISOLATION_POLICY", microTenantID)
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
		// Adjust this error handling to match how your client library exposes HTTP response details
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[INFO] Got Policy Set Isolation Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("action", resp.Action)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("zpn_cbi_profile_id", resp.ZpnCbiProfileID)
	_ = d.Set("zpn_isolation_profile_id", resp.ZpnIsolationProfileID)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))

	return nil
}

func resourcePolicyIsolationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	// Check if policy_set_id is provided by the user, otherwise fetch it
	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(zClient, "ISOLATION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy isolation rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyIsolationRule(d, policySetID)
	if err != nil {
		return err
	}
	if err := ValidateConditions(req.Conditions, zClient, GetString(d.Get("microtenant_id"))); err != nil {
		return err
	}

	if _, err := policysetcontroller.UpdateRule(service, policySetID, ruleID, req); err != nil {
		return err
	}

	return resourcePolicyIsolationRuleRead(d, meta)
}

func resourcePolicyIsolationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	// Check if policy_set_id is provided by the user, otherwise fetch it based on policy_type
	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		// Assuming "ISOLATION_POLICY" as policy type for demonstration
		policySetID, err = fetchPolicySetIDByType(zClient, "ISOLATION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}

	log.Printf("[INFO] Deleting policy isolation rule with id %v\n", d.Id())

	if _, err := policysetcontroller.Delete(service, policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyIsolationRule(d *schema.ResourceData, policySetID string) (*policysetcontroller.PolicyRule, error) {
	conditions, err := ExpandPolicyConditions(d)
	if err != nil {
		return nil, err
	}
	return &policysetcontroller.PolicyRule{
		ID:                    d.Get("id").(string),
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Action:                d.Get("action").(string),
		ActionID:              d.Get("action_id").(string),
		CustomMsg:             d.Get("custom_msg").(string),
		BypassDefaultRule:     d.Get("bypass_default_rule").(bool),
		DefaultRule:           d.Get("default_rule").(bool),
		Operator:              d.Get("operator").(string),
		PolicyType:            d.Get("policy_type").(string),
		ZpnCbiProfileID:       d.Get("zpn_cbi_profile_id").(string),
		ZpnIsolationProfileID: d.Get("zpn_isolation_profile_id").(string),
		Priority:              d.Get("priority").(string),
		MicroTenantID:         GetString(d.Get("microtenant_id")),
		Conditions:            conditions,
		PolicySetID:           policySetID,
	}, nil
}
