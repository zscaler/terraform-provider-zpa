package zpa

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/policysetcontroller"
)

func dataSourcePolicyType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePolicyTypeRead,
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
			"modified_by": {
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
			"sorted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ACCESS_POLICY", "GLOBAL_POLICY",
					"TIMEOUT_POLICY", "REAUTH_POLICY",
					"CLIENT_FORWARDING_POLICY", "BYPASS_POLICY",
					"ISOLATION_POLICY", "INSPECTION_POLICY",
					"SIEM_POLICY",
				}, false),
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
						"modified_by": {
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
									"modified_by": {
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
												"modified_by": {
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

func dataSourcePolicyTypeRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	log.Printf("[INFO] Getting data for policy type\n")
	var resp *policysetcontroller.PolicySet
	var err error
	policyType, policyTypeIsSet := d.GetOk("policy_type")
	if policyTypeIsSet {
		resp, _, err = zClient.policysetcontroller.GetByPolicyType(policyType.(string))
	} else {
		resp, _, err = zClient.policysetcontroller.GetByPolicyType("GLOBAL_POLICY")
	}
	if err != nil {
		return err
	}

	log.Printf("[INFO] Getting data for Policy Type:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("creation_time", resp.CreationTime)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("modified_by", resp.ModifiedBy)
	_ = d.Set("modified_time", resp.ModifiedTime)
	_ = d.Set("name", resp.Name)
	_ = d.Set("sorted", resp.Sorted)
	_ = d.Set("policy_type", resp.PolicyType)

	if err := d.Set("rules", flattenPolicySetRules(resp)); err != nil {
		return err
	}

	return nil
}

func flattenPolicySetRules(policySetRules *policysetcontroller.PolicySet) []interface{} {
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
			"modified_by":                 ruleItem.ModifiedBy,
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

func flattenRuleConditions(conditions policysetcontroller.PolicyRule) []interface{} {
	ruleConditions := make([]interface{}, len(conditions.Conditions))
	for i, ruleCondition := range conditions.Conditions {
		ruleConditions[i] = map[string]interface{}{
			"creation_time": ruleCondition.CreationTime,
			"id":            ruleCondition.ID,
			"modified_by":   ruleCondition.ModifiedBy,
			"modified_time": ruleCondition.ModifiedTime,
			"negated":       ruleCondition.Negated,
			"operator":      ruleCondition.Operator,
			"operands":      flattenConditionOperands(ruleCondition),
		}
	}

	return ruleConditions
}

func flattenConditionOperands(operands policysetcontroller.Conditions) []interface{} {
	conditionOperands := make([]interface{}, len(operands.Operands))
	for i, conditionOperand := range operands.Operands {
		conditionOperands[i] = map[string]interface{}{
			"creation_time": conditionOperand.CreationTime,
			"id":            conditionOperand.ID,
			"idp_id":        conditionOperand.IdpID,
			"lhs":           conditionOperand.LHS,
			"modified_by":   conditionOperand.ModifiedBy,
			"modified_time": conditionOperand.ModifiedTime,
			"name":          conditionOperand.Name,
			"object_type":   conditionOperand.ObjectType,
			"rhs":           conditionOperand.RHS,
		}
	}

	return conditionOperands
}
