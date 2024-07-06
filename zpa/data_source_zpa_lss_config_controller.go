package zpa

import (
	"fmt"
	"html"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/lssconfigcontroller"
)

func dataSourceLSSConfigController() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceLSSConfigControllerRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_message": {
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
						"filter": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"lss_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lss_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_log_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_tls": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"connector_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_rule": {
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
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_rule": {
							Type:     schema.TypeBool,
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
						"lss_default_rule": {
							Type:     schema.TypeBool,
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
					},
				},
			},
		},
	}
}

func dataSourceLSSConfigControllerRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.LSSConfigController

	var resp *lssconfigcontroller.LSSResource
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for lss config controller %s\n", id)
		res, _, err := lssconfigcontroller.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for lss config controller %s\n", name)
		res, _, err := lssconfigcontroller.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		if err := d.Set("config", flattenLSSConfig(resp.LSSConfig)); err != nil {
			return err
		}

		_ = d.Set("connector_groups", flattenConnectorGroups(resp.ConnectorGroups))

		if err := d.Set("policy_rule", flattenLSSPolicyRule(resp.PolicyRule)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any lss config controller with name '%s' or id", id)
	}

	return nil
}

func flattenLSSConfig(lssConfig *lssconfigcontroller.LSSConfig) interface{} {
	return []map[string]interface{}{
		{
			"audit_message":   html.UnescapeString(lssConfig.AuditMessage),
			"description":     lssConfig.Description,
			"enabled":         lssConfig.Enabled,
			"filter":          lssConfig.Filter,
			"id":              lssConfig.ID,
			"name":            lssConfig.Name,
			"lss_host":        lssConfig.LSSHost,
			"lss_port":        lssConfig.LSSPort,
			"source_log_type": lssConfig.SourceLogType,
			"use_tls":         lssConfig.UseTLS,
			"format":          html.UnescapeString(lssConfig.Format),
		},
	}
}

func flattenConnectorGroups(lssConnectorGroup []lssconfigcontroller.ConnectorGroups) []interface{} {
	lssConnectorGroups := make([]interface{}, len(lssConnectorGroup))
	for i, val := range lssConnectorGroup {
		lssConnectorGroups[i] = map[string]interface{}{
			"id": val.ID,
		}
	}

	return lssConnectorGroups
}

func flattenLSSPolicyRule(policySetRules *lssconfigcontroller.PolicyRule) interface{} {
	if policySetRules == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"action":                      policySetRules.Action,
			"action_id":                   policySetRules.ActionID,
			"creation_time":               policySetRules.CreationTime,
			"custom_msg":                  policySetRules.CustomMsg,
			"description":                 policySetRules.Description,
			"id":                          policySetRules.ID,
			"isolation_default_rule":      policySetRules.IsolationDefaultRule,
			"modifiedby":                  policySetRules.ModifiedBy,
			"modified_time":               policySetRules.ModifiedTime,
			"operator":                    policySetRules.Operator,
			"policy_set_id":               policySetRules.PolicySetID,
			"policy_type":                 policySetRules.PolicyType,
			"priority":                    policySetRules.Priority,
			"reauth_default_rule":         policySetRules.ReauthDefaultRule,
			"reauth_idle_timeout":         policySetRules.ReauthIdleTimeout,
			"reauth_timeout":              policySetRules.ReauthTimeout,
			"rule_order":                  policySetRules.RuleOrder,
			"zpn_cbi_profile_id":          policySetRules.ZpnCbiProfileID,
			"zpn_inspection_profile_id":   policySetRules.ZpnInspectionProfileID,
			"zpn_inspection_profile_name": policySetRules.ZpnInspectionProfileName,
			"conditions":                  flattenLSSRuleConditions(policySetRules),
		},
	}
}

func flattenLSSRuleConditions(conditions *lssconfigcontroller.PolicyRule) []interface{} {
	ruleConditions := make([]interface{}, len(conditions.Conditions))
	for i, ruleCondition := range conditions.Conditions {
		ruleConditions[i] = map[string]interface{}{
			"creation_time": ruleCondition.CreationTime,
			"id":            ruleCondition.ID,
			"modifiedby":    ruleCondition.ModifiedBy,
			"modified_time": ruleCondition.ModifiedTime,
			"operator":      ruleCondition.Operator,
			"operands":      flattenLSSConditionOperands(ruleCondition),
		}
	}

	return ruleConditions
}

func flattenLSSConditionOperands(operands lssconfigcontroller.Conditions) []interface{} {
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
