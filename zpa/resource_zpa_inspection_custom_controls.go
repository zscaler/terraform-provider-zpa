package zpa

import (
	"errors"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/client"
	"github.com/zscaler/terraform-provider-zpa/gozscaler/inspectioncontrol/inspection_custom_controls"
)

func resourceInspectionCustomControls() *schema.Resource {
	return &schema.Resource{
		Create: resourceInspectionCustomControlsCreate,
		Read:   resourceInspectionCustomControlsRead,
		Update: resourceInspectionCustomControlsUpdate,
		Delete: resourceInspectionCustomControlsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("custom_id", id)
				} else {
					resp, _, err := zClient.inspection_custom_controls.GetByName(id)
					if err == nil {
						d.SetId(resp.ID)
						d.Set("custom_id", resp.ID)
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
				Required: true,
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The performed action",
				ValidateFunc: validation.StringInSlice([]string{
					"PASS",
					"BLOCK",
					"REDIRECT",
				}, false),
			},
			"action_value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"associated_inspection_profile_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Name of the inspection profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"control_number": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"control_rule_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The control rule in JSON format that has the conditions and type of control for the inspection control",
			},
			"default_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The performed action",
				ValidateFunc: validation.StringInSlice([]string{
					"PASS",
					"BLOCK",
					"REDIRECT",
				}, false),
			},
			"default_action_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This is used to provide the redirect URL if the default action is set to REDIRECT",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the custom control",
			},
			"paranoia_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "OWASP Predefined Paranoia Level. Range: [1-4], inclusive",
			},
			"rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Rules of the custom controls applied as conditions (JSON)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"names": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Name of the rules. If rules.type is set to REQUEST_HEADERS, REQUEST_COOKIES, or RESPONSE_HEADERS, the rules.name field is required.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Type value for the rules. ",
							ValidateFunc: validation.StringInSlice([]string{
								"REQUEST_HEADERS",
								"REQUEST_URI",
								"QUERY_STRING",
								"REQUEST_COOKIES",
								"REQUEST_METHOD",
								"REQUEST_BODY",
								"RESPONSE_HEADERS",
								"RESPONSE_BODY",
							}, false),
						},
						"conditions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lhs": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Signifies the key for the object type",
										ValidateFunc: validation.StringInSlice([]string{
											"SIZE",
											"VALUE",
										}, false),
									},
									"op": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Denotes the operation type.",
										ValidateFunc: validation.StringInSlice([]string{
											"RX",
											"EQ",
											"LE",
											"GE",
											"CONTAINS",
											"STARTS_WITH",
											"ENDS_WITH",
										}, false),
									},
									"rhs": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Denotes the value for the given object type. Its value depends on the key.",
									},
								},
							},
						},
					},
				},
			},
			"severity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Severity of the control number",
				ValidateFunc: validation.StringInSlice([]string{
					"CRITICAL",
					"ERROR",
					"WARNING",
					"INFO",
				}, false),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rules to be applied to the request or response type",
				ValidateFunc: validation.StringInSlice([]string{
					"REQUEST",
					"RESPONSE",
				}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceInspectionCustomControlsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandInspectionCustomControls(d)
	log.Printf("[INFO] Creating custom inspection control with request\n%+v\n", req)
	if err := valdateRules(req); err != nil {
		return err
	}
	resp, _, err := zClient.inspection_custom_controls.Create(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created custom inspection control request. ID: %v\n", resp)

	d.SetId(resp.ID)
	return resourceInspectionCustomControlsRead(d, m)
}

func resourceInspectionCustomControlsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.inspection_custom_controls.Get(d.Id())
	if err != nil {
		if errResp, ok := err.(*client.ErrorResponse); ok && errResp.IsObjectNotFound() {
			log.Printf("[WARN] Removing custom inspection control %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	log.Printf("[INFO] Getting custom inspection control:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("action", resp.Action)
	_ = d.Set("action_value", resp.ActionValue)
	_ = d.Set("control_number", resp.ControlNumber)
	_ = d.Set("control_rule_json", resp.ControlRuleJson)
	_ = d.Set("default_action", resp.DefaultAction)
	_ = d.Set("default_action_value", resp.DefaultActionValue)
	_ = d.Set("description", resp.Description)
	_ = d.Set("name", resp.Name)
	_ = d.Set("paranoia_level", resp.ParanoiaLevel)
	_ = d.Set("severity", resp.Severity)
	_ = d.Set("version", resp.Version)
	_ = d.Set("type", resp.Type)

	if err := d.Set("rules", flattenInspectionCustomRules(resp.Rules)); err != nil {
		return err
	}
	return nil
}

func resourceInspectionCustomControlsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating custom inspection control ID: %v\n", id)
	req := expandInspectionCustomControls(d)
	if err := valdateRules(req); err != nil {
		return err
	}
	if _, err := zClient.inspection_custom_controls.Update(id, &req); err != nil {
		return err
	}

	return resourceInspectionCustomControlsRead(d, m)
}

func resourceInspectionCustomControlsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting custom inspection control ID: %v\n", d.Id())

	if _, err := zClient.inspection_custom_controls.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] custom inspection control deleted")
	return nil
}

