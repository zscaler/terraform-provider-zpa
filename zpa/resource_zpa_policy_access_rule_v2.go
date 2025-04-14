package zpa

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func resourcePolicyAccessRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyAccessV2Create,
		ReadContext:   resourcePolicyAccessV2Read,
		UpdateContext: resourcePolicyAccessV2Update,
		DeleteContext: resourcePolicyAccessV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFuncV2([]string{"ACCESS_POLICY", "GLOBAL_POLICY"}),
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
					"ALLOW",
					"DENY",
					"REQUIRE_APPROVAL",
				}, false),
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
			"policy_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_msg": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This is for providing a customer message for the user.",
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
										Computed:    true,
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
													Computed: true,
												},
												"lhs": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
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
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourcePolicyAccessV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := validateObjectTypeUniqueness(d); err != nil {
		return diag.FromErr(err)
	}

	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "ACCESS_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ACCESS_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandCreatePolicyRuleV2(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating zpa policy rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

	policysetcontrollerv2, _, err := policysetcontrollerv2.CreateRule(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policysetcontrollerv2.ID)

	return resourcePolicyAccessV2Read(ctx, d, meta)
}

func resourcePolicyAccessV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Set Terraform state
	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", v2PolicyRule.Name)
	_ = d.Set("description", v2PolicyRule.Description)
	_ = d.Set("action", v2PolicyRule.Action)
	_ = d.Set("operator", v2PolicyRule.Operator)
	_ = d.Set("policy_set_id", policySetID)
	_ = d.Set("custom_msg", v2PolicyRule.CustomMsg)
	_ = d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))
	_ = d.Set("app_server_groups", flattenCommonAppServerGroups(resp.AppServerGroups))
	_ = d.Set("app_connector_groups", flattenCommonAppConnectorGroups(resp.AppConnectorGroups))

	return nil
}

func resourcePolicyAccessV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := validateObjectTypeUniqueness(d); err != nil {
		return diag.FromErr(err)
	}

	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "ACCESS_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ACCESS_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating access policy rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyRuleV2(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

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

	return resourcePolicyAccessV2Read(ctx, d, meta)
}

func resourcePolicyAccessV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Assume "ACCESS_POLICY" is the policy type for this resource. Adjust as needed.
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ACCESS_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Deleting access policy set rule with id %v\n", d.Id())

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// func expandCreatePolicyRuleV2(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
// 	conditions, err := ExpandPolicyConditionsV2(d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &policysetcontrollerv2.PolicyRule{
// 		ID:                 d.Get("id").(string),
// 		Name:               d.Get("name").(string),
// 		Description:        d.Get("description").(string),
// 		Action:             d.Get("action").(string),
// 		CustomMsg:          d.Get("custom_msg").(string),
// 		Operator:           d.Get("operator").(string),
// 		PolicySetID:        policySetID,
// 		Conditions:         conditions,
// 		AppServerGroups:    expandCommonServerGroups(d),
// 		AppConnectorGroups: expandCommonAppConnectorGroups(d),

// 	}, nil
// }

func expandCreatePolicyRuleV2(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}
	rule := &policysetcontrollerv2.PolicyRule{
		ID:                 d.Get("id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Action:             d.Get("action").(string),
		CustomMsg:          d.Get("custom_msg").(string),
		Operator:           d.Get("operator").(string),
		PolicySetID:        policySetID,
		Conditions:         conditions,
		AppServerGroups:    expandCommonServerGroups(d),
		AppConnectorGroups: expandCommonAppConnectorGroups(d),
	}

	// Conditionally set credential if the user actually set it in TF.
	// ADDING CONDITION TO EXPLICITLY IGNORE THE CREDENTIAL ATTRIBUTE
	if val, ok := d.GetOk("credential"); ok {
		c := val.(map[string]interface{})
		rule.Credential = &policysetcontrollerv2.Credential{
			ID:   c["id"].(string),
			Name: c["name"].(string),
		}
	}

	return rule, nil
}
