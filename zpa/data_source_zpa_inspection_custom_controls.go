package zpa

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInspectionCustomControls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInspectionCustomControlsRead,
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
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lhs": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"op": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rhs": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"names": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Name of the rules. If rules.type is set to REQUEST_HEADERS, REQUEST_COOKIES, or RESPONSE_HEADERS, the rules.name field is required.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
		},
	}
}

func dataSourceInspectionCustomControlsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *inspection_custom_controls.InspectionCustomControl
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for custom inspection control %s\n", id)
		res, _, err := inspection_custom_controls.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res
	}
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for custom inspection control name %s\n", name)
		res, _, err := inspection_custom_controls.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err) // Wrap error using diag.FromErr
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("action", resp.Action)
		_ = d.Set("action_value", resp.ActionValue)
		_ = d.Set("control_number", resp.ControlNumber)
		_ = d.Set("control_rule_json", resp.ControlRuleJson)
		_ = d.Set("control_type", resp.ControlType)
		_ = d.Set("creation_time", resp.CreationTime)
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
		_ = d.Set("type", resp.Type)

		if err := d.Set("rules", flattenInspectionCustomRules(resp.Rules)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any custom inspection controls with name '%s' or id '%s'", name, id))
	}

	return nil
}

func flattenInspectionCustomRules(rule []inspection_custom_controls.Rules) []interface{} {
	rules := make([]interface{}, len(rule))
	for i, rule := range rule {
		rules[i] = map[string]interface{}{
			"conditions": flattenInspectionRuleConditions(rule),
			"names":      rule.Names,
			"type":       rule.Type,
		}
	}

	return rules
}

func flattenInspectionRuleConditions(condition inspection_custom_controls.Rules) []interface{} {
	conditions := make([]interface{}, len(condition.Conditions))
	for i, val := range condition.Conditions {
		conditions[i] = map[string]interface{}{
			"lhs": val.LHS,
			"rhs": val.RHS,
			"op":  val.OP,
		}
	}

	return conditions
}
