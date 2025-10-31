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

func resourcePolicyIsolationRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyIsolationRuleV2Create,
		ReadContext:   resourcePolicyIsolationRuleV2Read,
		UpdateContext: resourcePolicyIsolationRuleV2Update,
		DeleteContext: resourcePolicyIsolationRuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFuncV2([]string{"ISOLATION_POLICY"}),
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
					"ISOLATE",
					"BYPASS_ISOLATE",
				}, false),
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_isolation_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
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
											"EDGE_CONNECTOR_GROUP",
											"PLATFORM",
											"IDP",
											"SAML",
											"SCIM",
											"SCIM_GROUP",
											"CHROME_ENTERPRISE",
											"CHROME_POSTURE_PROFILE",
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

func resourcePolicyIsolationRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "ISOLATION_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ISOLATION_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandPolicyIsolationRule(d, policySetID) // ensure this function now accepts policySetID as a parameter
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating zpa policy isolation rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

	resp, _, err := policysetcontrollerv2.CreateRule(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourcePolicyIsolationRuleV2Read(ctx, d, meta)
}

func resourcePolicyIsolationRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ISOLATION_POLICY", microTenantID)
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
	d.SetId(resp.ID)
	d.Set("name", v2PolicyRule.Name)
	d.Set("description", v2PolicyRule.Description)
	d.Set("action", v2PolicyRule.Action)
	d.Set("policy_set_id", policySetID) // Here, you're setting it based on fetched ID
	d.Set("zpn_isolation_profile_id", v2PolicyRule.ZpnIsolationProfileID)
	d.Set("microtenant_id", v2PolicyRule.MicroTenantID)
	d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))

	return nil
}

func resourcePolicyIsolationRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "ISOLATION_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ISOLATION_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating policy isolation rule ID: %v\n", ruleID)
	req, err := expandPolicyIsolationRule(d, policySetID) // Adjusted to use the fetched policySetID
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

	if _, err := policysetcontrollerv2.UpdateRule(ctx, service, policySetID, ruleID, req); err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyIsolationRuleV2Read(ctx, d, meta)
}

func resourcePolicyIsolationRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Assume "ISOLATION_POLICY" is the policy type for this resource. Adjust as needed.
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "ISOLATION_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandPolicyIsolationRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	return &policysetcontrollerv2.PolicyRule{
		ID:                    d.Get("id").(string),
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Action:                d.Get("action").(string),
		ZpnIsolationProfileID: d.Get("zpn_isolation_profile_id").(string),
		MicroTenantID:         GetString(d.Get("microtenant_id")),
		PolicySetID:           policySetID,
		Conditions:            conditions,
	}, nil
}
