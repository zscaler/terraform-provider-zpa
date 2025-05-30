package zpa

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func resourcePolicyBrowserProtectionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyBrowserProtectionRuleCreate,
		ReadContext:   resourcePolicyBrowserProtectionRuleRead,
		UpdateContext: resourcePolicyBrowserProtectionRuleUpdate,
		DeleteContext: resourcePolicyBrowserProtectionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFuncV2([]string{"CLIENTLESS_SESSION_PROTECTION_POLICY"}),
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the name of the policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the description of the access policy.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "  This is for providing the rule action.",
				ValidateFunc: validation.StringInSlice([]string{
					"MONITOR",
					"DO_NOT_MONITOR",
				}, false),
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "This is for proviidng the set of conditions for the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AND",
								"OR",
							}, false),
						},
						"operands": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "This signifies the various policy criteria.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "This denotes a list of values for the given object type. The value depend upon the key. If rhs is defined this list will be ignored",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"object_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "  This is for specifying the policy critiera.",
										ValidateFunc: validation.StringInSlice([]string{
											"APP",
											"APP_GROUP",
											"CLIENT_TYPE",
											"SAML",
											"SCIM",
											"SCIM_GROUP",
											"USER_PORTAL",
										}, false),
									},
									"entry_values": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rhs": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"lhs": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourcePolicyBrowserProtectionRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "CLIENTLESS_SESSION_PROTECTION_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CLIENTLESS_SESSION_PROTECTION_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

	req, err := expandBrowserProtectionPolicyRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating zpa policy browser protection rule with request\n%+v\n", req)

	policysetcontrollerv2, _, err := policysetcontrollerv2.CreateRule(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policysetcontrollerv2.ID)

	return resourcePolicyBrowserProtectionRuleRead(ctx, d, meta)
}

func resourcePolicyBrowserProtectionRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CLIENTLESS_SESSION_PROTECTION_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, d.Id())
	if err != nil {
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	v2PolicyRule := ConvertV1ResponseToV2Request(*resp)

	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", v2PolicyRule.Name)
	_ = d.Set("description", v2PolicyRule.Description)
	_ = d.Set("action", v2PolicyRule.Action)
	_ = d.Set("policy_set_id", policySetID)
	_ = d.Set("microtenant_id", v2PolicyRule.MicroTenantID)
	_ = d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))

	return nil
}

func resourcePolicyBrowserProtectionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "CLIENTLESS_SESSION_PROTECTION_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CLIENTLESS_SESSION_PROTECTION_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating policy browser protection rule ID: %v\n", ruleID)
	req, err := expandBrowserProtectionPolicyRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

	// Checking the current state of the rule to handle cases where it might have been deleted outside Terraform
	_, respErr, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if respErr != nil && (respErr.StatusCode == http.StatusNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	_, err = policysetcontrollerv2.UpdateRule(ctx, service, policySetID, ruleID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyBrowserProtectionRuleRead(ctx, d, meta)
}

func resourcePolicyBrowserProtectionRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Assume "CLIENTLESS_SESSION_PROTECTION_POLICY" is the policy type for this resource. Adjust as needed.
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CLIENTLESS_SESSION_PROTECTION_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete policy browser protection rule: %w", err))
	}

	return nil
}

func expandBrowserProtectionPolicyRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	return &policysetcontrollerv2.PolicyRule{
		ID:          d.Get("id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Action:      d.Get("action").(string),
		PolicySetID: policySetID,
		Conditions:  conditions,
	}, nil
}
