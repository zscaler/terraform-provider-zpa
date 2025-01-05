package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_predefined_controls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInspectionPredefinedControls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInspectionPredefinedControlsRead,
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
			"protocol_type": {
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

func dataSourceInspectionPredefinedControlsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *inspection_predefined_controls.PredefinedControls
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for predefined controls %s\n", id)
		res, _, err := inspection_predefined_controls.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		version, versionSet := d.Get("version").(string)
		if !versionSet || version == "" {
			return diag.FromErr(fmt.Errorf("when the name is set, version must be set as well"))
		}
		log.Printf("[INFO] Getting data for predefined controls name %s\n", name)
		res, _, err := inspection_predefined_controls.GetByName(ctx, service, name, version)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
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
		_ = d.Set("control_type", resp.ControlType)
		_ = d.Set("default_action", resp.DefaultAction)
		_ = d.Set("default_action_value", resp.DefaultActionValue)
		_ = d.Set("description", resp.Description)
		_ = d.Set("modifiedby", resp.ModifiedBy)
		_ = d.Set("modified_time", resp.ModifiedTime)
		_ = d.Set("name", resp.Name)
		_ = d.Set("paranoia_level", resp.ParanoiaLevel)
		_ = d.Set("protocol_type", resp.ProtocolType)
		_ = d.Set("severity", resp.Severity)
		_ = d.Set("version", resp.Version)
		_ = d.Set("associated_inspection_profile_names", flattenInspectionProfileNames(resp.AssociatedInspectionProfileNames))
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any predefined inspection controls with name '%s' or id '%s'", name, id))
	}

	return nil
}

func flattenInspectionProfileNames(names []common.AssociatedProfileNames) []interface{} {
	flattenedList := make([]interface{}, len(names))
	for i, val := range names {
		flattenedList[i] = map[string]interface{}{
			"id":   val.ID,
			"name": val.Name,
		}
	}
	return flattenedList
}
