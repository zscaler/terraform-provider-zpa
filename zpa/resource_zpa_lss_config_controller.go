package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/lssconfigcontroller"
)

func getPolicyRuleResourceSchema() map[string]*schema.Schema {
	return MergeSchema(
		CommonPolicySchema(), map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "  This is for providing the rule action.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"DENY",
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
						"negated": {
							Type:     schema.TypeBool,
							Optional: true,
						},
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
										}, false),
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
		Create: resourceLSSConfigControllerCreate,
		Read:   resourceLSSConfigControllerRead,
		Update: resourceLSSConfigControllerUpdate,
		Delete: resourceLSSConfigControllerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("id", id)
				} else {
					resp, _, err := zClient.lssconfigcontroller.GetByName(id)
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
								"zpn_sys_auth_log",
								"zpn_http_insp",
								"zpn_ast_comprehensive_stats",
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

func resourceLSSConfigControllerCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandLSSResource(d)
	log.Printf("[INFO] Creating zpa lss config controller with request\n%+v\n", req)

	resp, _, err := zClient.lssconfigcontroller.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created lss config controller request. ID: %v\n", resp)
	d.SetId(resp.ID)

	return resourceLSSConfigControllerRead(d, m)
}

func resourceLSSConfigControllerRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.lssconfigcontroller.Get(d.Id())
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing lss config controller %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
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

func resourceLSSConfigControllerUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating lss config controller ID: %v\n", id)
	req := expandLSSResource(d)

	if _, err := zClient.lssconfigcontroller.Update(id, &req); err != nil {
		return err
	}

	return resourceLSSConfigControllerRead(d, m)
}

func resourceLSSConfigControllerDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting lss config controller ID: %v\n", d.Id())

	if _, err := zClient.lssconfigcontroller.Delete(d.Id()); err != nil {
		return err
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
		DefaultRule:       polictSet["default_rule"].(bool),
		Description:       polictSet["description"].(string),
		ID:                polictSet["id"].(string),
		Name:              polictSet["name"].(string),
		Operator:          polictSet["operator"].(string),
		PolicyType:        polictSet["policy_type"].(string),
		Priority:          polictSet["priority"].(string),
		ReauthDefaultRule: polictSet["reauth_default_rule"].(bool),
		ReauthIdleTimeout: polictSet["reauth_idle_timeout"].(string),
		ReauthTimeout:     polictSet["reauth_timeout"].(string),
		RuleOrder:         polictSet["rule_order"].(string),
		LssDefaultRule:    polictSet["lss_default_rule"].(bool),
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
					Negated:  conditionSet["negated"].(bool),
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
			op := lssconfigcontroller.PolicyRuleResourceOperands{
				Values:     SetToStringSlice(valuesSet),
				ObjectType: operandSet["object_type"].(string),
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
