package zpa

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_profile"
)

func resourceInspectionProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceInspectionProfileCreate,
		Read:   resourceInspectionProfileRead,
		Update: resourceInspectionProfileUpdate,
		Delete: resourceInspectionProfileDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.InspectionProfile

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("profile_id", id)
				} else {
					resp, _, err := inspection_profile.GetByName(service, id)
					if err == nil {
						d.SetId(resp.ID)
						d.Set("profile_id", resp.ID)
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"associate_all_controls": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"controls_info": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"WEBSOCKET_PREDEFINED",
								"WEBSOCKET_CUSTOM",
								"THREATLABZ",
								"CUSTOM",
								"PREDEFINED",
							}, false),
						},
						// "count": {
						// 	Type:     schema.TypeString,
						// 	Optional: true,
						// 	Computed: true,
						// },
					},
				},
			},
			"custom_controls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The set of AppProtection controls used to define how inspections are managed",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier of the custom control",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The action of the custom control",
							ValidateFunc: validation.StringInSlice([]string{
								"PASS",
								"BLOCK",
								"REDIRECT",
							}, false),
						},
						"action_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Denotes the action. Supports any string",
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the AppProtection profile",
			},
			"global_control_actions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The actions of the predefined, custom, or override controls",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"incarnation_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"paranoia_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OWASP Predefined Paranoia Level",
			},
			"predefined_controls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The predefined controls",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier of the predefined control",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The action of the predefined control",
							ValidateFunc: validation.StringInSlice([]string{
								"PASS",
								"BLOCK",
								"REDIRECT",
							}, false),
						},
						"control_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The control type of the custom control",
							ValidateFunc: validation.StringInSlice([]string{
								"WEBSOCKET_PREDEFINED",
								"WEBSOCKET_CUSTOM",
								"THREATLABZ",
								"CUSTOM",
								"PREDEFINED",
							}, false),
						},
						"action_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value for the predefined controls action. This field is only required if the action is set to REDIRECT",
						},
						"protocol_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The protocol type of the predefined control",
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP",
								"HTTPS",
								"FTP",
								"RDP",
								"SSH",
								"WEBSOCKET",
							}, false),
						},
					},
				},
			},
			"predefined_controls_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol for the AppProtection application",
			},
			"zs_defined_control_choice": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates the user's choice for the ThreatLabZ Controls. Supported values: ALL and SPECIFIC",
				ValidateFunc: validation.StringInSlice([]string{
					"ALL",
					"SPECIFIC",
				}, false),
			},
			// "threatlabz_controls": {
			// 	Type:        schema.TypeSet,
			// 	Optional:    true,
			// 	Description: "The ThreatLabZ predefined controls",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeString,
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			// "websocket_controls": {
			// 	Type:        schema.TypeSet,
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "The WebSocket controls.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeString,
			// 				},
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

func validateInspectionProfile(profile *inspection_profile.InspectionProfile) error {
	for _, d := range profile.CustomControls {
		if d.Action == "REDIRECT" && d.ActionValue == "" {
			return errors.New("when action is REDIRECT, action value must be set")
		}
	}
	for _, d := range profile.PredefinedControls {
		if d.Action == "REDIRECT" && d.ActionValue == "" {
			return errors.New("when action is REDIRECT, action value must be set")
		}
	}
	return nil
}

func resourceInspectionProfileCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.InspectionProfile

	req := expandInspectionProfile(d)
	log.Printf("[INFO] Creating inspection profile with request\n%+v\n", req)
	if err := validateInspectionProfile(&req); err != nil {
		return err
	}
	// injectPredefinedControls(zClient, &req)
	resp, _, err := inspection_profile.Create(service, req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created inspection profile  request. ID: %v\n", resp)

	d.SetId(resp.ID)
	if v, ok := d.GetOk("associate_all_controls"); ok && v.(bool) {
		p, _, err := inspection_profile.Get(service, resp.ID)
		if err != nil {
			return err
		}
		inspection_profile.PutAssociate(service, resp.ID, p)
	}
	return resourceInspectionProfileRead(d, m)
}

func resourceInspectionProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.InspectionProfile

	resp, _, err := inspection_profile.Get(service, d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing inspection profile %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	log.Printf("[INFO] Getting inspection profile:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("associate_all_controls", d.Get("associate_all_controls"))
	_ = d.Set("global_control_actions", resp.GlobalControlActions)
	_ = d.Set("paranoia_level", resp.ParanoiaLevel)
	if resp.PredefinedControlsVersion != "" {
		_ = d.Set("predefined_controls_version", resp.PredefinedControlsVersion)
	}
	if len(resp.ControlInfoResource) > 0 {
		if err := d.Set("controls_info", flattenControlInfoResource(resp.ControlInfoResource)); err != nil {
			return err
		}
	}

	if err := d.Set("custom_controls", flattenCustomControlsSimple(resp.CustomControls)); err != nil {
		return err
	}

	if err := d.Set("predefined_controls", flattenPredefinedControlsSimple(resp.PredefinedControls)); err != nil {
		return fmt.Errorf("error setting predefined_controls: %s", err)
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
	return nil
}

func flattenPredefinedControlsSimple(predControl []inspection_profile.CustomCommonControls) []interface{} {
	if len(predControl) == 0 {
		return nil
	}

	predControls := make([]interface{}, len(predControl))
	for i, control := range predControl {
		controlMap := make(map[string]interface{})
		controlMap["id"] = control.ID
		controlMap["action"] = control.Action
		controlMap["control_type"] = control.ControlType
		controlMap["protocol_type"] = control.ProtocolType
		// Include other fields if necessary
		predControls[i] = controlMap
	}

	return predControls
}

func flattenCustomControlsSimple(customControl []inspection_profile.InspectionCustomControl) []interface{} {
	customControls := make([]interface{}, len(customControl))
	for i, custom := range customControl {
		customControls[i] = map[string]interface{}{
			"id":           custom.ID,
			"action":       custom.Action,
			"action_value": custom.ActionValue,
		}
	}

	return customControls
}

func resourceInspectionProfileUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.InspectionProfile

	id := d.Id()
	log.Printf("[INFO] Updating inspection profile ID: %v\n", id)
	req := expandInspectionProfile(d)
	if err := validateInspectionProfile(&req); err != nil {
		return err
	}

	if _, _, err := inspection_profile.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	// injectPredefinedControls(zClient, &req)
	if _, err := inspection_profile.Update(service, id, &req); err != nil {
		return err
	}
	if v, ok := d.GetOk("associate_all_controls"); ok && v.(bool) {
		p, _, err := inspection_profile.Get(service, req.ID)
		if err != nil {
			return err
		}
		inspection_profile.PutAssociate(service, req.ID, p)
	}
	return resourceInspectionProfileRead(d, m)
}

func resourceInspectionProfileDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.InspectionProfile

	log.Printf("[INFO] Deleting inspection profile ID: %v\n", d.Id())

	if _, err := inspection_profile.Delete(service, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] inspection profile deleted")
	return nil
}

func expandInspectionProfile(d *schema.ResourceData) inspection_profile.InspectionProfile {
	inspection_profile := inspection_profile.InspectionProfile{
		ID:                        d.Id(),
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		GlobalControlActions:      SetToStringList(d, "global_control_actions"),
		IncarnationNumber:         d.Get("incarnation_number").(string),
		ParanoiaLevel:             d.Get("paranoia_level").(string),
		PredefinedControlsVersion: d.Get("predefined_controls_version").(string),
		ControlInfoResource:       expandControlsInfo(d),
		CustomControls:            expandCustomControls(d),
		PredefinedControls:        expandPredefinedControls(d),
		// WebSocketControls:         expandWebSocketControl(d),
		// ThreatLabzControls: expandThreatLabzControls(d),
	}
	return inspection_profile
}

