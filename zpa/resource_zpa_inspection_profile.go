package zpa

import (
	"errors"
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

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("profile_id", id)
				} else {
					resp, _, err := zClient.inspection_profile.GetByName(id)
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
			"common_global_override_actions_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
								"ZSCALER",
								"CUSTOM",
								"PREDEFINED",
							}, false),
						},
						"count": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"custom_controls": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PASS",
								"BLOCK",
								"REDIRECT",
							}, false),
						},
						"action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"global_control_actions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"incarnation_number": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"paranoia_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"predefined_controls": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PASS",
								"BLOCK",
								"REDIRECT",
							}, false),
						},
						"control_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"WEBSOCKET_PREDEFINED",
								"WEBSOCKET_CUSTOM",
								"ZSCALER",
								"CUSTOM",
								"PREDEFINED",
							}, false),
						},
						"action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"web_socket_controls": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PASS",
								"BLOCK",
								"REDIRECT",
							}, false),
						},
						"control_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"WEBSOCKET_PREDEFINED",
								"WEBSOCKET_CUSTOM",
								"ZSCALER",
								"CUSTOM",
								"PREDEFINED",
							}, false),
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP",
								"HTTPS",
								"FTP",
								"RDP",
								"SSH",
								"WEBSOCKET",
							}, false),
						},
						"action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
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

	req := expandInspectionProfile(d)
	log.Printf("[INFO] Creating inspection profile with request\n%+v\n", req)
	if err := validateInspectionProfile(&req); err != nil {
		return err
	}
	//injectPredefinedControls(zClient, &req)
	resp, _, err := zClient.inspection_profile.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created inspection profile  request. ID: %v\n", resp)

	d.SetId(resp.ID)
	if v, ok := d.GetOk("associate_all_controls"); ok && v.(bool) {
		p, _, err := zClient.inspection_profile.Get(resp.ID)
		if err != nil {
			return err
		}
		zClient.inspection_profile.PutAssociate(resp.ID, p)
	}
	return resourceInspectionProfileRead(d, m)
}

func resourceInspectionProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.inspection_profile.Get(d.Id())
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
	_ = d.Set("common_global_override_actions_config", resp.CommonGlobalOverrideActionsConfig)
	_ = d.Set("associate_all_controls", d.Get("associate_all_controls"))
	_ = d.Set("description", resp.Description)
	_ = d.Set("global_control_actions", resp.GlobalControlActions)
	_ = d.Set("name", resp.Name)
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
		return err
	}

	if err := d.Set("web_socket_controls", flattenPredefinedControlsSimple(resp.WebSocketControls)); err != nil {
		return err
	}
	return nil

}

func flattenPredefinedControlsSimple(predControl []inspection_profile.CustomCommonControls) []interface{} {
	predControls := make([]interface{}, len(predControl))
	for i, predControl := range predControl {
		predControls[i] = map[string]interface{}{
			"id":            predControl.ID,
			"action":        predControl.Action,
			"action_value":  predControl.ActionValue,
			"control_type":  predControl.ControlType,
			"protocol_type": predControl.ProtocolType,
		}
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

	id := d.Id()
	log.Printf("[INFO] Updating inspection profile ID: %v\n", id)
	req := expandInspectionProfile(d)
	if err := validateInspectionProfile(&req); err != nil {
		return err
	}

	if _, _, err := zClient.inspection_profile.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	//injectPredefinedControls(zClient, &req)
	if _, err := zClient.inspection_profile.Update(id, &req); err != nil {
		return err
	}
	if v, ok := d.GetOk("associate_all_controls"); ok && v.(bool) {
		p, _, err := zClient.inspection_profile.Get(req.ID)
		if err != nil {
			return err
		}
		zClient.inspection_profile.PutAssociate(req.ID, p)
	}
	return resourceInspectionProfileRead(d, m)
}

func resourceInspectionProfileDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting inspection profile ID: %v\n", d.Id())

	if _, err := zClient.inspection_profile.Delete(d.Id()); err != nil {
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
			Count:       controlItem["count"].(string),
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
	predControlsInterface, ok := d.GetOk("predefined_controls")
	if ok {
		predControl := predControlsInterface.(*schema.Set)
		log.Printf("[INFO] predefined control data: %+v\n", predControl)
		var predControls []inspection_profile.CustomCommonControls
		for _, predControl := range predControl.List() {
			predControl, ok := predControl.(map[string]interface{})
			if ok {
				actionValue := ""
				if predControl["action_value"] != nil {
					actionValue = predControl["action_value"].(string)
				}
				predControls = append(predControls, inspection_profile.CustomCommonControls{
					ID:           predControl["id"].(string),
					Action:       predControl["action"].(string),
					ControlType:  predControl["control_type"].(string),
					ProtocolType: predControl["protocol_type"].(string),
					ActionValue:  actionValue,
				})
			}
		}
		return predControls
	}

	return []inspection_profile.CustomCommonControls{}
}
