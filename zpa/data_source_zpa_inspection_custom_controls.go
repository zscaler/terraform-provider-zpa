package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/inspectioncontrol/inspection_custom_controls"
)

func dataSourceInspectionCustomControls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInspectionCustomControlsRead,
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
			"rules": dataInspectionRulesSchema(),
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
		},
	}
}

func dataSourceInspectionCustomControlsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *inspection_custom_controls.InspectionCustomControl
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for custom inspection control %s\n", id)
		res, _, err := zClient.inspection_custom_controls.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for custom inspection control name %s\n", name)
		res, _, err := zClient.inspection_custom_controls.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("action", resp.Action)
		_ = d.Set("action_value", resp.ActionValue)
		_ = d.Set("control_number", resp.ControlNumber)
		_ = d.Set("control_rule_json", resp.ControlRuleJson)
		_ = d.Set("creation_time", resp.CreationTime)
		_ = d.Set("default_action", resp.DefaultAction)
		_ = d.Set("default_action_value", resp.DefaultActionValue)
		_ = d.Set("description", resp.Description)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("paranoia_level", resp.ParanoiaLevel)
		_ = d.Set("severity", resp.Severity)
		_ = d.Set("version", resp.Version)
		_ = d.Set("type", resp.Type)

		if err := d.Set("rules", flattenInspectionRules(resp.Rules)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any predefined inspection controls with name '%s' or id '%s'", name, id)
	}

	return nil
}
