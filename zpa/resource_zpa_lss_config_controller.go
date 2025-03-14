package zpa

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

func getPolicyRuleResourceSchema() map[string]*schema.Schema {
	return MergeSchema(
		CommonPolicySchema(), map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "  This is for providing the rule action.",
				ValidateFunc: validation.StringInSlice([]string{
					"LOG",
				}, false),
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This is for proviidng the set of conditions for the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AND",
								"OR",
							}, false),
						},
						"operands": {
							Type:        schema.TypeList,
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
										Required:    true,
										Description: "  This is for specifying the policy critiera.",
										ValidateFunc: validation.StringInSlice([]string{
											"APP",
											"APP_GROUP",
											"CLIENT_TYPE",
											"IDP",
											"SCIM",
											"SCIM_GROUP",
											"SAML",
										}, false),
									},
									"entry_values": {
										Type:     schema.TypeList,
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
		},
	)
}

func resourceLSSConfigController() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLSSConfigControllerCreate,
		ReadContext:   resourceLSSConfigControllerRead,
		UpdateContext: resourceLSSConfigControllerUpdate,
		DeleteContext: resourceLSSConfigControllerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := lssconfigcontroller.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(resp.ID)
						_ = d.Set("id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_rule_resource": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: getPolicyRuleResourceSchema(),
				},
			},
			"connector_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "App Connector Group(s) to be added to the LSS configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_message": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the LSS configuration",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether this LSS configuration is enabled or not. Supported values: true, false",
						},
						"filter": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Filter for the LSS configuration. Format given by the following API to get status codes: /mgmtconfig/v2/admin/lssConfig/statusCodes",
						},

						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Format of the log type. Format given by the following API to get formats: /mgmtconfig/v2/admin/lssConfig/logType/formats",
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Description:  "Name of the LSS configuration",
						},
						"lss_host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Host of the LSS configuration",
						},
						"lss_port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Port of the LSS configuration",
						},
						"source_log_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log type of the LSS configuration",
							ValidateFunc: validation.StringInSlice([]string{
								"zpn_trans_log",
								"zpn_auth_log",
								"zpn_ast_auth_log",
								"zpn_http_trans_log",
								"zpn_audit_log",
								"zpn_ast_comprehensive_stats",
								"zpn_sys_auth_log",
								"zpn_waf_http_exchanges_log",
								"zpn_pbroker_comprehensive_stats",
							}, false),
						},
						"use_tls": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceLSSConfigControllerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandLSSResource(d)
	log.Printf("[INFO] Creating zpa lss config controller with request\n%+v\n", req)

	sourceLogType := d.Get("config.0.source_log_type").(string)

	conditions := d.Get("policy_rule_resource.0.conditions").([]interface{})
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		operandsInterface := conditionMap["operands"].([]interface{})

		var operands []lssconfigcontroller.PolicyRuleResourceOperands
		for _, operandInterface := range operandsInterface {
			operandMap := operandInterface.(map[string]interface{})
			objectType := operandMap["object_type"].(string)
			valuesInterface := operandMap["values"].(*schema.Set).List()

			var values []string
			for _, valueInterface := range valuesInterface {
				values = append(values, valueInterface.(string))
			}

			entryValuesInterface, exists := operandMap["entry_values"]
			var entryValues []lssconfigcontroller.OperandsResourceLHSRHSValue
			if exists && entryValuesInterface != nil {
				for _, entryValueInterface := range entryValuesInterface.([]interface{}) {
					entryValueMap := entryValueInterface.(map[string]interface{})
					entryValue := lssconfigcontroller.OperandsResourceLHSRHSValue{
						LHS: entryValueMap["lhs"].(string),
						RHS: entryValueMap["rhs"].(string),
					}
					entryValues = append(entryValues, entryValue)
				}
			}

			operand := lssconfigcontroller.PolicyRuleResourceOperands{
				ObjectType:                  objectType,
				Values:                      values,
				OperandsResourceLHSRHSValue: &entryValues,
			}
			operands = append(operands, operand)
		}

		// Validate operand object types and values here
		for _, operand := range operands {
			err := validateLSSConfigControllerFilters(sourceLogType, operand.ObjectType, "", operand.Values, operands)
			if err != nil {
				return diag.FromErr(err) // handle the error as appropriate
			}
		}
	}

	// Validate filters within the config block separately
	if filterSet, exists := d.GetOk("config.0.filter"); exists {
		for _, filter := range filterSet.(*schema.Set).List() {
			filterStr := filter.(string)
			// For filter validation, passing empty string as the objectType and nil for values and operands
			err := validateLSSConfigControllerFilters(sourceLogType, "", filterStr, nil, nil)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	resp, _, err := lssconfigcontroller.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created lss config controller request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceLSSConfigControllerRead(ctx, d, meta)
}

func resourceLSSConfigControllerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, _, err := lssconfigcontroller.Get(ctx, service, d.Id())
	if err != nil {
		if err.(*errorx.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing lss config controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting lss config controller:\n%+v\n", resp)
	d.SetId(resp.ID)
	if resp.PolicyRule != nil {
		_ = d.Set("policy_rule_id", resp.PolicyRule.ID)
	}
	_ = d.Set("config", flattenLSSConfig(resp.LSSConfig))
	_ = d.Set("connector_groups", flattenConnectorGroupsSimple(resp.ConnectorGroups))
	return nil
}

func resourceLSSConfigControllerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id := d.Id()
	req := expandLSSResource(d)
	log.Printf("[INFO] Updating zpa lss config controller with request\n%+v\n", req)

	sourceLogType := d.Get("config.0.source_log_type").(string)

	conditions := d.Get("policy_rule_resource.0.conditions").([]interface{})
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		operandsInterface := conditionMap["operands"].([]interface{})

		var operands []lssconfigcontroller.PolicyRuleResourceOperands
		for _, operandInterface := range operandsInterface {
			operandMap := operandInterface.(map[string]interface{})
			objectType := operandMap["object_type"].(string)
			valuesInterface := operandMap["values"].(*schema.Set).List()

			var values []string
			for _, valueInterface := range valuesInterface {
				values = append(values, valueInterface.(string))
			}

			entryValuesInterface, exists := operandMap["entry_values"]
			var entryValues []lssconfigcontroller.OperandsResourceLHSRHSValue
			if exists && entryValuesInterface != nil {
				for _, entryValueInterface := range entryValuesInterface.([]interface{}) {
					entryValueMap := entryValueInterface.(map[string]interface{})
					entryValue := lssconfigcontroller.OperandsResourceLHSRHSValue{
						LHS: entryValueMap["lhs"].(string),
						RHS: entryValueMap["rhs"].(string),
					}
					entryValues = append(entryValues, entryValue)
				}
			}

			operand := lssconfigcontroller.PolicyRuleResourceOperands{
				ObjectType:                  objectType,
				Values:                      values,
				OperandsResourceLHSRHSValue: &entryValues,
			}
			operands = append(operands, operand)
		}

		// Validate operand object types and values here
		for _, operand := range operands {
			err := validateLSSConfigControllerFilters(sourceLogType, operand.ObjectType, "", operand.Values, operands)
			if err != nil {
				return diag.FromErr(err) // handle the error as appropriate
			}
		}
	}

	// Validate filters within the config block separately
	if filterSet, exists := d.GetOk("config.0.filter"); exists {
		for _, filter := range filterSet.(*schema.Set).List() {
			filterStr := filter.(string)
			// Pass empty string as objectType and nil for values and operands when validating filters
			err := validateLSSConfigControllerFilters(sourceLogType, "", filterStr, nil, nil)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if _, _, err := lssconfigcontroller.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, err := lssconfigcontroller.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceLSSConfigControllerRead(ctx, d, meta)
}

func resourceLSSConfigControllerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Deleting lss config controller ID: %v\n", d.Id())

	if _, err := lssconfigcontroller.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] lss config controller deleted")
	return nil
}

