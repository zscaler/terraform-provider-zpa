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

func resourcePolicyPortalAccessRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyPortalAccessRuleRuleCreate,
		ReadContext:   resourcePolicyPortalAccessRuleRuleRead,
		UpdateContext: resourcePolicyPortalAccessRuleRuleUpdate,
		DeleteContext: resourcePolicyPortalAccessRuleRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFuncV2([]string{"PRIVILEGED_PORTAL_POLICY"}),
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
					"CHECK_PRIVILEGED_PORTAL_CAPABILITIES",
				}, false),
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:        schema.TypeSet,
				Optional:    true,
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
										Description: "  This is for specifying the policy critiera.",
										ValidateFunc: validation.StringInSlice([]string{
											"COUNTRY_CODE",
											"PRIVILEGE_PORTAL",
											"SAML",
											"SCIM",
											"SCIM_GROUP",
										}, false),
									},
									"entry_values": {
										Type:     schema.TypeSet,
										Optional: true,
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
			"privileged_portal_capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_file": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Allows a User to delete files to reclaim space. Allowing deletion will prevent auditing of the file.",
						},
						"access_uninspected_file": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Allows a User like an Admin to see all files marked Uninspected from other users in the tenant.",
						},
						"request_approvals": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Allows a User to request approvals",
						},
						"review_approvals": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Allows a User to review approvals",
						},
					},
				},
			},
			"microtenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePolicyPortalAccessRuleRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}
	// Automatically determining policy_set_id for "PRIVILEGED_PORTAL_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "PRIVILEGED_PORTAL_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandPrivilegedPortalCapabilitiesRule(d, policySetID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Creating zpa policy privileged portal capabilities rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return diag.FromErr(err)
	}

	resp, _, err := policysetcontrollerv2.CreateRule(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourcePolicyPortalAccessRuleRuleRead(ctx, d, meta)
}

func resourcePolicyPortalAccessRuleRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "PRIVILEGED_PORTAL_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontrollerv2.GetPolicyRule(ctx, service, policySetID, d.Id())
	if err != nil {
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing privileged portal capabilities %s from state because it no longer exists in ZPA", d.Id())
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
	_ = d.Set("microtenant_id", v2PolicyRule.MicroTenantID)
	_ = d.Set("policy_set_id", policySetID) // Here, you're setting it based on fetched ID
	_ = d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))
	if len(resp.PrivilegedCapabilities.Capabilities) > 0 {
		_ = d.Set("privileged_portal_capabilities", flattenPrivilegedPortalCapabilities(resp.PrivilegedPortalCapabilities))
	}
	return nil
}

func resourcePolicyPortalAccessRuleRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "PRIVILEGED_PORTAL_POLICY"
	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "PRIVILEGED_PORTAL_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return diag.FromErr(err)
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating policy privileged portal capabilities rule ID: %v\n", ruleID)
	req, err := expandPrivilegedPortalCapabilitiesRule(d, policySetID)
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

	return resourcePolicyPortalAccessRuleRuleRead(ctx, d, meta)
}

func resourcePolicyPortalAccessRuleRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(ctx, zClient, "PRIVILEGED_PORTAL_POLICY", microTenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting privileged portal capabilities set rule with id %v\n", d.Id())

	if _, err := policysetcontrollerv2.Delete(ctx, service, policySetID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenPrivilegedPortalCapabilities(capabilities policysetcontrollerv2.PrivilegedPortalCapabilities) []interface{} {
	if len(capabilities.Capabilities) == 0 {
		return nil
	}

	capMap := make(map[string]bool)
	for _, cap := range capabilities.Capabilities {
		switch cap {
		case "DELETE_FILE":
			capMap["delete_file"] = true
		case "ACCESS_UNINSPECTED_FILE":
			capMap["access_uninspected_file"] = true
		case "REQUEST_APPROVALS":
			capMap["request_approvals"] = true
		case "REVIEW_APPROVALS":
			capMap["review_approvals"] = true
		}
	}

	return []interface{}{
		map[string]interface{}{
			"delete_file":             capMap["delete_file"],
			"access_uninspected_file": capMap["access_uninspected_file"],
			"request_approvals":       capMap["request_approvals"],
			"review_approvals":        capMap["review_approvals"],
		},
	}
}

func expandPrivilegedPortalCapabilitiesRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	// Initialize an empty slice for capabilities
	capabilities := []string{}

	// Check if the privileged_portal_capabilities block exists
	if v, ok := d.GetOk("privileged_portal_capabilities"); ok {
		privCapsList := v.([]interface{})
		if len(privCapsList) > 0 {
			// Extract the map from the first item of the list (MaxItems: 1)
			privCapsMap := privCapsList[0].(map[string]interface{})

			// Convert Boolean values to the API expected string values
			if privCapsMap["delete_file"].(bool) {
				capabilities = append(capabilities, "DELETE_FILE")
			}
			if privCapsMap["access_uninspected_file"].(bool) {
				capabilities = append(capabilities, "ACCESS_UNINSPECTED_FILE")
			}
			if privCapsMap["request_approvals"].(bool) {
				capabilities = append(capabilities, "REQUEST_APPROVALS")
			}
			if privCapsMap["review_approvals"].(bool) {
				capabilities = append(capabilities, "REVIEW_APPROVALS")
			}
		}
	}

	// Construct the PrivilegedPortalCapabilities struct
	privilegedPortalCapabilities := policysetcontrollerv2.PrivilegedPortalCapabilities{
		Capabilities: capabilities,
	}

	// Construct the PolicyRule struct
	policyRule := &policysetcontrollerv2.PolicyRule{
		ID:                           d.Get("id").(string),
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Action:                       d.Get("action").(string),
		MicroTenantID:                d.Get("microtenant_id").(string),
		PolicySetID:                  policySetID,
		Conditions:                   conditions,
		PrivilegedPortalCapabilities: privilegedPortalCapabilities,
	}

	return policyRule, nil
}
