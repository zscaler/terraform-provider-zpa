package zpa

import (
	"log"

	"github.com/SecurityGeekIO/terraform-provider-zpa/gozscaler/policysetglobal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGlobalAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlobalAccessPolicyRead,
		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modifiedby": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bypass_default_rule": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isolation_default_rule": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"modifiedby": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reauth_default_rule": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"reauth_idle_timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reauth_timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_order": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zpn_cbi_profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zpn_inspection_profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zpn_inspection_profile_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"modifiedby": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"modified_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"negated": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operands": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"creation_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"idp_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"lhs": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"modifiedby": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"modified_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"object_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"rhs": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
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
	}
}

func dataSourceGlobalAccessPolicyRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	log.Printf("[INFO] Getting data for global policy set\n")

	resp, _, err := zClient.policysetglobal.Get()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting Policy Set Global Rules:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("creation_time", resp.CreationTime)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("modifiedby", resp.ModifiedBy)
	_ = d.Set("modified_time", resp.ModifiedTime)
	_ = d.Set("name", resp.Name)
	_ = d.Set("policy_type", resp.PolicyType)

	if err := d.Set("rules", flattenPolicySetRules(resp)); err != nil {
		return err
	}

	return nil
}

func flattenPolicySetRules(policySetRules *policysetglobal.PolicySet) []interface{} {
	ruleItems := make([]interface{}, len(policySetRules.Rules))
	for i, ruleItem := range policySetRules.Rules {
		ruleItems[i] = map[string]interface{}{
			"action":                      ruleItem.Action,
			"action_id":                   ruleItem.ActionID,
			"creation_time":               ruleItem.CreationTime,
			"custom_msg":                  ruleItem.CustomMsg,
			"description":                 ruleItem.Description,
			"id":                          ruleItem.ID,
			"isolation_default_rule":      ruleItem.IsolationDefaultRule,
			"modifiedby":                  ruleItem.ModifiedBy,
			"modified_time":               ruleItem.ModifiedTime,
			"operator":                    ruleItem.Operator,
			"policy_set_id":               ruleItem.PolicySetID,
			"policy_type":                 ruleItem.PolicyType,
			"priority":                    ruleItem.Priority,
			"reauth_default_rule":         ruleItem.ReauthDefaultRule,
			"reauth_idle_timeout":         ruleItem.ReauthIdleTimeout,
			"reauth_timeout":              ruleItem.ReauthTimeout,
			"rule_order":                  ruleItem.RuleOrder,
			"zpn_cbi_profile_id":          ruleItem.ZpnCbiProfileID,
			"zpn_inspection_profile_id":   ruleItem.ZpnInspectionProfileID,
			"zpn_inspection_profile_name": ruleItem.ZpnInspectionProfileName,
			"conditions":                  flattenRuleConditions(ruleItem),
		}
	}

	return ruleItems
}

func flattenRuleConditions(conditions policysetglobal.Rules) []interface{} {
	ruleConditions := make([]interface{}, len(conditions.Conditions))
	for i, ruleCondition := range conditions.Conditions {
		ruleConditions[i] = map[string]interface{}{
			"creation_time": ruleCondition.CreationTime,
			"id":            ruleCondition.ID,
			"modifiedby":    ruleCondition.ModifiedBy,
			"modified_time": ruleCondition.ModifiedTime,
			"negated":       ruleCondition.Negated,
			"operator":      ruleCondition.Operator,
			"operands":      flattenConditionOperands(ruleCondition),
		}
	}

	return ruleConditions
}

func flattenConditionOperands(operands policysetglobal.Conditions) []interface{} {
	conditionOperands := make([]interface{}, len(*operands.Operands))
	for i, conditionOperand := range *operands.Operands {
		conditionOperands[i] = map[string]interface{}{
			"creation_time": conditionOperand.CreationTime,
			"id":            conditionOperand.ID,
			"idp_id":        conditionOperand.IdpID,
			"lhs":           conditionOperand.LHS,
			"modifiedby":    conditionOperand.ModifiedBy,
			"modified_time": conditionOperand.ModifiedTime,
			"name":          conditionOperand.Name,
			"object_type":   conditionOperand.ObjectType,
			"rhs":           conditionOperand.RHS,
		}
	}

	return conditionOperands
}
