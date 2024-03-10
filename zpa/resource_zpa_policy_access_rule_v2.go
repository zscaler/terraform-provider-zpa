package zpa

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
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
				Optional: true,
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
	// zClient := m.(*Client)
	service := m.(*Client).policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))
	req, err := expandCreatePolicyRuleV2(d)
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
	service := m.(*Client).policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return err
	}
	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", globalPolicySet.ID, d.Id())
	resp, _, err := service.GetPolicyRule(globalPolicySet.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	// Assuming you have a function to convert v1 response to v2 request format if necessary
	// If your data is already in the v2 format, you may not need this conversion
	v2PolicyRule := policysetcontrollerv2.ConvertV1ResponseToV2Request(*resp)

	// Set Terraform state
	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", v2PolicyRule.Name)
	_ = d.Set("description", v2PolicyRule.Description)
	_ = d.Set("action", v2PolicyRule.Action)
	_ = d.Set("operator", v2PolicyRule.Operator)
	_ = d.Set("policy_set_id", v2PolicyRule.PolicySetID)
	_ = d.Set("custom_msg", v2PolicyRule.CustomMsg)
	_ = d.Set("conditions", flattenConditionsV2(v2PolicyRule.Conditions))
	_ = d.Set("app_server_groups", flattenPolicyRuleServerGroupsV2(resp.AppServerGroups))
	_ = d.Set("app_connector_groups", flattenPolicyRuleAppConnectorGroupsV2(resp.AppConnectorGroups))

	return nil
}

func resourcePolicyAccessV2Update(d *schema.ResourceData, m interface{}) error {
	// zClient := m.(*Client)
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return err
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyRuleV2(d)
	if err != nil {
		return err
	}

	if err := ValidatePolicyRuleConditions(d); err != nil {
		return err
	}
	if _, _, err := service.GetPolicyRule(globalPolicySet.ID, ruleID); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	serviceUpdate := m.(*Client).policysetcontrollerv2.WithMicroTenant(GetString(d.Get("microtenant_id")))
	if _, err := serviceUpdate.UpdateRule(globalPolicySet.ID, ruleID, req); err != nil {
		return err
	}

	return resourcePolicyAccessV2Read(d, m)
}

