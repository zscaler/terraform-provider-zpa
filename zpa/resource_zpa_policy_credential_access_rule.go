package zpa

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontrollerv2"
)

func resourcePolicyCredentialAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCredentialAccessRuleCreate,
		Read:   resourcePolicyCredentialAccessRuleRead,
		Update: resourcePolicyCredentialAccessRuleUpdate,
		Delete: resourcePolicyCredentialAccessRuleDelete,
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
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

func resourcePolicyCredentialAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetControllerV2

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	// Automatically determining policy_set_id for "CREDENTIAL_POLICY"
	policySetID, err := fetchPolicySetIDByType(zClient, "CREDENTIAL_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandCredentialPolicyRule(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy credential rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	resp, _, err := policysetcontrollerv2.CreateRule(service, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePolicyCredentialAccessRuleRead(d, meta)
}

func resourcePolicyCredentialAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetControllerV2

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(zClient, "CREDENTIAL_POLICY", microTenantID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontrollerv2.GetPolicyRule(service, policySetID, d.Id())
	if err != nil {
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
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

func resourcePolicyCredentialAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetControllerV2

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(zClient, "CREDENTIAL_POLICY", microTenantID)
	if err != nil {
		return err
	}

	d.Set("policy_set_id", policySetID)
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy credential rule ID: %v\n", ruleID)
	req, err := expandCredentialPolicyRule(d, policySetID)
	if err != nil {
		return err
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	_, respErr, err := policysetcontrollerv2.GetPolicyRule(service, policySetID, ruleID)
	if err != nil {
		if respErr != nil && (respErr.StatusCode == http.StatusNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	if _, err := policysetcontrollerv2.UpdateRule(service, policySetID, ruleID, req); err != nil {
		return err
	}
	return resourcePolicyCredentialAccessRuleRead(d, meta)
}

func resourcePolicyCredentialAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetControllerV2

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	policySetID, err := fetchPolicySetIDByType(zClient, "CREDENTIAL_POLICY", microTenantID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy credential rule with id %v\n", d.Id())

	if _, err := policysetcontrollerv2.Delete(service, policySetID, d.Id()); err != nil {
		return fmt.Errorf("failed to delete policy credential rule: %w", err)
	}

	return nil
}

func flattenCredential(credential *policysetcontrollerv2.Credential) []interface{} {
	if credential == nil || credential.ID == "" {
		return []interface{}{}
	}

	m := make(map[string]interface{})
	m["id"] = credential.ID

	return []interface{}{m}
}

func expandCredentialPolicyRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}
	credential := expandCredential(d)

	return &policysetcontrollerv2.PolicyRule{
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Action:        d.Get("action").(string),
		MicroTenantID: d.Get("microtenant_id").(string),
		PolicySetID:   policySetID,
		Conditions:    conditions,
		Credential:    credential,
	}, nil
}

func expandCredential(d *schema.ResourceData) *policysetcontrollerv2.Credential {
	if v, ok := d.GetOk("credential"); ok && len(v.([]interface{})) > 0 {
		credentialMap := v.([]interface{})[0].(map[string]interface{})
		return &policysetcontrollerv2.Credential{
			ID: credentialMap["id"].(string),
		}
	}
	return nil
}
