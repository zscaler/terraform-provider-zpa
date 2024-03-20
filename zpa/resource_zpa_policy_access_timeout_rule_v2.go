package zpa

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontrollerv2"
)

func resourcePolicyTimeoutRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyTimeoutRuleV2Create,
		Read:   resourcePolicyTimeoutRuleV2Read,
		Update: resourcePolicyTimeoutRuleV2Update,
		Delete: resourcePolicyTimeoutRuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"TIMEOUT_POLICY", "REAUTH_POLICY"}),
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
					"RE_AUTH",
				}, false),
			},
			"custom_msg": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is for providing a customer message for the user.",
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reauth_idle_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reauth_timeout": {
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
											"IDP",
											"POSTURE",
											"PLATFORM",
											"SAML",
											"SCIM",
											"SCIM_GROUP",
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

func resourcePolicyTimeoutRuleV2Create(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	service := client.policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Automatically determining policy_set_id for "TIMEOUT_POLICY"
	policySetID, err := fetchPolicySetIDByType(client, "TIMEOUT_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	// Validate reauth_idle_timeout
	if idleTimeout, ok := d.GetOk("reauth_idle_timeout"); ok {
		if err := validateTimeoutIntervals(idleTimeout.(string)); err != nil {
			return err
		}
	}

	// Validate reauth_timeout
	if timeout, ok := d.GetOk("reauth_timeout"); ok {
		if err := validateTimeoutIntervals(timeout.(string)); err != nil {
			return err
		}
	}

	req, err := expandTimeOutPolicyRule(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy timeout rule with request\n%+v\n", req)

	policysetcontrollerv2, _, err := service.CreateRule(req)
	if err != nil {
		return err
	}
	d.SetId(policysetcontrollerv2.ID)

	return resourcePolicyTimeoutRuleV2Read(d, m)
}

func resourcePolicyTimeoutRuleV2Read(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(client, "TIMEOUT_POLICY", microTenantID)
	if err != nil {
		return err
	}

	service := client.policysetcontrollerv2.WithMicroTenant(microTenantID)
	resp, respErr, err := service.GetPolicyRule(policySetID, d.Id())
	if err != nil {
		// Adjust this error handling to match how your client library exposes HTTP response details
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	v2PolicyRule := policysetcontrollerv2.ConvertV1ResponseToV2Request(*resp)

	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", v2PolicyRule.Name)
	_ = d.Set("description", v2PolicyRule.Description)
	_ = d.Set("action", v2PolicyRule.Action)
	_ = d.Set("policy_set_id", policySetID)
	_ = d.Set("custom_msg", v2PolicyRule.CustomMsg)
	_ = d.Set("reauth_idle_timeout", secondsToHumanReadable(resp.ReauthIdleTimeout))
	_ = d.Set("reauth_timeout", secondsToHumanReadable(resp.ReauthTimeout))
	_ = d.Set("microtenant_id", v2PolicyRule.MicroTenantID)
	_ = d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))

	return nil
}

func resourcePolicyTimeoutRuleV2Update(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	service := client.policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Automatically determining policy_set_id for "TIMEOUT_POLICY"
	policySetID, err := fetchPolicySetIDByType(client, "TIMEOUT_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating policy timeout rule ID: %v\n", ruleID)
	req, err := expandTimeOutPolicyRule(d, policySetID)
	if err != nil {
		return err
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}
	// Validate reauth_idle_timeout
	if idleTimeout, ok := d.GetOk("reauth_idle_timeout"); ok {
		if err := validateTimeoutIntervals(idleTimeout.(string)); err != nil {
			return err
		}
	}

	// Validate reauth_timeout
	if timeout, ok := d.GetOk("reauth_timeout"); ok {
		if err := validateTimeoutIntervals(timeout.(string)); err != nil {
			return err
		}
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

	return resourcePolicyTimeoutRuleV2Read(d, m)
}

func resourcePolicyTimeoutRuleV2Delete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	// Assume "TIMEOUT_POLICY" is the policy type for this resource. Adjust as needed.
	policySetID, err := fetchPolicySetIDByType(client, "TIMEOUT_POLICY", microTenantID)
	if err != nil {
		return err
	}

	service := client.policysetcontrollerv2.WithMicroTenant(microTenantID)
	if _, err := service.Delete(policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandTimeOutPolicyRule(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	reauthIdleTimeoutInSeconds, err := parseHumanReadableTimeout(d.Get("reauth_idle_timeout").(string))
	if err != nil {
		return nil, err
	}

	reauthTimeoutInSeconds, err := parseHumanReadableTimeout(d.Get("reauth_timeout").(string))
	if err != nil {
		return nil, err
	}

	return &policysetcontrollerv2.PolicyRule{
		ID:                d.Get("id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		CustomMsg:         d.Get("custom_msg").(string),
		Action:            d.Get("action").(string),
		ReauthIdleTimeout: strconv.Itoa(reauthIdleTimeoutInSeconds),
		ReauthTimeout:     strconv.Itoa(reauthTimeoutInSeconds),
		PolicySetID:       policySetID,
		Conditions:        conditions,
	}, nil
}

func validateTimeoutIntervals(input string) error {
	// Allow "Never" without further checks
	if strings.ToLower(input) == "never" {
		return nil
	}

	timeoutInSeconds, err := parseHumanReadableTimeout(input)
	if err != nil {
		return err
	}

	// Ensure other time intervals meet the minimum requirement
	if timeoutInSeconds >= 0 && timeoutInSeconds < 600 { // 10 minutes in seconds
		return fmt.Errorf("timeout interval must be at least 10 minutes or 'Never'")
	}

	return nil
}