func resourcePolicyAccessV2Delete(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	if _, err := service.Delete(globalPolicySet.ID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyRuleV2(d *schema.ResourceData) (*policysetcontrollerv2.PolicyRule, error) {
	policySetID, ok := d.Get("policy_set_id").(string)
	if !ok {
		log.Printf("[ERROR] policy_set_id is not set\n")
		return nil, fmt.Errorf("policy_set_id is not set")
	}
	log.Printf("[INFO] action_id:%v\n", d.Get("action_id"))
	conditions, err := ExpandPolicyConditionsV2(d)
	if err != nil {
		return nil, err
	}
	return &policysetcontrollerv2.PolicyRule{
		ID:          d.Get("id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Action:      d.Get("action").(string),
		CustomMsg:   d.Get("custom_msg").(string),
		Operator:    d.Get("operator").(string),
		PolicySetID: policySetID,
		// MicroTenantID:      d.Get("microtenant_id").(string),
		Conditions:         conditions,
		AppServerGroups:    expandPolicySetControllerAppServerGroupsV2(d),
		AppConnectorGroups: expandPolicysetControllerAppConnectorGroupsV2(d),
	}, nil
}

func ExpandPolicyConditionsV2(d *schema.ResourceData) ([]policysetcontrollerv2.PolicyRuleResourceConditions, error) {
	conditionInterface, ok := d.GetOk("conditions")
	if ok {
		conditions := conditionInterface.([]interface{})
		log.Printf("[INFO] conditions data: %+v\n", conditions)
		var conditionSets []policysetcontrollerv2.PolicyRuleResourceConditions
		for _, condition := range conditions {
			conditionSet, _ := condition.(map[string]interface{})
			if conditionSet != nil {
				operands, err := expandOperandsListV2(conditionSet["operands"])
				if err != nil {
					return nil, err
				}
				conditionSets = append(conditionSets, policysetcontrollerv2.PolicyRuleResourceConditions{
					ID:       conditionSet["id"].(string),
					Operator: conditionSet["operator"].(string),
					Operands: operands,
				})
			}
		}
		return conditionSets, nil
	}

	return []policysetcontrollerv2.PolicyRuleResourceConditions{}, nil
}

func expandOperandsListV2(ops interface{}) ([]policysetcontrollerv2.PolicyRuleResourceOperands, error) {
	if ops != nil {
		operands := ops.([]interface{})
		log.Printf("[INFO] operands data: %+v\n", operands)
		var operandsSets []policysetcontrollerv2.PolicyRuleResourceOperands
		for _, operand := range operands {
			operandSet, _ := operand.(map[string]interface{})
			id, _ := operandSet["id"].(string)
			IdpID, _ := operandSet["idp_id"].(string)

			// Extracting Values from TypeSet
			var values []string
			if valuesInterface, valuesOk := operandSet["values"].(*schema.Set); valuesOk && valuesInterface != nil {
				for _, v := range valuesInterface.List() {
					if strVal, ok := v.(string); ok {
						values = append(values, strVal)
					}
				}
			}

			// Extracting EntryValues
			var entryValues []policysetcontrollerv2.OperandsResourceLHSRHSValue
			if entryValuesInterface, ok := operandSet["entry_values"].([]interface{}); ok {
				for _, ev := range entryValuesInterface {
					entryValueMap, _ := ev.(map[string]interface{})
					lhs, _ := entryValueMap["lhs"].(string)
					rhs, _ := entryValueMap["rhs"].(string)

					entryValues = append(entryValues, policysetcontrollerv2.OperandsResourceLHSRHSValue{
						LHS: lhs,
						RHS: rhs,
					})
				}
			}

			log.Printf("[DEBUG] Extracted values: %+v\n", values)
			log.Printf("[DEBUG] Extracted entryValues: %+v\n", entryValues)

			op := policysetcontrollerv2.PolicyRuleResourceOperands{
				ID:                id,
				ObjectType:        operandSet["object_type"].(string),
				IDPID:             IdpID,
				Values:            values,
				EntryValuesLHSRHS: entryValues,
			}

			operandsSets = append(operandsSets, op)
		}

		return operandsSets, nil
	}
	return []policysetcontrollerv2.PolicyRuleResourceOperands{}, nil
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

// ValidatePolicyRuleConditions ensures that the necessary values are provided for specific object types.
func ValidatePolicyRuleConditions(d *schema.ResourceData) error {
	conditions, ok := d.GetOk("conditions")
	if !ok {
		// If conditions are not provided, there's nothing to validate
		return nil
	}

	validClientTypes := []string{
		"zpn_client_type_zapp",
		"zpn_client_type_exporter",
		"zpn_client_type_ip_anchoring",
		"zpn_client_type_browser_isolation",
		"zpn_client_type_machine_tunnel",
		"zpn_client_type_edge_connector",
		"zpn_client_type_exporter_noauth",
		"zpn_client_type_slogger",
		"zpn_client_type_branch_connector",
	}

	validPlatformTypes := []string{"mac", "linux", "ios", "windows", "android"}

	conditionList := conditions.([]interface{})
	for _, condition := range conditionList {
		conditionMap := condition.(map[string]interface{})
		operands, ok := conditionMap["operands"].([]interface{})
		if !ok {
			// No operands to validate
			continue
		}

		for _, operand := range operands {
			operandMap := operand.(map[string]interface{})
			objectType := operandMap["object_type"].(string)
			valuesSet, valuesPresent := operandMap["values"].(*schema.Set)

			switch objectType {
			case "APP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("an Application Segment ID must be provided when object_type = APP")
				}
			case "APP_GROUP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Segment Group ID must be provided when object_type = APP_GROUP")
				}
			case "MACHINE_GRP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Machine Group ID must be provided when object_type = MACHINE_GRP")
				}
			case "LOCATION":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Location ID must be provided when object_type = LOCATION")
				}
			case "EDGE_CONNECTOR_GROUP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Edge Connector Group ID must be provided when object_type = EDGE_CONNECTOR_GROUP")
				}
			case "BRANCH_CONNECTOR_GROUP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Branch Connector Group ID must be provided when object_type = BRANCH_CONNECTOR_GROUP")
				}
			case "CLIENT_TYPE":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("please provide one of the valid Client Types: %v", validClientTypes)
				}
				for _, v := range valuesSet.List() {
					value := v.(string)
					if !contains(validClientTypes, value) {
						return fmt.Errorf("invalid Client Type '%s'. Please provide one of the valid Client Types: %v", value, validClientTypes)
					}
				}
			case "PLATFORM":
				entryValues, ok := operandMap["entry_values"].([]interface{})
				if !ok || len(entryValues) == 0 {
					return fmt.Errorf("please provide one of the valid platform types: %v", validPlatformTypes)
				}
				for _, ev := range entryValues {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)
					if !lhsOk || !contains(validPlatformTypes, lhs) {
						return fmt.Errorf("please provide one of the valid platform types: %v", validPlatformTypes)
					}
					if !rhsOk || rhs != "true" {
						return fmt.Errorf("rhs value must be 'true' for PLATFORM object_type")
					}
				}
			case "POSTURE":
				entryValues, ok := operandMap["entry_values"].([]interface{})
				if !ok || len(entryValues) == 0 {
					return fmt.Errorf("please provide a valid Posture UDID")
				}
				for _, ev := range entryValues {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)
					if !lhsOk || !contains(validPlatformTypes, lhs) {
						return fmt.Errorf("please provide a valid Posture UDID")
					}
					if !rhsOk || (rhs != "true" && rhs != "false") {
						return fmt.Errorf("rhs value must be 'true' or 'false' for POSTURE object_type")
					}
				}
			case "TRUSTED_NETWORK":
				entryValues, ok := operandMap["entry_values"].([]interface{})
				if !ok || len(entryValues) == 0 {
					return fmt.Errorf("please provide a valid Network ID")
				}
				for _, ev := range entryValues {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)
					if !lhsOk || !contains(validPlatformTypes, lhs) {
						return fmt.Errorf("please provide a valid Network ID")
					}
					if !rhsOk || (rhs != "true" && rhs != "false") {
						return fmt.Errorf("rhs value must be 'true' or 'false' for TRUSTED_NETWORK object_type")
					}
				}
			case "COUNTRY_CODE":
				entryValues, ok := operandMap["entry_values"].([]interface{})
				if !ok || len(entryValues) == 0 {
					return fmt.Errorf("please provide a valid country code in 'entry_values'")
				}

				var invalidCodes []string
				for _, ev := range entryValues {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)

					// Validate 'lhs' as a country code
					if lhsOk {
						_, errors := validateCountryCode(lhs, "lhs")
						if len(errors) > 0 {
							// Collect invalid country codes instead of returning immediately
							invalidCodes = append(invalidCodes, lhs)
						}
					} else {
						return fmt.Errorf("a valid ISO-3166 Alpha-2 country code must be provided in 'lhs'")
					}

					// Ensure 'rhs' is "true"
					if !rhsOk || rhs != "true" {
						return fmt.Errorf("rhs value must be 'true' for COUNTRY_CODE object_type")
					}
				}

				// If there are any invalid country codes, return an aggregated error message
				if len(invalidCodes) > 0 {
					return fmt.Errorf("'%s' is not a valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", strings.Join(invalidCodes, "', '"))
				}
			case "SAML":
				entryValues, ok := operandMap["entry_values"].([]interface{})
				if !ok || len(entryValues) == 0 {
					return fmt.Errorf("entry_values must be provided for SAML object_type")
				}
				for _, ev := range entryValues {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs := evMap["rhs"].(string) // Directly accessing the value, assuming zero value (empty string) is not valid

					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid SAML attribute ID and cannot be empty for SAML object_type")
					}
					if rhs == "" {
						return fmt.Errorf("RHS must be a valid string i.e email adddress, department, group etc and cannot be empty for SAML object_type")
					}
				}
			case "SCIM":
				entryValues, ok := operandMap["entry_values"].([]interface{})
				if !ok || len(entryValues) == 0 {
					return fmt.Errorf("entry_values must be provided for SCIM object_type")
				}
				for _, ev := range entryValues {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs := evMap["rhs"].(string) // Directly accessing the value, assuming zero value (empty string) is not valid

					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid IdP ID and cannot be empty for SCIM object_type")
					}
					if rhs == "" {
						return fmt.Errorf("RHS must be a valid string i.e email adddress, department, group etc and cannot be empty for SCIM object_type")
					}
				}
			case "SCIM_GROUP":
				entryValues, ok := operandMap["entry_values"].([]interface{})
				if !ok || len(entryValues) == 0 {
					return fmt.Errorf("entry_values must be provided for SCIM_GROUP object_type")
				}
				for _, ev := range entryValues {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs := evMap["rhs"].(string) // Directly accessing the value, assuming zero value (empty string) is not valid

					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid IdP ID and cannot be empty for SCIM_GROUP object_type")
					}
					if rhs == "" {
						return fmt.Errorf("RHS must be a valid scim group ID and cannot be empty for SCIM_GROUP object_type")
					}
				}
			}

		}
	}
	return nil
}

