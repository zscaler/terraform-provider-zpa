package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func resourcePolicyRedictionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyRedictionRuleCreate,
		ReadContext:   resourcePolicyRedictionRuleRead,
		UpdateContext: resourcePolicyRedictionRuleUpdate,
		DeleteContext: resourcePolicyRedictionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"REDIRECTION_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(),
			map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "  This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"REDIRECT_DEFAULT",
						"REDIRECT_PREFERRED",
						"REDIRECT_ALWAYS",
					}, false),
				},
				"conditions": GetPolicyConditionsSchema([]string{
					"CLIENT_TYPE",
				}),
				"service_edge_groups": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "List of the service edge group IDs.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeSet,
								Required: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		),
	}
}

// validatePolicyRedirectionRuleAction validates the "action" attribute against "service_edge_groups" requirements.
func validatePolicyRedirectionRuleAction(d *schema.ResourceData) error {
	action := d.Get("action").(string)

	raw := d.Get("service_edge_groups")
	var serviceEdgeGroups []interface{}

	switch v := raw.(type) {
	case *schema.Set:
		serviceEdgeGroups = v.List()
	case []interface{}:
		serviceEdgeGroups = v
	case nil:
		serviceEdgeGroups = []interface{}{}
	default:
		return fmt.Errorf("service_edge_groups has unexpected type %T", raw)
	}

	switch action {
	case "REDIRECT_PREFERRED", "REDIRECT_ALWAYS":
		if len(serviceEdgeGroups) == 0 {
			return fmt.Errorf("one or more ZPA Private Service Edge groups must be selected when the Private Service Edge Selection Method is %s", action)
		}
	case "REDIRECT_DEFAULT":
		if len(serviceEdgeGroups) > 0 {
			return fmt.Errorf("ZPA Private Service Edge groups must be empty when the Private Service Edge Selection Method is REDIRECT_DEFAULT")
		}
	}

	return nil
}

func resourcePolicyRedictionRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Validate the "action" and "service_edge_groups" attributes
	if err := validatePolicyRedirectionRuleAction(d); err != nil {
		return diag.FromErr(err)
	}

	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	req, err := expandCreatePolicyRedirectionRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating zpa policy redirection rule with request\n%+v\n", req)
	if err := ValidateConditions(ctx, req.Conditions, zClient, req.MicroTenantID); err == nil {
		policysetcontroller, _, err := policysetcontroller.CreateRule(ctx, service, req)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(policysetcontroller.ID)

		return resourcePolicyRedictionRuleRead(ctx, d, meta)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't validate the zpa policy redirection (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name))
	}
}

func resourcePolicyRedictionRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	globalPolicySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, "REDIRECTION_POLICY")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", globalPolicySet.ID, d.Id())
	resp, _, err := policysetcontroller.GetPolicyRule(ctx, service, globalPolicySet.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*errorx.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Got Policy Set Redirection Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("action", resp.Action)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))
	_ = d.Set("service_edge_groups", flattenServiceEdgeGroupSimple(resp.ServiceEdgeGroups))
	return nil
}

func resourcePolicyRedictionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Validate the "action" and "service_edge_groups" attributes
	if err := validatePolicyRedirectionRuleAction(d); err != nil {
		return diag.FromErr(err)
	}

	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	globalPolicySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, "REDIRECTION_POLICY")
	if err != nil {
		return diag.FromErr(err)
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy redirection rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyRedirectionRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := ValidateConditions(ctx, req.Conditions, zClient, req.MicroTenantID); err == nil {
		if _, _, err := policysetcontroller.GetPolicyRule(ctx, service, globalPolicySet.ID, ruleID); err != nil {
			if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
				d.SetId("")
				return nil
			}
		}

		if _, err := policysetcontroller.UpdateRule(ctx, service, globalPolicySet.ID, ruleID, req); err != nil {
			return diag.FromErr(err)
		}

		return resourcePolicyRedictionRuleRead(ctx, d, meta)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't validate the zpa policy redirection (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name))
	}
}

func resourcePolicyRedictionRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	globalPolicySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, "REDIRECTION_POLICY")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting policy redirection rule with id %v\n", d.Id())

	if _, err := policysetcontroller.Delete(ctx, service, globalPolicySet.ID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandCreatePolicyRedirectionRule(d *schema.ResourceData) (*policysetcontroller.PolicyRule, error) {
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
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Action:        d.Get("action").(string),
		ActionID:      d.Get("action_id").(string),
		CustomMsg:     d.Get("custom_msg").(string),
		Operator:      d.Get("operator").(string),
		PolicySetID:   policySetID,
		PolicyType:    d.Get("policy_type").(string),
		Priority:      d.Get("priority").(string),
		MicroTenantID: GetString(d.Get("microtenant_id")),
		Conditions:    conditions,
		ServiceEdgeGroups: func() []serviceedgegroup.ServiceEdgeGroup {
			groups := expandCommonServiceEdgeGroups(d)
			if groups == nil {
				return []serviceedgegroup.ServiceEdgeGroup{}
			}
			return groups
		}(),
		// ServiceEdgeGroups: expandCommonServiceEdgeGroups(d),
	}, nil
}