func expandControlsInfo(d *schema.ResourceData) []inspection_profile.ControlInfoResource {
	var controlItems []inspection_profile.ControlInfoResource
	controlInterface, ok := d.GetOk("controls_info")
	if !ok {
		return controlItems
	}
	controlInfo, ok := controlInterface.(*schema.Set)
	if !ok {
		return controlItems
	}
	for _, controlItemObj := range controlInfo.List() {
		controlItem, ok := controlItemObj.(map[string]interface{})
		if !ok {
			return controlItems
		}
		controlItems = append(controlItems, inspection_profile.ControlInfoResource{
			ControlType: controlItem["control_type"].(string),
			// Count:       controlItem["count"].(string),
		})
	}
	return controlItems
}

func expandCustomControls(d *schema.ResourceData) []inspection_profile.InspectionCustomControl {
	customControlsInterface, ok := d.GetOk("custom_controls")
	if ok {
		control := customControlsInterface.(*schema.Set)
		log.Printf("[INFO] custom control data: %+v\n", control)
		var customControls []inspection_profile.InspectionCustomControl
		for _, customControl := range control.List() {
			customControl, ok := customControl.(map[string]interface{})
			if ok {
				actionValue := ""
				if customControl["action_value"] != nil {
					actionValue = customControl["action_value"].(string)
				}
				customControls = append(customControls, inspection_profile.InspectionCustomControl{
					ID:          customControl["id"].(string),
					Action:      customControl["action"].(string),
					ActionValue: actionValue,
				})

			}
		}
		return customControls
	}

	return []inspection_profile.InspectionCustomControl{}
}

func expandPredefinedControls(d *schema.ResourceData) []inspection_profile.CustomCommonControls {
	if v, ok := d.GetOk("predefined_controls"); ok {
		predefinedControlsSet := v.(*schema.Set)
		var predefinedControls []inspection_profile.CustomCommonControls

		for _, v := range predefinedControlsSet.List() {
			controlMap := v.(map[string]interface{})

			control := inspection_profile.CustomCommonControls{
				ID:           controlMap["id"].(string),
				Action:       controlMap["action"].(string),
				ControlType:  controlMap["control_type"].(string),
				ProtocolType: controlMap["protocol_type"].(string),
			}

			// Only add action_value if it's set in the schema
			if actionValue, ok := controlMap["action_value"].(string); ok && actionValue != "" {
				control.ActionValue = actionValue
			}

			predefinedControls = append(predefinedControls, control)
		}

		return predefinedControls
	}

	return nil
}

/*
func expandThreatLabzControls(d *schema.ResourceData) []inspection_profile.ThreatLabzControls {
	threatLabzInterface, ok := d.GetOk("threatlabz_controls")
	if ok {
		threatLabzControl := threatLabzInterface.(*schema.Set)
		log.Printf("[INFO] threatlabz control data: %+v\n", threatLabzControl)
		var threatLabzControls []inspection_profile.ThreatLabzControls
		for _, threatLabzControl := range threatLabzControl.List() {
			threatLabzControl, _ := threatLabzControl.(map[string]interface{})
			if threatLabzControl != nil {
				for _, id := range threatLabzControl["id"].([]interface{}) {
					threatLabzControls = append(threatLabzControls, inspection_profile.ThreatLabzControls{
						ID: id.(string),
					})
				}
			}
		}
		return threatLabzControls
	}

	return []inspection_profile.ThreatLabzControls{}
}


	func expandWebSocketControl(d *schema.ResourceData) []inspection_profile.WebSocketControls {
		websocketInterface, ok := d.GetOk("websocket_controls")
		if ok {
			websocketControl := websocketInterface.(*schema.Set)
			log.Printf("[INFO] websocket control data: %+v\n", websocketControl)
			var websocketControls []inspection_profile.WebSocketControls
			for _, websocketControl := range websocketControl.List() {
				websocketControl, _ := websocketControl.(map[string]interface{})
				if websocketControl != nil {
					for _, id := range websocketControl["id"].([]interface{}) {
						websocketControls = append(websocketControls, inspection_profile.WebSocketControls{
							ID: id.(string),
						})
					}
				}
			}
			return websocketControls
		}

		return []inspection_profile.WebSocketControls{}
	}


func flattenIDList(idList []string) *schema.Set {
	set := schema.NewSet(schema.HashString, []interface{}{})

	for _, id := range idList {
		if id != "" {
			set.Add(id)
		}
	}

	return set
}
*/
