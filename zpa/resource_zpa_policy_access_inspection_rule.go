package zpa

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func resourcePolicyInspectionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyInspectionRuleCreate,
		ReadContext:   resourcePolicyInspectionRuleRead,
		UpdateContext: resourcePolicyInspectionRuleUpdate,
		DeleteContext: resourcePolicyInspectionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"INSPECTION_POLICY"}),
		},

		Schema: MergeSchema(
			InspectionPolicySchema(),
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

func InspectionPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "This is the description of the access policy.",
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"policy_set_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "This is the name of the policy.",
		},
		"operator": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"AND",
				"OR",
			}, false),
		},
		"zpn_inspection_profile_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"microtenant_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func resourcePolicyInspectionRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
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
		policySetID, err = fetchPolicySetIDByType(ctx, zClient, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	req, err := expandCreatePolicyInspectionRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating zpa policy inspection rule with request\n%+v\n", req)

	if err := ValidateConditions(ctx, req.Conditions, zClient, GetString(d.Get("microtenant_id"))); err != nil {
		return diag.FromErr(err)
	}
	resp, _, err := policysetcontroller.CreateRule(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourcePolicyInspectionRuleRead(ctx, d, meta)
}

func resourcePolicyInspectionRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "INSPECTION_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting Policy Set Rule: policySetID:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontroller.GetPolicyRule(ctx, service, policySetID, d.Id())
	if err != nil {
		// Adjust this error handling to match how your client library exposes HTTP response details
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Got Policy Set Inspection Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("action", resp.Action)
	_ = d.Set("description", resp.Description)
	_ = d.Set("name", resp.Name)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("zpn_inspection_profile_id", resp.ZpnInspectionProfileID)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))

	return nil
}

func resourcePolicyInspectionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
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
		policySetID, err = fetchPolicySetIDByType(ctx, zClient, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy inspection rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyInspectionRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := ValidateConditions(ctx, req.Conditions, zClient, GetString(d.Get("microtenant_id"))); err != nil {
		return diag.FromErr(err)
	}

	if _, err := policysetcontroller.UpdateRule(ctx, service, policySetID, ruleID, req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyInspectionRuleRead(ctx, d, meta)
}

func resourcePolicyInspectionRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
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
		policySetID, err = fetchPolicySetIDByType(ctx, zClient, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	log.Printf("[INFO] Deleting policy inspection rule with id %v\n", d.Id())

	if _, err := policysetcontroller.Delete(ctx, service, policySetID, d.Id()); err != nil {
		return diag.FromErr(err)
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
		Description:            d.Get("description").(string),
		ID:                     d.Get("id").(string),
		Name:                   d.Get("name").(string),
		Operator:               d.Get("operator").(string),
		PolicySetID:            policySetID,
		MicroTenantID:          GetString(d.Get("microtenant_id")),
		ZpnInspectionProfileID: d.Get("zpn_inspection_profile_id").(string),
		Conditions:             conditions,
	}, nil
}
