package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_profile"
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

func dataSourceInspectionProfileRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.InspectionProfile

	var resp *inspection_profile.InspectionProfile
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for inspection profile  %s\n", id)
		res, _, err := inspection_profile.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for inspection profile name %s\n", name)
		res, _, err := inspection_profile.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("description", resp.Description)
		_ = d.Set("global_control_actions", resp.GlobalControlActions)
		_ = d.Set("incarnation_number", resp.IncarnationNumber)
		_ = d.Set("modified_by", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("paranoia_level", resp.ParanoiaLevel)
		_ = d.Set("predefined_controls_version", resp.PredefinedControlsVersion)

		if err := d.Set("controls_info", flattenControlInfoResource(resp.ControlInfoResource)); err != nil {
			return err
		}

		if err := d.Set("custom_controls", flattenCustomControls(resp.CustomControls)); err != nil {
			return err
		}

		if err := d.Set("predefined_controls", flattenPredefinedControls(resp.PredefinedControls)); err != nil {
			return err
		}

		// Flattening ThreatLabz Controls
		// threatLabzIDs := make([]string, len(resp.ThreatLabzControls))
		// for i, control := range resp.ThreatLabzControls {
		// 	threatLabzIDs[i] = control.ID
		// }
		// if err := d.Set("threatlabz_controls", flattenIDList(threatLabzIDs)); err != nil {
		// 	return err
		// }

		// // Flattening WebSocket Controls
		// websocketIDs := make([]string, len(resp.WebSocketControls))
		// for i, control := range resp.WebSocketControls {
		// 	websocketIDs[i] = control.ID
		// }
		// if err := d.Set("websocket_controls", flattenIDList(websocketIDs)); err != nil {
		// 	return err
		// }
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
			"type":                                custom.Type,
			"version":                             custom.Version,
			"associated_inspection_profile_names": flattenAssociatedInspectionProfileNames(custom.AssociatedInspectionProfileNames),
			"rules":                               flattenInspectionRules(custom.Rules),
		}
	}

	return customControls
}