func expandLSSResource(d *schema.ResourceData) lssconfigcontroller.LSSResource {
	policy, err := expandPolicyRuleResource(d)
	if err != nil {
		log.Printf("[ERROR] failed reading policy rule resource: %v\n", err)
	}
	req := lssconfigcontroller.LSSResource{
		ID:                 d.Get("id").(string),
		PolicyRuleResource: policy,
		LSSConfig:          expandLSSConfigController(d),
		ConnectorGroups:    expandConnectorGroups(d),
	}
	return req
}

func expandPolicyRuleResource(d *schema.ResourceData) (*lssconfigcontroller.PolicyRuleResource, error) {
	policyObj, ok := d.GetOk("policy_rule_resource")
	if !ok {
		return nil, nil
	}
	policyList := policyObj.([]interface{})
	if len(policyList) == 0 {
		return nil, nil
	}
	polictSet := policyList[0].(map[string]interface{})
	conditions, err := ExpandPolicyRuleResourceConditions(polictSet)
	if err != nil {
		return nil, err
	}

	return &lssconfigcontroller.PolicyRuleResource{
		Action:            polictSet["action"].(string),
		ActionID:          polictSet["action_id"].(string),
		CustomMsg:         polictSet["custom_msg"].(string),
		Description:       polictSet["description"].(string),
		ID:                polictSet["id"].(string),
		Name:              polictSet["name"].(string),
		Operator:          polictSet["operator"].(string),
		PolicyType:        polictSet["policy_type"].(string),
		Priority:          polictSet["priority"].(string),
		ReauthIdleTimeout: polictSet["reauth_idle_timeout"].(string),
		ReauthTimeout:     polictSet["reauth_timeout"].(string),
		RuleOrder:         polictSet["rule_order"].(string),
		Conditions:        conditions,
	}, nil
}