func expandInspectionCustomControls(d *schema.ResourceData) inspection_custom_controls.InspectionCustomControl {
	custom_control := inspection_custom_controls.InspectionCustomControl{
		Action:                           d.Get("action").(string),
		ActionValue:                      d.Get("action_value").(string),
		ControlNumber:                    d.Get("control_number").(string),
		ControlRuleJson:                  d.Get("control_rule_json").(string),
		DefaultAction:                    d.Get("default_action").(string),
		DefaultActionValue:               d.Get("default_action_value").(string),
		Description:                      d.Get("description").(string),
		Name:                             d.Get("name").(string),
		ParanoiaLevel:                    d.Get("paranoia_level").(string),
		Severity:                         d.Get("severity").(string),
		Type:                             d.Get("type").(string),
		Version:                          d.Get("version").(string),
		AssociatedInspectionProfileNames: expandAssociatedInspectionProfileNames(d),
		Rules:                            expandInspectionCustomControlsRules(d),
	}
	return custom_control
}

func expandAssociatedInspectionProfileNames(d *schema.ResourceData) []inspection_custom_controls.AssociatedProfileNames {
	inspectionProfileInterface, ok := d.GetOk("associated_inspection_profile_names")
	if ok {
		inspectionProfile := inspectionProfileInterface.(*schema.Set)
		log.Printf("[INFO] associated inspection profile names data: %+v\n", inspectionProfile)
		var inspectionProfiles []inspection_custom_controls.AssociatedProfileNames
		for _, inspectionProfile := range inspectionProfile.List() {
			inspectionProfile, ok := inspectionProfile.(map[string]interface{})
			if ok {
				for _, id := range inspectionProfile["id"].(*schema.Set).List() {
					inspectionProfiles = append(inspectionProfiles, inspection_custom_controls.AssociatedProfileNames{
						ID: id.(string),
					})
				}
			}
		}
		return inspectionProfiles
	}

	return []inspection_custom_controls.AssociatedProfileNames{}
}

func expandInspectionCustomControlsRules(d *schema.ResourceData) []inspection_custom_controls.Rules {
	rulesObj, ok := d.GetOk("rules")
	if !ok {
		return nil
	}
	rulesInterfaces := rulesObj.([]interface{})
	var rules []inspection_custom_controls.Rules
	for _, ruleObj := range rulesInterfaces {
		ruleMap, ok := ruleObj.(map[string]interface{})
		if !ok {
			continue
		}
		var names []string
		ruleNamesSet, ok := ruleMap["names"].(*schema.Set)
		if ok {
			for _, name := range ruleNamesSet.List() {
				names = append(names, name.(string))
			}

		}
		rules = append(rules, inspection_custom_controls.Rules{
			Names:      names,
			Type:       ruleMap["type"].(string),
			Conditions: expandCustomControlRuleConditions(ruleMap["conditions"]),
		})
	}
	return rules
}

