package zpa

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/client"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/inspection_profile"
)

func resourceinspection_profile() *schema.Resource {
	return &schema.Resource{
		Create: resourceinspection_profileCreate,
		Read:   resourceinspection_profileRead,
		Update: resourceinspection_profileUpdate,
		Delete: resourceinspection_profileDelete,
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
								"CUSTOM",
								"PREDEFINED",
								"ZSCALER",
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
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
				Type:     schema.TypeList,
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
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"predefined_controls_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceinspection_profileCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandinspection_profile(d)
	log.Printf("[INFO] Creating segment group with request\n%+v\n", req)

	resp, _, err := zClient.inspection_profile.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created segment group request. ID: %v\n", resp)

	d.SetId(resp.ID)
	return resourceinspection_profileRead(d, m)
}

func resourceinspection_profileRead(d *schema.ResourceData, m interface{}) error {
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
	_ = d.Set("creation_time", resp.CreationTime)
	_ = d.Set("description", resp.Description)
	_ = d.Set("global_control_actions", resp.GlobalControlActions)
	_ = d.Set("incarnation_number", resp.IncarnationNumber)
	_ = d.Set("modifiedby", resp.ModifiedBy)
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
	return nil

}

func resourceinspection_profileUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating inspection profile ID: %v\n", id)
	req := expandinspection_profile(d)

	if _, err := zClient.inspection_profile.Update(id, &req); err != nil {
		return err
	}

	return resourceinspection_profileRead(d, m)
}

func resourceinspection_profileDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting inspection profile ID: %v\n", d.Id())

	if _, err := zClient.inspection_profile.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] inspection profile deleted")
	return nil
}

func expandinspection_profile(d *schema.ResourceData) inspection_profile.inspection_profile {
	inspection_profile := inspection_profile.inspection_profile{
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

func expandCustomControls(d *schema.ResourceData) []inspection_profile.CustomControls {
	customControlsInterface, ok := d.GetOk("custom_controls")
	if ok {
		control := customControlsInterface.(*schema.Set)
		log.Printf("[INFO] custom control data: %+v\n", control)
		var customControls []inspection_profile.CustomControls
		for _, customControl := range control.List() {
			customControl, ok := customControl.(map[string]interface{})
			if ok {
				for _, id := range customControl["id"].(*schema.Set).List() {
					customControls = append(customControls, inspection_profile.CustomControls{
						ID: id.(string),
					})
				}
			}
		}
		return customControls
	}

	return []inspection_profile.CustomControls{}
}

func expandPredefinedControls(d *schema.ResourceData) []inspection_profile.PredefinedControls {
	predControlsInterface, ok := d.GetOk("predefined_controls")
	if ok {
		predControl := predControlsInterface.(*schema.Set)
		log.Printf("[INFO] predefined control data: %+v\n", predControl)
		var predControls []inspection_profile.PredefinedControls
		for _, predControl := range predControl.List() {
			predControl, ok := predControl.(map[string]interface{})
			if ok {
				for _, id := range predControl["id"].(*schema.Set).List() {
					predControls = append(predControls, inspection_profile.PredefinedControls{
						ID: id.(string),
					})
				}
			}
		}
		return predControls
	}

	return []inspection_profile.PredefinedControls{}
}
