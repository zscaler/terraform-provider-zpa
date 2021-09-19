package zpa

import (
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/policysetrule"
)

type listrules struct {
	orders map[string]int
	sync.Mutex
}

var rules = listrules{
	orders: make(map[string]int),
}

func resourcePolicyAccessRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourcePolicySetCreate,
		Read:     resourcePolicySetRead,
		Update:   resourcePolicySetUpdate,
		Delete:   resourcePolicySetDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "  This is for providing the rule action.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"DENY",
				}, false),
			},
			"action_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bypass_default_rule": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"custom_msg": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is for providing a customer message for the user.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the description of the access policy rule.",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the name of the policy rule.",
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AND",
					"OR",
				}, false),
			},
			"policy_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reauth_default_rule": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reauth_idle_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reauth_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_server_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of the server group IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
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
				Description: "List of app-connector IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This is for proviidng the set of conditions for the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"idp_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"lhs": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "This signifies the key for the object type. String ID example: id ",
									},
									"rhs": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "This denotes the value for the given object type. Its value depends upon the key.",
									},
									"object_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "  This is for specifying the policy critiera.",
										ValidateFunc: validation.StringInSlice([]string{
											"USER",
											"USER_GROUP",
											"LOCATION",
											"APP",
											"APP_GROUP",
											"SAML",
											"POSTURE",
											"CLIENT_TYPE",
											"IDP",
											"TRUSTED_NETWORK",
											"EDGE_CONNECTOR_GROUP",
											"MACHINE_GRP",
											"SCIM",
											"SCIM_GROUP",
										}, false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourcePolicySetCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandCreatePolicyRule(d)
	log.Printf("[INFO] Creating zpa policy rule with request\n%+v\n", req)
	if ValidateConditions(req.Conditions, zClient) {
		policysetrule, _, err := zClient.policysetrule.Create(&req)
		if err != nil {
			return err
		}
		d.SetId(policysetrule.ID)
		order, ok := d.GetOk("rule_order")
		if ok {
			reorder(order, policysetrule.PolicySetID, policysetrule.ID, zClient)
		}
		return resourcePolicySetRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy rule (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}
}

func resourcePolicySetRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	globalPolicySet, _, err := zClient.policysetglobal.Get()
	if err != nil {
		return err
	}
	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", globalPolicySet.ID, d.Id())
	resp, _, err := zClient.policysetrule.Get(globalPolicySet.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("action", resp.Action)
	_ = d.Set("action_id", resp.ActionID)
	_ = d.Set("custom_msg", resp.CustomMsg)
	_ = d.Set("description", resp.Description)
	_ = d.Set("name", resp.Name)
	_ = d.Set("bypass_default_rule", resp.BypassDefaultRule)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("priority", resp.Priority)
	_ = d.Set("reauth_default_rule", resp.ReauthDefaultRule)
	_ = d.Set("reauth_idle_timeout", resp.ReauthIdleTimeout)
	_ = d.Set("reauth_timeout", resp.ReauthTimeout)
	_ = d.Set("rule_order", resp.RuleOrder)
	_ = d.Set("conditions", flattenPolicyRuleConditions(resp.Conditions))
	_ = d.Set("app_server_groups", flattenPolicyRuleServerGroups(resp.AppServerGroups))
	_ = d.Set("app_connector_groups", flattenPolicyRuleAppConnectorGroups(resp.AppConnectorGroups))

	return nil
}

func resourcePolicySetUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicySet, _, err := zClient.policysetglobal.Get()
	if err != nil {
		return err
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy rule ID: %v\n", ruleID)
	req := expandCreatePolicyRule(d)
	if ValidateConditions(req.Conditions, zClient) {
		if _, err := zClient.policysetrule.Update(globalPolicySet.ID, ruleID, &req); err != nil {
			return err
		}
		if d.HasChange("rule_order") {
			order, ok := d.GetOk("rule_order")
			if ok {
				reorder(order, globalPolicySet.ID, ruleID, zClient)
			}
		}
		return resourcePolicySetRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy rule (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}

}

func resourcePolicySetDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicySet, _, err := zClient.policysetglobal.Get()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	if _, err := zClient.policysetrule.Delete(globalPolicySet.ID, d.Id()); err != nil {
		return err
	}

	return nil

}

func expandCreatePolicyRule(d *schema.ResourceData) policysetrule.PolicyRule {
	policySetID, ok := d.Get("policy_set_id").(string)
	if !ok {
		log.Printf("[ERROR] policy_set_id is not set\n")
	}
	log.Printf("[INFO] action_id:%v\n", d.Get("action_id"))
	return policysetrule.PolicyRule{
		Action:             d.Get("action").(string),
		ActionID:           d.Get("action_id").(string),
		CustomMsg:          d.Get("custom_msg").(string),
		Description:        d.Get("description").(string),
		ID:                 d.Get("id").(string),
		Name:               d.Get("name").(string),
		Operator:           d.Get("operator").(string),
		PolicySetID:        policySetID,
		PolicyType:         d.Get("policy_type").(string),
		Priority:           d.Get("priority").(string),
		ReauthDefaultRule:  d.Get("reauth_default_rule").(bool),
		ReauthIdleTimeout:  d.Get("reauth_idle_timeout").(string),
		ReauthTimeout:      d.Get("reauth_timeout").(string),
		RuleOrder:          d.Get("rule_order").(string),
		Conditions:         expandConditionSet(d),
		AppServerGroups:    expandPolicySetRuleAppServerGroups(d),
		AppConnectorGroups: expandPolicySetRuleAppConnectorGroups(d),
	}
}

func expandPolicySetRuleAppServerGroups(d *schema.ResourceData) []policysetrule.AppServerGroups {
	appServerGroupsInterface, ok := d.GetOk("app_server_groups")
	if ok {
		appServer := appServerGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", appServer)
		var appServerGroups []policysetrule.AppServerGroups
		for _, appServerGroup := range appServer.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].([]interface{}) {
					appServerGroups = append(appServerGroups, policysetrule.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appServerGroups
	}

	return []policysetrule.AppServerGroups{}
}

func expandPolicySetRuleAppConnectorGroups(d *schema.ResourceData) []policysetrule.AppConnectorGroups {
	appConnectorGroupsInterface, ok := d.GetOk("app_connector_groups")
	if ok {
		appConnector := appConnectorGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app connector groups data: %+v\n", appConnector)
		var appConnectorGroups []policysetrule.AppConnectorGroups
		for _, appConnectorGroup := range appConnector.List() {
			appConnectorGroup, _ := appConnectorGroup.(map[string]interface{})
			if appConnectorGroup != nil {
				for _, id := range appConnectorGroup["id"].([]interface{}) {
					appConnectorGroups = append(appConnectorGroups, policysetrule.AppConnectorGroups{
						ID: id.(string),
					})
				}

			}
		}
		return appConnectorGroups
	}

	return []policysetrule.AppConnectorGroups{}
}

func expandConditionSet(d *schema.ResourceData) []policysetrule.Conditions {
	conditionInterface, ok := d.GetOk("conditions")
	if ok {
		conditions := conditionInterface.([]interface{})
		log.Printf("[INFO] conditions data: %+v\n", conditions)
		var conditionSets []policysetrule.Conditions
		for _, condition := range conditions {
			conditionSet, _ := condition.(map[string]interface{})
			if conditionSet != nil {
				conditionSets = append(conditionSets, policysetrule.Conditions{
					ID:       conditionSet["id"].(string),
					Negated:  conditionSet["negated"].(bool),
					Operator: conditionSet["operator"].(string),
					Operands: expandOperandsList(conditionSet["operands"]),
				})
			}
		}
		return conditionSets
	}

	return []policysetrule.Conditions{}
}

func expandOperandsList(ops interface{}) []policysetrule.Operands {
	if ops != nil {
		operands := ops.([]interface{})
		log.Printf("[INFO] operands data: %+v\n", operands)
		var operandsSets []policysetrule.Operands
		for _, operand := range operands {
			operandSet, _ := operand.(map[string]interface{})
			id, _ := operandSet["id"].(string)
			IdpID, _ := operandSet["idp_id"].(string)
			if operandSet != nil {
				operandsSets = append(operandsSets, policysetrule.Operands{
					ID:         id,
					IdpID:      IdpID,
					LHS:        operandSet["lhs"].(string),
					ObjectType: operandSet["object_type"].(string),
					RHS:        operandSet["rhs"].(string),
					Name:       operandSet["name"].(string),
				})
			}
		}

		return operandsSets
	}
	return []policysetrule.Operands{}
}
func flattenPolicyRuleConditions(conditions []policysetrule.Conditions) []interface{} {
	ruleConditions := make([]interface{}, len(conditions))
	for i, ruleConditionItems := range conditions {
		ruleConditions[i] = map[string]interface{}{
			"id":       ruleConditionItems.ID,
			"negated":  ruleConditionItems.Negated,
			"operator": ruleConditionItems.Operator,
			"operands": flattenPolicyRuleOperands(ruleConditionItems.Operands),
		}
	}

	return ruleConditions
}

func flattenPolicyRuleOperands(conditionOperand []policysetrule.Operands) []interface{} {
	conditionOperands := make([]interface{}, len(conditionOperand))
	for i, operandItems := range conditionOperand {
		conditionOperands[i] = map[string]interface{}{
			"id":          operandItems.ID,
			"idp_id":      operandItems.IdpID,
			"lhs":         operandItems.LHS,
			"object_type": operandItems.ObjectType,
			"rhs":         operandItems.RHS,
		}
	}

	return conditionOperands
}

func flattenPolicyRuleServerGroups(appServerGroup []policysetrule.AppServerGroups) []interface{} {
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

func flattenPolicyRuleAppConnectorGroups(appConnectorGroups []policysetrule.AppConnectorGroups) []interface{} {
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