func expandCustomControlRuleConditions(conditionsObj interface{}) []inspection_custom_controls.Conditions {
	conditionsInterface, ok := conditionsObj.([]interface{})
	if !ok {
		return nil
	}
	var conditions []inspection_custom_controls.Conditions
	for _, conditionObj := range conditionsInterface {
		condition, ok := conditionObj.(map[string]interface{})
		if !ok {
			continue
		}
		conditions = append(conditions, inspection_custom_controls.Conditions{
			LHS: condition["lhs"].(string),
			RHS: condition["rhs"].(string),
			OP:  condition["op"].(string),
		})
	}

	return conditions
}

func valdateRules(customCtl inspection_custom_controls.InspectionCustomControl) error {

	for _, rule := range customCtl.Rules {
		if customCtl.Type == "RESPONSE" {
			if rule.Type != "RESPONSE_HEADERS" && rule.Type != "RESPONSE_BODY" {
				return errors.New("when type == RESPONSE rules.type must be: RESPONSE_HEADERS || RESPONSE_BODY")
			}
		} else if customCtl.Type == "REQUEST" {
			if (rule.Type == "REQUEST_HEADERS" || rule.Type == "REQUEST_COOKIES") && len(rule.Names) == 0 {
				return errors.New("when type == REQUEST and rules.type is: REQUEST_HEADERS || REQUEST_COOKIES the rules.names must be set")
			}
			if (rule.Type == "REQUEST_URI" || rule.Type == "QUERY_STRING" || rule.Type == "REQUEST_BODY" || rule.Type == "REQUEST_METHOD") && len(rule.Names) > 0 {
				return errors.New("when type == REQUEST and rules.type is: REQUEST_URI || QUERY_STRING || REQUEST_BODY || REQUEST_METHOD the rules.name is not allowed")
			}
		}
		for _, cond := range rule.Conditions {
			if in(rule.Type, []string{"REQUEST_HEADERS", "REQUEST_COOKIES", "REQUEST_URI", "QUERY_STRING", "REQUEST_BODY"}) {
				if cond.LHS == "SIZE" && (!in(cond.OP, []string{"EQ", "LE", "GE"}) || !isNumber(cond.RHS)) {
					return errors.New("when rules.type is: " + rule.Type + " the conditions.lhs must be == SIZE && conditions.op == EQ, LE, GE && condition.rhs must be a number(string)")
				}
				if cond.LHS == "VALUE" && (!in(cond.OP, []string{"CONTAINS", "STARTS_WITH", "ENDS_WITH", "RX"})) {
					return errors.New("when rules.type is: " + rule.Type + " the conditions.lhs must be == VALUE && conditions.op must be == CONTAINS, STARTS_WITH, ENDS_WITH, RX and rhs must be a string value")
				}
			}
			if rule.Type == "REQUEST_METHOD" {
				if cond.LHS == "SIZE" && (!in(cond.OP, []string{"EQ", "LE", "GE"}) || !isNumber(cond.RHS)) {
					return errors.New("when rules.type is: " + rule.Type + " the conditions.lhs must be == SIZE && conditions.op == EQ, LE, GE && condition.rhs must be a number(string)")
				}
				if cond.LHS == "VALUE" && (!in(cond.OP, []string{"CONTAINS", "STARTS_WITH", "ENDS_WITH", "RX"}) || !in(cond.RHS, []string{"GET", "POST", "PUT", "PATCH", "CONNECT", "HEAD", "OPTIONS", "DELETE", "TRACE"})) {
					return errors.New("when rules.type is: " + rule.Type + " the conditions.lhs must be == VALUE && conditions.op must be == CONTAINS, STARTS_WITH, ENDS_WITH, RX && condition.rhs== GET,POST,PUT,PATCH,CONNECT,HEAD,OPTIONS,DELETE,TRACE")
				}
			}
		}

	}
	return nil
}

func isNumber(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}

func in(val string, list []string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
