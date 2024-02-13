package zpa

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

func resourcePolicyRedictionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyRedictionRuleCreate,
		Read:   resourcePolicyRedictionRuleRead,
		Update: resourcePolicyRedictionRuleUpdate,
		Delete: resourcePolicyRedictionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"REDIRECTION_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(),
			map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "  This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"REDIRECT_DEFAULT",
						"REDIRECT_PREFERRED",
						"REDIRECT_ALWAYS",
					}, false),
				},
				"conditions": GetPolicyConditionsSchema([]string{
					"CLIENT_TYPE",
				}),
				"service_edge_groups": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					Description: "List of the service edge group IDs.",
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
			},
		),
	}
}

// validatePolicyRedirectionRuleAction validates the "action" attribute against "service_edge_groups" requirements.
func validatePolicyRedirectionRuleAction(d *schema.ResourceData) error {
	action := d.Get("action").(string)
	serviceEdgeGroups := d.Get("service_edge_groups").(*schema.Set).List()

	switch action {
	case "REDIRECT_PREFERRED", "REDIRECT_ALWAYS":
		if len(serviceEdgeGroups) == 0 {
			return fmt.Errorf("one or more ZPA Private Service Edge groups must be selected when the Private Service Edge Selection Method is %s", action)
		}
	case "REDIRECT_DEFAULT":
		if len(serviceEdgeGroups) > 0 {
			return fmt.Errorf("zpa Private Service Edge groups must be empty when the Private Service Edge Selection Method is REDIRECT_DEFAULT")
		}
	}

	return nil
}

func resourcePolicyRedictionRuleCreate(d *schema.ResourceData, m interface{}) error {
	// Validate the "action" and "service_edge_groups" attributes
	if err := validatePolicyRedirectionRuleAction(d); err != nil {
		return err
	}

	zClient := m.(*Client)
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))

	req, err := expandCreatePolicyRedirectionRule(d)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy redirection rule with request\n%+v\n", req)
	if err := ValidateConditions(req.Conditions, zClient, req.MicroTenantID); err == nil {
		policysetcontroller, _, err := service.Create(req)
		if err != nil {
			return err
		}
		d.SetId(policysetcontroller.ID)

		return resourcePolicyRedictionRuleRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy redirection (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}
}

func resourcePolicyRedictionRuleRead(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))

	globalPolicySet, _, err := service.GetByPolicyType("REDIRECTION_POLICY")
	if err != nil {
		return err
	}
	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", globalPolicySet.ID, d.Id())
	resp, _, err := service.GetPolicyRule(globalPolicySet.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Got Policy Set Redirection Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("action", resp.Action)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))
	_ = d.Set("service_edge_groups", flattenPolicyRuleServiceEdgeGroups(resp.ServiceEdgeGroups))
	return nil
}

func resourcePolicyRedictionRuleUpdate(d *schema.ResourceData, m interface{}) error {

	// Validate the "action" and "service_edge_groups" attributes
	if err := validatePolicyRedirectionRuleAction(d); err != nil {
		return err
	}

	zClient := m.(*Client)
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("REDIRECTION_POLICY")
	if err != nil {
		return err
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy redirection rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyRedirectionRule(d)
	if err != nil {
		return err
	}
	if err := ValidateConditions(req.Conditions, zClient, req.MicroTenantID); err == nil {
		if _, _, err := service.GetPolicyRule(globalPolicySet.ID, ruleID); err != nil {
			if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
				d.SetId("")
				return nil
			}
		}

		if _, err := service.Update(globalPolicySet.ID, ruleID, req); err != nil {
			return err
		}

		return resourcePolicyRedictionRuleRead(d, m)
	} else {
		return fmt.Errorf("couldn't validate the zpa policy redirection (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}
}

func resourcePolicyRedictionRuleDelete(d *schema.ResourceData, m interface{}) error {
	service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))
	globalPolicySet, _, err := service.GetByPolicyType("REDIRECTION_POLICY")
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy redirection rule with id %v\n", d.Id())

	if _, err := service.Delete(globalPolicySet.ID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyRedirectionRule(d *schema.ResourceData) (*policysetcontroller.PolicyRule, error) {
	policySetID, ok := d.Get("policy_set_id").(string)
	if !ok {
		return nil, fmt.Errorf("policy_set_id is not set")
	}
	log.Printf("[INFO] action_id:%v\n", d.Get("action_id"))
	conditions, err := ExpandPolicyConditions(d)
	if err != nil {
		return nil, err
	}
	return &policysetcontroller.PolicyRule{
		ID:                d.Get("id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Action:            d.Get("action").(string),
		ActionID:          d.Get("action_id").(string),
		CustomMsg:         d.Get("custom_msg").(string),
		Operator:          d.Get("operator").(string),
		PolicySetID:       policySetID,
		PolicyType:        d.Get("policy_type").(string),
		Priority:          d.Get("priority").(string),
		MicroTenantID:     GetString(d.Get("microtenant_id")),
		Conditions:        conditions,
		ServiceEdgeGroups: expandPolicysetControllerServiceEdgeGroups(d),
	}, nil
}

func expandPolicysetControllerServiceEdgeGroups(d *schema.ResourceData) []policysetcontroller.ServiceEdgeGroups {
	serviceEdgeGroupsInterface, ok := d.GetOk("service_edge_groups")
	if ok {
		edgeGroup := serviceEdgeGroupsInterface.(*schema.Set)
		log.Printf("[INFO] service edge groups data: %+v\n", edgeGroup)
		var edgeGroups []policysetcontroller.ServiceEdgeGroups
		for _, edgeGroup := range edgeGroup.List() {
			edgeGroup, _ := edgeGroup.(map[string]interface{})
			if edgeGroup != nil {
				for _, id := range edgeGroup["id"].(*schema.Set).List() {
					edgeGroups = append(edgeGroups, policysetcontroller.ServiceEdgeGroups{
						ID: id.(string),
					})
				}
			}
		}
		return edgeGroups
	}

	return []policysetcontroller.ServiceEdgeGroups{}
}

func flattenPolicyRuleServiceEdgeGroups(serviceEdgeGroup []policysetcontroller.ServiceEdgeGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(serviceEdgeGroup))
	for i, serviceEdgeGroup := range serviceEdgeGroup {
		ids[i] = serviceEdgeGroup.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}
