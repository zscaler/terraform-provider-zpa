package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func resourcePolicyCredentialAccessRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCredentialAccessRuleCreate,
		ReadContext:   resourcePolicyCredentialAccessRuleRead,
		UpdateContext: resourcePolicyCredentialAccessRuleUpdate,
		DeleteContext: resourcePolicyCredentialAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFuncV2([]string{"CREDENTIAL_POLICY"}),
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
					"INJECT_CREDENTIALS",
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
											"CONSOLE",
											"SAML",
											"SCIM",
											"SCIM_GROUP",
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
			"credential": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"credential", "credential_pool"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"credential_pool": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"credential", "credential_pool"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
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

func resourcePolicyCredentialAccessRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "CREDENTIAL_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CREDENTIAL_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandCredentialPolicyRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating zpa policy credential rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

	resp, _, err := policysetcontrollerv2.CreateRule(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourcePolicyCredentialAccessRuleRead(ctx, d, meta)
}

func resourcePolicyCredentialAccessRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CREDENTIAL_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", policySetID, d.Id())
	resp, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
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
	_ = d.Set("credential", flattenCredential(resp.Credential))
	_ = d.Set("credential_pool", flattenCredential(resp.CredentialPool))

	// Ensure microtenant_id is being correctly set in state
	if v2PolicyRule.MicroTenantID != "" {
		log.Printf("[INFO] Setting microtenant_id in state: %s\n", v2PolicyRule.MicroTenantID)
		_ = d.Set("microtenant_id", v2PolicyRule.MicroTenantID)
	} else {
		log.Printf("[WARN] microtenant_id is empty in response.")
		_ = d.Set("microtenant_id", "")
	}

	return nil
}

func resourcePolicyCredentialAccessRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CREDENTIAL_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("policy_set_id", policySetID)
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy credential rule ID: %v\n", ruleID)

	req, err := expandCredentialPolicyRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve the current state from the API
	resp, _, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, ruleID)
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Handle credential or credential pool fallback logic
	if req.Credential == nil && req.CredentialPool == nil {
		// If both are nil in request, check if we should fall back to API values
		if resp != nil {
			if resp.Credential != nil && resp.Credential.ID != "" {
				req.Credential = resp.Credential
			} else if resp.CredentialPool != nil && resp.CredentialPool.ID != "" {
				req.CredentialPool = resp.CredentialPool
			}
		}
	}

	// Final validation - ensure we have either credential or credential pool
	if (req.Credential == nil || req.Credential.ID == "") && (req.CredentialPool == nil || req.CredentialPool.ID == "") {
		return diag.Errorf("either credential or credential_pool block must be present and contain an ID during update")
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, ruleID, req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyCredentialAccessRuleRead(ctx, d, meta)
}

func resourcePolicyCredentialAccessRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "CREDENTIAL_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting policy credential rule with id %v\n", d.Id())

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete policy credential rule: %w", err))
	}

	return nil
}

func flattenCredential(credential *policysetcontrollerv2.Credential) []interface{} {
	if credential == nil || credential.ID == "" {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id": credential.ID,
		},
	}
}

func expandCredentialPolicyRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	var credential *policysetcontrollerv2.Credential
	var credentialPool *policysetcontrollerv2.Credential

	if v, ok := d.GetOk("credential"); ok {
		if items := v.([]interface{}); len(items) > 0 {
			m := items[0].(map[string]interface{})
			if id, ok := m["id"].(string); ok && id != "" {
				credential = &policysetcontrollerv2.Credential{ID: id}
			}
		}
	}

	if v, ok := d.GetOk("credential_pool"); ok {
		if items := v.([]interface{}); len(items) > 0 {
			m := items[0].(map[string]interface{})
			if id, ok := m["id"].(string); ok && id != "" {
				credentialPool = &policysetcontrollerv2.Credential{ID: id}
			}
		}
	}

	// Validate mutual exclusivity
	if credential != nil && credentialPool != nil {
		return nil, fmt.Errorf("only one of 'credential' or 'credential_pool' can be set")
	}

	return &policysetcontrollerv2.PolicyRule{
		ID:             d.Get("id").(string),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Action:         d.Get("action").(string),
		MicroTenantID:  d.Get("microtenant_id").(string),
		PolicySetID:    policySetID,
		Conditions:     conditions,
		Credential:     credential,
		CredentialPool: credentialPool,
	}, nil
}