// flattenConditions flattens the conditions part of the policy rule into a format suitable for Terraform schema.
func flattenConditionsV2(conditions []policysetcontrollerv2.PolicyRuleResourceConditions) []interface{} {
	if conditions == nil {
		return nil
	}

	c := make([]interface{}, len(conditions)) // Simplified slice initialization
	for i, condition := range conditions {
		condMap := make(map[string]interface{})
		condMap["operator"] = condition.Operator
		condMap["operands"] = flattenOperandsV2(condition.Operands)
		c[i] = condMap
	}
	return c
}

// flattenOperands flattens the operands part of the conditions into a format suitable for Terraform schema.
func flattenOperandsV2(operands []policysetcontrollerv2.PolicyRuleResourceOperands) []interface{} {
	if operands == nil {
		return nil
	}

	o := make([]interface{}, len(operands)) // Simplified slice initialization
	for i, operand := range operands {
		operandMap := make(map[string]interface{})
		operandMap["object_type"] = operand.ObjectType

		if len(operand.Values) > 0 {
			operandMap["values"] = operand.Values
		} else {
			operandMap["values"] = []interface{}{} // Ensure "values" key exists with an empty slice if no values are present.
		}

		entryValues := make([]interface{}, len(operand.EntryValuesLHSRHS))
		for j, entryValue := range operand.EntryValuesLHSRHS {
			entryValues[j] = map[string]interface{}{
				"lhs": entryValue.LHS,
				"rhs": entryValue.RHS,
			}
		}

		if len(entryValues) > 0 {
			operandMap["entry_values"] = entryValues
		} else {
			operandMap["entry_values"] = []interface{}{} // Ensure "entry_values" key exists with an empty slice if no entry values are present.
		}

		o[i] = operandMap
	}
	return o
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
