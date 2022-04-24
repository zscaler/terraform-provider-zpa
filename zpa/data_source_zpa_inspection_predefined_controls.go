package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/inspectioncontrol/inspection_predefined_controls"
)

func dataSourceInspectionPredefinedControls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInspectionPredefinedControlsRead,
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
			"control_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_number": {
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
			"modifiedby": {
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
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:         schema.TypeString,
				RequiredWith: []string{"name"},
				Optional:     true,
			},
		},
	}
}

func dataSourceInspectionPredefinedControlsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *inspection_predefined_controls.PredefinedControls
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for predefined controls %s\n", id)
		res, _, err := zClient.inspection_predefined_controls.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		version, versionSet := d.Get("version").(string)
		if !versionSet || version == "" {
			return fmt.Errorf("when the name is set, version must be set as well")
		}
		log.Printf("[INFO] Getting data for predefined controls name %s\n", name)
		res, _, err := zClient.inspection_predefined_controls.GetByName(name, version)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("action", resp.Action)
		_ = d.Set("action_value", resp.ActionValue)
		_ = d.Set("attachment", resp.Attachment)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("control_group", resp.ControlGroup)
		_ = d.Set("control_number", resp.ControlNumber)
		_ = d.Set("default_action", resp.DefaultAction)
		_ = d.Set("default_action_value", resp.DefaultActionValue)
		_ = d.Set("description", resp.Description)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("paranoia_level", resp.ParanoiaLevel)
		_ = d.Set("severity", resp.Severity)
		_ = d.Set("version", resp.Version)

	} else {
		return fmt.Errorf("couldn't find any predefined inspection controls with name '%s' or id '%s'", name, id)
	}

	return nil
}
