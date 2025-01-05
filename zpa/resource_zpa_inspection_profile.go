package zpa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceInspectionProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInspectionProfileCreate,
		ReadContext:   resourceInspectionProfileRead,
		UpdateContext: resourceInspectionProfileUpdate,
		DeleteContext: resourceInspectionProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("profile_id", id)
				} else {
					resp, _, err := inspection_profile.GetByName(ctx, service, id)
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the AppProtection profile",
			},
			"api_profile": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
			},
			"override_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				ValidateFunc: validation.StringInSlice([]string{
					"COMMON",
					"NONE",
					"SPECIFIC",
				}, false),
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
			"global_control_actions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The actions of the predefined, custom, or override controls",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"common_global_override_actions_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"paranoia_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OWASP Predefined Paranoia Level",
			},
			"predefined_controls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The predefined controls",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"predefined_api_controls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The predefined controls",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"threat_labz_controls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The ThreatLabZ predefined controls",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"websocket_controls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The WebSocket predefined controls",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"predefined_controls_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "OWASP_CRS/3.3.0",
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

func resourceInspectionProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandInspectionProfile(d)
	log.Printf("[INFO] Creating inspection profile with request\n%+v\n", req)
	if err := validateInspectionProfile(&req); err != nil {
		return diag.FromErr(err)
	}
	// injectPredefinedControls(zClient, &req)
	resp, _, err := inspection_profile.Create(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created inspection profile  request. ID: %v\n", resp)

	d.SetId(resp.ID)
	if v, ok := d.GetOk("associate_all_controls"); ok && v.(bool) {
		p, _, err := inspection_profile.Get(ctx, service, resp.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		inspection_profile.PutAssociate(ctx, service, resp.ID, p)
	}
	return resourceInspectionProfileRead(ctx, d, meta)
}

func resourceInspectionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, _, err := inspection_profile.Get(ctx, service, d.Id())
	if err != nil {
		if errResp, ok := err.(*errorx.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing inspection profile %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Getting inspection profile:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("api_profile", resp.APIProfile)
	_ = d.Set("override_action", resp.OverrideAction)
	_ = d.Set("associate_all_controls", d.Get("associate_all_controls"))
	_ = d.Set("common_global_override_actions_config", resp.CommonGlobalOverrideActionsConfig)
	_ = d.Set("global_control_actions", resp.GlobalControlActions)
	_ = d.Set("paranoia_level", resp.ParanoiaLevel)

	// Ensure the predefined_controls_version is set correctly
	if resp.PredefinedControlsVersion != "" {
		_ = d.Set("predefined_controls_version", resp.PredefinedControlsVersion)
	} else {
		_ = d.Set("predefined_controls_version", "OWASP_CRS/3.3.0")
	}

	if len(resp.ControlInfoResource) > 0 {
		if err := d.Set("controls_info", flattenControlInfoResource(resp.ControlInfoResource)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("custom_controls", flattenCustomControlsSimple(resp.CustomControls)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("predefined_controls", flattenPredefinedControlsSimple(resp.PredefinedControls)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting predefined_controls: %s", err))
	}

	if err := d.Set("predefined_api_controls", flattenPredefinedApiControlsSimple(resp.PredefinedAPIControls)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting predefined_api_controls: %s", err))
	}

	if resp.WebSocketControls != nil {
		if err := d.Set("websocket_controls", flattenWebSocketControls(resp.WebSocketControls)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting websocket_controls: %s", err))
		}
	}

	if resp.ThreatLabzControls != nil {
		if err := d.Set("threat_labz_controls", flattenThreatLabzControls(resp.ThreatLabzControls)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting threat_labz_controls: %s", err))
		}
	}

	return nil
}

func resourceInspectionProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id := d.Id()
	log.Printf("[INFO] Updating inspection profile ID: %v\n", id)
	req := expandInspectionProfile(d)
	if err := validateInspectionProfile(&req); err != nil {
		return diag.FromErr(err)
	}

	if _, _, err := inspection_profile.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	// injectPredefinedControls(zClient, &req)
	if _, err := inspection_profile.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}
	if v, ok := d.GetOk("associate_all_controls"); ok && v.(bool) {
		p, _, err := inspection_profile.Get(ctx, service, req.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		inspection_profile.PutAssociate(ctx, service, req.ID, p)
	}
	return resourceInspectionProfileRead(ctx, d, meta)
}

func resourceInspectionProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Deleting inspection profile ID: %v\n", d.Id())

	if _, err := inspection_profile.Delete(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] inspection profile deleted")
	return nil
}

func expandInspectionProfile(d *schema.ResourceData) inspection_profile.InspectionProfile {
	inspection_profile := inspection_profile.InspectionProfile{
		ID:                                d.Id(),
		Name:                              d.Get("name").(string),
		Description:                       d.Get("description").(string),
		APIProfile:                        d.Get("api_profile").(bool),
		OverrideAction:                    d.Get("override_action").(string),
		CommonGlobalOverrideActionsConfig: d.Get("common_global_override_actions_config").(map[string]interface{}),
		GlobalControlActions:              SetToStringList(d, "global_control_actions"),
		ParanoiaLevel:                     d.Get("paranoia_level").(string),
		PredefinedControlsVersion:         d.Get("predefined_controls_version").(string),
		ControlInfoResource:               expandControlsInfo(d),
		CustomControls:                    expandCustomControls(d),
		PredefinedAPIControls:             expandPredefinedAPIControls(d),
		PredefinedControls:                expandPredefinedControls(d),
		ThreatLabzControls:                expandThreatLabzControls(d),
		WebSocketControls:                 expandWebSocketControls(d),
	}

	if inspection_profile.PredefinedControlsVersion == "" {
		inspection_profile.PredefinedControlsVersion = "OWASP_CRS/3.3.0"
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

func expandThreatLabzControls(d *schema.ResourceData) []inspection_profile.ThreatLabzControls {
	controlsData := d.Get("threat_labz_controls").([]interface{})
	var controls []inspection_profile.ThreatLabzControls
	for _, item := range controlsData {
		controlMap := item.(map[string]interface{})
		controls = append(controls, inspection_profile.ThreatLabzControls{
			ID:          controlMap["id"].(string),
			Action:      controlMap["action"].(string),
			ActionValue: controlMap["action_value"].(string),
		})
	}
	return controls
}

func expandWebSocketControls(d *schema.ResourceData) []inspection_profile.WebSocketControls {
	controlsData := d.Get("websocket_controls").([]interface{})
	var controls []inspection_profile.WebSocketControls
	for _, item := range controlsData {
		controlMap := item.(map[string]interface{})
		controls = append(controls, inspection_profile.WebSocketControls{
			ID:          controlMap["id"].(string),
			Action:      controlMap["action"].(string),
			ActionValue: controlMap["action_value"].(string),
		})
	}
	return controls
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

func flattenThreatLabzControls(controls []inspection_profile.ThreatLabzControls) []interface{} {
	var controlsData []interface{}
	for _, control := range controls {
		controlMap := make(map[string]interface{})
		controlMap["id"] = control.ID
		controlMap["action"] = control.Action
		controlMap["action_value"] = control.ActionValue
		controlsData = append(controlsData, controlMap)
	}
	return controlsData
}

func flattenWebSocketControls(webSocketControls []inspection_profile.WebSocketControls) []interface{} {
	var controlsData []interface{}
	for _, webSocketControl := range webSocketControls {
		controlMap := make(map[string]interface{})
		controlMap["id"] = webSocketControl.ID
		controlMap["action"] = webSocketControl.Action
		controlMap["action_value"] = webSocketControl.ActionValue
		controlsData = append(controlsData, controlMap)
	}
	return controlsData
}