func ExpandPolicyRuleResourceConditions(d map[string]interface{}) ([]lssconfigcontroller.PolicyRuleResourceConditions, error) {
	conditionInterface, ok := d["conditions"]
	if ok {
		conditions := conditionInterface.([]interface{})
		log.Printf("[INFO] conditions data: %+v\n", conditions)
		var conditionSets []lssconfigcontroller.PolicyRuleResourceConditions
		for _, condition := range conditions {
			conditionSet, _ := condition.(map[string]interface{})
			if conditionSet != nil {
				operands, err := expandPolicyRuleResourceOperandsList(conditionSet["operands"])
				if err != nil {
					return nil, err
				}
				conditionSets = append(conditionSets, lssconfigcontroller.PolicyRuleResourceConditions{
					Operator: conditionSet["operator"].(string),
					Operands: &operands,
				})
			}
		}
		return conditionSets, nil
	}

	return []lssconfigcontroller.PolicyRuleResourceConditions{}, nil
}

func expandPolicyRuleResourceOperandsList(ops interface{}) ([]lssconfigcontroller.PolicyRuleResourceOperands, error) {
	if ops != nil {
		operands := ops.([]interface{})
		log.Printf("[INFO] operands data: %+v\n", operands)
		var operandsSets []lssconfigcontroller.PolicyRuleResourceOperands
		for _, operand := range operands {
			operandSet, _ := operand.(map[string]interface{})
			valuesSet := operandSet["values"].(*schema.Set)

			// Expanding entryValues from schema
			entryValuesInterface, exists := operandSet["entry_values"]
			var entryValues []lssconfigcontroller.OperandsResourceLHSRHSValue
			if exists && entryValuesInterface != nil {
				for _, entryValueInterface := range entryValuesInterface.([]interface{}) {
					entryValueMap := entryValueInterface.(map[string]interface{})
					entryValue := lssconfigcontroller.OperandsResourceLHSRHSValue{
						LHS: entryValueMap["lhs"].(string),
						RHS: entryValueMap["rhs"].(string),
					}
					entryValues = append(entryValues, entryValue)
				}
			}

			op := lssconfigcontroller.PolicyRuleResourceOperands{
				Values:     SetToStringSlice(valuesSet),
				ObjectType: operandSet["object_type"].(string),
			}

			if len(entryValues) > 0 {
				op.OperandsResourceLHSRHSValue = &entryValues // Setting expanded entryValues only if it is not empty
			}

			operandsSets = append(operandsSets, op)
		}
		return operandsSets, nil
	}
	return []lssconfigcontroller.PolicyRuleResourceOperands{}, nil
}

func expandLSSConfigController(d *schema.ResourceData) *lssconfigcontroller.LSSConfig {
	configInterface, ok := d.GetOk("config")
	if ok {
		configList := configInterface.([]interface{})
		if len(configList) == 0 {
			return nil
		}
		config, _ := configList[0].(map[string]interface{})
		filterSet, _ := config["filter"].(*schema.Set)
		return &lssconfigcontroller.LSSConfig{
			ID:            d.Get("id").(string),
			AuditMessage:  config["audit_message"].(string),
			Description:   config["description"].(string),
			Enabled:       config["enabled"].(bool),
			Filter:        SetToStringSlice(filterSet),
			Format:        config["format"].(string),
			Name:          config["name"].(string),
			LSSHost:       config["lss_host"].(string),
			LSSPort:       config["lss_port"].(string),
			SourceLogType: config["source_log_type"].(string),
			UseTLS:        config["use_tls"].(bool),
		}
	}
	return nil
}

func expandConnectorGroups(d *schema.ResourceData) []lssconfigcontroller.ConnectorGroups {
	appConnectorGroupsInterface, ok := d.GetOk("connector_groups")
	if ok {
		appConnector := appConnectorGroupsInterface.(*schema.Set)
		log.Printf("[INFO] connector groups data: %+v\n", appConnector)
		var appConnectorGroups []lssconfigcontroller.ConnectorGroups
		for _, appConnectorGroup := range appConnector.List() {
			appConnectorGroup, ok := appConnectorGroup.(map[string]interface{})
			if ok {
				for _, id := range appConnectorGroup["id"].([]interface{}) {
					appConnectorGroups = append(appConnectorGroups, lssconfigcontroller.ConnectorGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appConnectorGroups
	}

	return []lssconfigcontroller.ConnectorGroups{}
}

func flattenConnectorGroupsSimple(lssConnectorGroup []lssconfigcontroller.ConnectorGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(lssConnectorGroup))
	for i, item := range lssConnectorGroup {
		ids[i] = item.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}
