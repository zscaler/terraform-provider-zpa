package zpa

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

func resourcePolicyInspectionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyInspectionRuleCreate,
		Read:   resourcePolicyInspectionRuleRead,
		Update: resourcePolicyInspectionRuleUpdate,
		Delete: resourcePolicyInspectionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"INSPECTION_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(),
			map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"INSPECT",
						"BYPASS_INSPECT",
					}, false),
				},
				"conditions": GetPolicyConditionsSchema([]string{
					"APP",
					"APP_GROUP",
					"CLIENT_TYPE",
					"EDGE_CONNECTOR_GROUP",
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

func resourcePolicyInspectionRuleCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
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
		policySetID, err = fetchPolicySetIDByType(zClient, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}
	req, err := expandCreatePolicyInspectionRule(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy inspection rule with request\n%+v\n", req)

	if err := ValidateConditions(req.Conditions, zClient, GetString(d.Get("microtenant_id"))); err != nil {
		return err
	}
	resp, _, err := policysetcontroller.CreateRule(service, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePolicyInspectionRuleRead(d, m)
}

func resourcePolicyInspectionRuleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(zClient, "INSPECTION_POLICY", microTenantID)
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

	log.Printf("[INFO] Got Policy Set Inspection Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("action", resp.Action)
	_ = d.Set("action_id", resp.ActionID)
	_ = d.Set("custom_msg", resp.CustomMsg)
	_ = d.Set("description", resp.Description)
	_ = d.Set("name", resp.Name)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("priority", resp.Priority)
	_ = d.Set("zpn_inspection_profile_id", resp.ZpnInspectionProfileID)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))

	return nil
}

func resourcePolicyInspectionRuleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
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
		policySetID, err = fetchPolicySetIDByType(zClient, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy inspection rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyInspectionRule(d, policySetID)
	if err != nil {
		return err
	}
	// Replace ValidatePolicyRuleConditions with ValidateConditions
	if err := ValidateConditions(req.Conditions, zClient, GetString(d.Get("microtenant_id"))); err != nil {
		return err
	}

	if _, err := policysetcontroller.UpdateRule(service, policySetID, ruleID, req); err != nil {
		return err
	}

	return resourcePolicyInspectionRuleRead(d, m)
}

func resourcePolicyInspectionRuleDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
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
		// Assuming "INSPECTION_POLICY" as policy type for demonstration
		policySetID, err = fetchPolicySetIDByType(zClient, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return err
		}
	}
	log.Printf("[INFO] Deleting policy inspection rule with id %v\n", d.Id())

	if _, err := policysetcontroller.Delete(service, policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyInspectionRule(d *schema.ResourceData, policySetID string) (*policysetcontroller.PolicyRule, error) {
	conditions, err := ExpandPolicyConditions(d)
	if err != nil {
		return nil, err
	}
	return &policysetcontroller.PolicyRule{
		Action:                 d.Get("action").(string),
		ActionID:               d.Get("action_id").(string),
		CustomMsg:              d.Get("custom_msg").(string),
		Description:            d.Get("description").(string),
		ID:                     d.Get("id").(string),
		Name:                   d.Get("name").(string),
		Operator:               d.Get("operator").(string),
		PolicySetID:            policySetID,
		PolicyType:             d.Get("policy_type").(string),
		Priority:               d.Get("priority").(string),
		MicroTenantID:          GetString(d.Get("microtenant_id")),
		ZpnInspectionProfileID: d.Get("zpn_inspection_profile_id").(string),
		Conditions:             conditions,
	}, nil
}
