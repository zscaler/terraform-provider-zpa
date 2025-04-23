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

func resourcePolicyAccessRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyAccessCreate,
		ReadContext:   resourcePolicyAccessRead,
		UpdateContext: resourcePolicyAccessUpdate,
		DeleteContext: resourcePolicyAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"ACCESS_POLICY", "GLOBAL_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(), map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "  This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"ALLOW",
						"DENY",
						"REQUIRE_APPROVAL",
					}, false),
				},
				"app_server_groups": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					Description: "List of the server group IDs.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"app_connector_groups": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					Description: "List of app-connector IDs.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"conditions": GetPolicyConditionsSchema([]string{
					"APP",
					"APP_GROUP",
					"LOCATION",
					"IDP",
					"SAML",
					"SCIM",
					"SCIM_GROUP",
					"CLIENT_TYPE",
					"POSTURE",
					"TRUSTED_NETWORK",
					"BRANCH_CONNECTOR_GROUP",
					"EDGE_CONNECTOR_GROUP",
					"MACHINE_GRP",
					"COUNTRY_CODE",
					"PLATFORM",
					"RISK_FACTOR_TYPE",
					"CHROME_ENTERPRISE",
				}),
			},
		),
	}
}

func resourcePolicyAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(ctx, zClient, "ACCESS_POLICY", microTenantID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	req, err := expandCreatePolicyRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating ZPA policy access rule with request\n%+v\n", req)

	if err := ValidateConditions(ctx, req.Conditions, zClient, microTenantID); err != nil {
		return diag.FromErr(err)
	}

	resp, _, err := policysetcontroller.CreateRule(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourcePolicyAccessRead(ctx, d, meta)
}

func resourcePolicyAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ACCESS_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting Policy Set Rule: policySetID:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontroller.GetPolicyRule(ctx, service, policySetID, d.Id())
	if err != nil {
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("description", resp.Description)
	_ = d.Set("name", resp.Name)
	_ = d.Set("action", resp.Action)
	_ = d.Set("action_id", resp.ActionID)
	_ = d.Set("custom_msg", resp.CustomMsg)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", policySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("priority", resp.Priority)
	_ = d.Set("lss_default_rule", resp.LSSDefaultRule)
	_ = d.Set("microtenant_id", microTenantID)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))
	_ = d.Set("app_server_groups", flattenCommonAppServerGroupSimple(resp.AppServerGroups))
	_ = d.Set("app_connector_groups", flattenCommonAppConnectorGroups(resp.AppConnectorGroups))

	return nil
}

func resourcePolicyAccessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(ctx, zClient, "ACCESS_POLICY", microTenantID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	ruleID := d.Id()
	req, err := expandCreatePolicyRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := ValidateConditions(ctx, req.Conditions, zClient, microTenantID); err != nil {
		return diag.FromErr(err)
	}

	if _, err := policysetcontroller.UpdateRule(ctx, service, policySetID, ruleID, req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyAccessRead(ctx, d, meta)
}

func resourcePolicyAccessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(ctx, zClient, "ACCESS_POLICY", microTenantID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	if _, err := policysetcontroller.Delete(ctx, service, policySetID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandCreatePolicyRule(d *schema.ResourceData, policySetID string) (*policysetcontroller.PolicyRule, error) {
	conditions, err := ExpandPolicyConditions(d)
	if err != nil {
		return nil, err
	}
	return &policysetcontroller.PolicyRule{
		ID:                 d.Get("id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Action:             d.Get("action").(string),
		ActionID:           d.Get("action_id").(string),
		BypassDefaultRule:  d.Get("bypass_default_rule").(bool),
		CustomMsg:          d.Get("custom_msg").(string),
		DefaultRule:        d.Get("default_rule").(bool),
		Operator:           d.Get("operator").(string),
		PolicySetID:        policySetID,
		PolicyType:         d.Get("policy_type").(string),
		Priority:           d.Get("priority").(string),
		MicroTenantID:      d.Get("microtenant_id").(string),
		LSSDefaultRule:     d.Get("lss_default_rule").(bool),
		Conditions:         conditions,
		AppServerGroups:    expandCommonServerGroups(d),
		AppConnectorGroups: expandCommonAppConnectorGroups(d),
	}, nil
}
