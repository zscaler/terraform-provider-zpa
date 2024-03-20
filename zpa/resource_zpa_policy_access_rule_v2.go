package zpa

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontrollerv2"
)

func resourcePolicyAccessRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyAccessV2Create,
		Read:   resourcePolicyAccessV2Read,
		Update: resourcePolicyAccessV2Update,
		Delete: resourcePolicyAccessV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"ACCESS_POLICY", "GLOBAL_POLICY"}),
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
				Description: "This is for providing a customer message for the user.",
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
										Type:     schema.TypeSet,
										Optional: true,
										//Computed:    true,
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
													// Computed: true,
												},
												"lhs": {
													Type:     schema.TypeString,
													Optional: true,
													//Computed: true,
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

func resourcePolicyAccessV2Create(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	service := client.policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Automatically determining policy_set_id for "ACCESS_POLICY"
	policySetID, err := fetchPolicySetIDByType(client, "ACCESS_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	req, err := expandCreatePolicyRuleV2(d, policySetID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy rule with request\n%+v\n", req)

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}

	policysetcontrollerv2, _, err := service.CreateRule(req)
	if err != nil {
		return err
	}
	d.SetId(policysetcontrollerv2.ID)

	return resourcePolicyAccessV2Read(d, m)
}

func resourcePolicyAccessV2Read(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(client, "ACCESS_POLICY", microTenantID)
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

	// Set Terraform state
	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", v2PolicyRule.Name)
	_ = d.Set("description", v2PolicyRule.Description)
	_ = d.Set("action", v2PolicyRule.Action)
	_ = d.Set("operator", v2PolicyRule.Operator)
	_ = d.Set("policy_set_id", policySetID) // Here, you're setting it based on fetched ID
	_ = d.Set("custom_msg", v2PolicyRule.CustomMsg)
	_ = d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))
	_ = d.Set("app_server_groups", flattenPolicyRuleServerGroupsV2(resp.AppServerGroups))
	_ = d.Set("app_connector_groups", flattenPolicyRuleAppConnectorGroupsV2(resp.AppConnectorGroups))

	return nil
}

func resourcePolicyAccessV2Update(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	service := client.policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))

	// Automatically determining policy_set_id for "ACCESS_POLICY"
	policySetID, err := fetchPolicySetIDByType(client, "ACCESS_POLICY", GetString(d.Get("microtenant_id")))
	if err != nil {
		return err
	}

	// Setting the policy_set_id for further use
	d.Set("policy_set_id", policySetID)

	ruleID := d.Id()
	log.Printf("[INFO] Updating access policy rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyRuleV2(d, policySetID)
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

	return resourcePolicyAccessV2Read(d, m)
}

func resourcePolicyAccessV2Delete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	// Assume "ACCESS_POLICY" is the policy type for this resource. Adjust as needed.
	policySetID, err := fetchPolicySetIDByType(client, "ACCESS_POLICY", microTenantID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting access policy set rule with id %v\n", d.Id())

	service := client.policysetcontrollerv2.WithMicroTenant(microTenantID)
	if _, err := service.Delete(policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyRuleV2(d *schema.ResourceData, policySetID string) (*policysetcontrollerv2.PolicyRule, error) {

	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}

	return &policysetcontrollerv2.PolicyRule{
		ID:                 d.Get("id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Action:             d.Get("action").(string),
		CustomMsg:          d.Get("custom_msg").(string),
		Operator:           d.Get("operator").(string),
		PolicySetID:        policySetID,
		Conditions:         conditions,
		AppServerGroups:    expandPolicySetControllerAppServerGroupsV2(d),
		AppConnectorGroups: expandPolicysetControllerAppConnectorGroupsV2(d),
	}, nil
}

func expandPolicySetControllerAppServerGroupsV2(d *schema.ResourceData) []policysetcontrollerv2.AppServerGroups {
	appServerGroupsInterface, ok := d.GetOk("app_server_groups")
	if ok {
		appServer := appServerGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", appServer)
		var appServerGroups []policysetcontrollerv2.AppServerGroups
		for _, appServerGroup := range appServer.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].(*schema.Set).List() {
					appServerGroups = append(appServerGroups, policysetcontrollerv2.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appServerGroups
	}

	return []policysetcontrollerv2.AppServerGroups{}
}

func expandPolicysetControllerAppConnectorGroupsV2(d *schema.ResourceData) []policysetcontrollerv2.AppConnectorGroups {
	appConnectorGroupsInterface, ok := d.GetOk("app_connector_groups")
	if ok {
		appConnector := appConnectorGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app connector groups data: %+v\n", appConnector)
		var appConnectorGroups []policysetcontrollerv2.AppConnectorGroups
		for _, appConnectorGroup := range appConnector.List() {
			appConnectorGroup, _ := appConnectorGroup.(map[string]interface{})
			if appConnectorGroup != nil {
				for _, id := range appConnectorGroup["id"].(*schema.Set).List() {
					appConnectorGroups = append(appConnectorGroups, policysetcontrollerv2.AppConnectorGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appConnectorGroups
	}

	return []policysetcontrollerv2.AppConnectorGroups{}
}

func flattenPolicyRuleServerGroupsV2(appServerGroup []policysetcontrollerv2.AppServerGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(appServerGroup))
	for i, serverGroup := range appServerGroup {
		ids[i] = serverGroup.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenPolicyRuleAppConnectorGroupsV2(appConnectorGroups []policysetcontrollerv2.AppConnectorGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(appConnectorGroups))
	for i, group := range appConnectorGroups {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}
