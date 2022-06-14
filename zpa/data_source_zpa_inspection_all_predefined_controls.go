package zpa

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/inspectioncontrol/inspection_predefined_controls"
)

func dataSourceInspectionAllPredefinedControls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInspectionAllPredefinedControlsRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInspectionAllPredefinedControlsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	version, versionSet := d.Get("version").(string)
	if !versionSet || version == "" {
		return fmt.Errorf("when the name is set, version must be set as well")
	}
	list, err := zClient.inspection_predefined_controls.GetAll(version)
	if err != nil {
		return err
	}
	d.SetId("predefined_controls")
	_ = d.Set("list", flattenList(list))
	return nil
}

func flattenList(list []inspection_predefined_controls.PredefinedControls) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, control := range list {
		result = append(result, map[string]interface{}{
			"action":                              control.Action,
			"action_value":                        control.ActionValue,
			"attachment":                          control.Attachment,
			"creation_time":                       control.CreationTime,
			"control_group":                       control.ControlGroup,
			"control_number":                      control.ControlNumber,
			"default_action":                      control.DefaultAction,
			"default_action_value":                control.DefaultActionValue,
			"description":                         control.Description,
			"modifiedby":                          control.ModifiedBy,
			"modified_time":                       control.ModifiedTime,
			"name":                                control.Name,
			"paranoia_level":                      control.ParanoiaLevel,
			"severity":                            control.Severity,
			"version":                             control.Version,
			"associated_inspection_profile_names": flattenInspectionProfileNames(control.AssociatedInspectionProfileNames),
		})
	}
	return result
}
