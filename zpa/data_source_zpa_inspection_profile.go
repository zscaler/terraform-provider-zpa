package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/inspectioncontrol/inspection_profile"
)

func dataSourceInspectionProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInspectionProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"check_control_deployment_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"zs_defined_control_choice": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"common_global_override_actions_config": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"controls_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"threat_labz_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"default_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_action_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"paranoia_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_deployment_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_metadata": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_processor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ruleset_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ruleset_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zscaler_info_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associated_customers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"customer_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exclude_constellation": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_partner": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"associated_inspection_profile_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"custom_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_rule_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_action_value": {
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
						"paranoia_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rules": dataInspectionRulesSchema(),
						"associated_inspection_profile_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_control_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"incarnation_number": {
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
			"paranoia_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"predefined_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associated_inspection_profile_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"attachment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_action_value": {
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
						"paranoia_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"predefined_controls_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_socket_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associated_inspection_profile_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"control_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_action_value": {
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
						"paranoia_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInspectionProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *inspection_profile.InspectionProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for inspection profile  %s\n", id)
		res, _, err := zClient.inspection_profile.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for inspection profile name %s\n", name)
		res, _, err := zClient.inspection_profile.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("common_global_override_actions_config", resp.CommonGlobalOverrideActionsConfig)
		_ = d.Set("check_control_deployment_status", resp.CheckControlDeploymentStatus)
		_ = d.Set("global_control_actions", resp.GlobalControlActions)
		_ = d.Set("incarnation_number", resp.IncarnationNumber)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("paranoia_level", resp.ParanoiaLevel)
		_ = d.Set("predefined_controls_version", resp.PredefinedControlsVersion)
		_ = d.Set("zs_defined_control_choice", resp.ZSDefinedControlChoice)

		if err := d.Set("controls_info", flattenControlInfoResource(resp.ControlInfoResource)); err != nil {
			return err
		}

		if err := d.Set("custom_controls", flattenCustomControls(resp.CustomControls)); err != nil {
			return err
		}

		if err := d.Set("predefined_controls", flattenPredefinedControls(resp.PredefinedControls)); err != nil {
			return err
		}

		if err := d.Set("web_socket_controls", flattenWebSocketControls(resp.WebSocketControls)); err != nil {
			return err
		}
		if err := d.Set("threat_labz_controls", flattenThreatLabzControls(resp.ThreatLabzControls)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any inspection profile with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenControlInfoResource(controlInfo []inspection_profile.ControlInfoResource) []interface{} {
	controlInfoResource := make([]interface{}, len(controlInfo))
	for i, val := range controlInfo {
		controlInfoResource[i] = map[string]interface{}{
			"control_type": val.ControlType,
			"count":        val.Count,
		}
	}

	return controlInfoResource
}

func flattenCustomControls(customControl []inspection_profile.InspectionCustomControl) []interface{} {
	customControls := make([]interface{}, len(customControl))
	for i, custom := range customControl {
		customControls[i] = map[string]interface{}{
			"id":                                  custom.ID,
			"name":                                custom.Name,
			"action":                              custom.Action,
			"action_value":                        custom.ActionValue,
			"control_number":                      custom.ControlNumber,
			"control_rule_json":                   custom.ControlRuleJson,
			"creation_time":                       custom.CreationTime,
			"default_action":                      custom.DefaultAction,
			"default_action_value":                custom.DefaultActionValue,
			"description":                         custom.Description,
			"modified_by":                         custom.ModifiedBy,
			"modified_time":                       custom.ModifiedTime,
			"paranoia_level":                      custom.ParanoiaLevel,
			"protocol_type":                       custom.ProtocolType,
			"type":                                custom.Type,
			"version":                             custom.Version,
			"associated_inspection_profile_names": flattenAssociatedInspectionProfileNames(custom.AssociatedInspectionProfileNames),
			"rules":                               flattenInspectionRules(custom.Rules),
		}
	}

	return customControls
}

func flattenAssociatedInspectionProfileNames(associatedInspectionProfileNames []inspection_profile.AssociatedProfileNames) []interface{} {
	rule := make([]interface{}, len(associatedInspectionProfileNames))
	for i, val := range associatedInspectionProfileNames {
		rule[i] = map[string]interface{}{
			"id":   val.ID,
			"name": val.Name,
		}
	}

	return rule
}

func flattenPredefinedControls(predControl []inspection_profile.CustomCommonControls) []interface{} {
	predControls := make([]interface{}, len(predControl))
	for i, predControl := range predControl {
		predControls[i] = map[string]interface{}{
			"id":                                  predControl.ID,
			"action":                              predControl.Action,
			"action_value":                        predControl.ActionValue,
			"attachment":                          predControl.Attachment,
			"control_group":                       predControl.ControlGroup,
			"control_number":                      predControl.ControlNumber,
			"control_type":                        predControl.ControlType,
			"creation_time":                       predControl.CreationTime,
			"default_action":                      predControl.DefaultAction,
			"default_action_value":                predControl.DefaultActionValue,
			"description":                         predControl.Description,
			"modified_by":                         predControl.ModifiedBy,
			"modified_time":                       predControl.ModifiedTime,
			"name":                                predControl.Name,
			"paranoia_level":                      predControl.ParanoiaLevel,
			"severity":                            predControl.Severity,
			"version":                             predControl.Version,
			"associated_inspection_profile_names": flattenAssociatedInspectionProfileNames(predControl.AssociatedInspectionProfileNames),
		}
	}

	return predControls
}

func flattenWebSocketControls(websocketControl []inspection_profile.CustomCommonControls) []interface{} {
	websocketControls := make([]interface{}, len(websocketControl))
	for i, websocketControl := range websocketControl {
		websocketControls[i] = map[string]interface{}{
			"id":                                  websocketControl.ID,
			"name":                                websocketControl.Name,
			"action":                              websocketControl.Action,
			"action_value":                        websocketControl.ActionValue,
			"control_number":                      websocketControl.ControlNumber,
			"control_type":                        websocketControl.ControlType,
			"creation_time":                       websocketControl.CreationTime,
			"default_action":                      websocketControl.DefaultAction,
			"default_action_value":                websocketControl.DefaultActionValue,
			"description":                         websocketControl.Description,
			"modified_by":                         websocketControl.ModifiedBy,
			"modified_time":                       websocketControl.ModifiedTime,
			"paranoia_level":                      websocketControl.ParanoiaLevel,
			"severity":                            websocketControl.Severity,
			"version":                             websocketControl.Version,
			"associated_inspection_profile_names": flattenAssociatedInspectionProfileNames(websocketControl.AssociatedInspectionProfileNames),
		}
	}

	return websocketControls
}

func flattenAssociatedCustomers(associatedCustomers []inspection_profile.AssociatedCustomers) []interface{} {
	rule := make([]interface{}, len(associatedCustomers))
	for i, val := range associatedCustomers {
		rule[i] = map[string]interface{}{
			"customer_id":           val.CustomerID,
			"exclude_constellation": val.ExcludeConstellation,
			"is_partner":            val.IsPartner,
			"name":                  val.Name,
		}
	}

	return rule
}

func flattenThreatLabzControls(threatLabzControl []inspection_profile.ThreatLabzControls) []interface{} {
	threatLabzControls := make([]interface{}, len(threatLabzControl))
	for i, threatLabzControl := range threatLabzControl {
		threatLabzControls[i] = map[string]interface{}{
			"id":                                  threatLabzControl.ID,
			"name":                                threatLabzControl.Name,
			"description":                         threatLabzControl.Description,
			"enabled":                             threatLabzControl.Enabled,
			"action":                              threatLabzControl.Action,
			"action_value":                        threatLabzControl.ActionValue,
			"control_group":                       threatLabzControl.ControlGroup,
			"control_number":                      threatLabzControl.ControlNumber,
			"control_type":                        threatLabzControl.ControlType,
			"attachment":                          threatLabzControl.Attachment,
			"creation_time":                       threatLabzControl.CreationTime,
			"default_action":                      threatLabzControl.DefaultAction,
			"default_action_value":                threatLabzControl.DefaultActionValue,
			"engine_version":                      threatLabzControl.EngineVersion,
			"last_deployment_time":                threatLabzControl.LastDeploymentTime,
			"rule_deployment_state":               threatLabzControl.RuleDeploymentState,
			"rule_metadata":                       threatLabzControl.RuleMetadata,
			"rule_processor":                      threatLabzControl.RuleProcessor,
			"ruleset_name":                        threatLabzControl.RulesetName,
			"ruleset_version":                     threatLabzControl.RulesetVersion,
			"zscaler_info_url":                    threatLabzControl.ZscalerInfoUrl,
			"modified_by":                         threatLabzControl.ModifiedBy,
			"modified_time":                       threatLabzControl.ModifiedTime,
			"paranoia_level":                      threatLabzControl.ParanoiaLevel,
			"severity":                            threatLabzControl.Severity,
			"version":                             threatLabzControl.Version,
			"associated_inspection_profile_names": flattenAssociatedInspectionProfileNames(threatLabzControl.AssociatedInspectionProfileNames),
			"associated_customers":                flattenAssociatedCustomers(threatLabzControl.AssociatedCustomers),
		}
	}

	return threatLabzControls
}
