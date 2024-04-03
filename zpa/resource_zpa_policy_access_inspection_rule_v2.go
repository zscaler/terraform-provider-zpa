package zpa

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontrollerv2"
)

func resourcePolicyInspectionRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyInspectionRuleV2Create,
		Read:   resourcePolicyInspectionRuleV2Read,
		Update: resourcePolicyInspectionRuleV2Update,
		Delete: resourcePolicyInspectionRuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFuncV2([]string{"INSPECTION_POLICY"}),
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
				Description: "This is for providing the rule action.",
				ValidateFunc: validation.StringInSlice([]string{
					"INSPECT",
					"BYPASS_INSPECT",
				}, false),
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zpn_inspection_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"conditions": {
				Type:        schema.TypeList,
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
							Type:        schema.TypeList,
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
											"IDP",
											"POSTURE",
											"SAML",
											"SCIM",
											"SCIM_GROUP",
											"TRUSTED_NETWORK",
										}, false),
									},
									"entry_values": {
										Type:     schema.TypeList,
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

func resourcePolicyInspectionRuleV2Create(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	service := client.policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Automatically determining policy_set_id for "INSPECTION_POLICY"
	policySetID, err := fetchPolicySetIDByType(client, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandPolicyInspectionRule(d, policySetID) // ensure this function now accepts policySetID as a parameter
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy inspection rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	resp, _, err := service.CreateRule(req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)
	return resourcePolicyInspectionRuleV2Read(d, m)
}

func resourcePolicyInspectionRuleV2Read(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(client, "INSPECTION_POLICY", microTenantID)
	if err != nil {
		return err
	}

	service := client.policysetcontrollerv2.WithMicroTenant(microTenantID)
	log.Printf("[INFO] Getting Policy Set Rule: policySetID:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := service.GetPolicyRule(policySetID, d.Id())
	if err != nil {
		// Adjust this error handling to match how your client library exposes HTTP response details
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	v2PolicyRule := policysetcontrollerv2.ConvertV1ResponseToV2Request(*resp)
	d.SetId(resp.ID)
	d.Set("name", v2PolicyRule.Name)
	d.Set("description", v2PolicyRule.Description)
	d.Set("action", v2PolicyRule.Action)
	d.Set("policy_set_id", policySetID) // Here, you're setting it based on fetched ID
	d.Set("zpn_inspection_profile_id", v2PolicyRule.ZpnInspectionProfileID)
	d.Set("microtenant_id", v2PolicyRule.MicroTenantID)
	d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))

	return nil
}

func resourcePolicyInspectionRuleV2Update(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	service := client.policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Automatically determining policy_set_id for "INSPECTION_POLICY"
	policySetID, err := fetchPolicySetIDByType(client, "INSPECTION_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating policy inspection rule ID: %v\n", ruleID)
	req, err := expandPolicyInspectionRule(d, policySetID) // Adjusted to use the fetched policySetID
	if err != nil {
		return err
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	// Checking the current state of the rule to handle cases where it might have been deleted outside Terraform
	_, respErr, err := service.GetPolicyRule(policySetID, ruleID)
	if err != nil {
		if respErr != nil && (respErr.StatusCode == http.StatusNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	_, err = service.UpdateRule(policySetID, ruleID, req)
	if err != nil {
		return err
	}

	return resourcePolicyInspectionRuleV2Read(d, m)
}

func resourcePolicyInspectionRuleV2Delete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	// Assume "INSPECTION_POLICY" is the policy type for this resource. Adjust as needed.
	policySetID, err := fetchPolicySetIDByType(client, "INSPECTION_POLICY", microTenantID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	service := client.policysetcontrollerv2.WithMicroTenant(microTenantID)
	if _, err := service.Delete(policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandPolicyInspectionRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	return &policysetcontrollerv2.PolicyRule{
		ID:                     d.Get("id").(string),
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		Action:                 d.Get("action").(string),
		ZpnInspectionProfileID: d.Get("zpn_inspection_profile_id").(string),
		MicroTenantID:          GetString(d.Get("microtenant_id")),
		PolicySetID:            policySetID,
		Conditions:             conditions,
	}, nil
}
