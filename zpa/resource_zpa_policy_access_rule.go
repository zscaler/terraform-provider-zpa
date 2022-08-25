package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/policysetcontroller"
)

func validateAccessPolicyRuleOrder(order string, zClient *Client) error {
	o, err := strconv.Atoi(order)
	if err != nil || o < 1 {
		return fmt.Errorf("order must be a valid number >= 1")
	}
	policy, _, err := zClient.policysetcontroller.GetByNameAndType("ACCESS_POLICY", "Zscaler Deception")
	if err == nil && policy != nil && o == 1 {
		return fmt.Errorf("policy Zscaler Deception exists, orders must be start from 2")
	}
	return nil
}

func resourcePolicyAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyAccessCreate,
		Read:   resourcePolicyAccessRead,
		Update: resourcePolicyAccessUpdate,
		Delete: resourcePolicyAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importPolicyStateContextFunc([]string{"ACCESS_POLICY", "GLOBAL_POLICY"}),
		},

		Schema: MergeSchema(
			CommonPolicySchema(), map[string]*schema.Schema{
				"action": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "  This is for providing the rule action.",
					ValidateFunc: validation.StringInSlice([]string{
						"ALLOW",
						"DENY",
					}, false),
				},
				"app_server_groups": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					Description: "List of the server group IDs.",
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
				"app_connector_groups": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					Description: "List of app-connector IDs.",
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
				"conditions": GetPolicyConditionsSchema([]string{
					"USER",
					"USER_GROUP",
					"LOCATION",
					"APP",
					"APP_GROUP",
					"SAML",
					"POSTURE",
					"CLIENT_TYPE",
					"IDP",
					"TRUSTED_NETWORK",
					"EDGE_CONNECTOR_GROUP",
					"MACHINE_GRP",
					"SCIM",
					"SCIM_GROUP",
				}),
			},
		),
	}
}

func resourcePolicyAccessCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req, err := expandCreatePolicyRule(d)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating zpa policy rule with request\n%+v\n", req)
	if err := validateAccessPolicyRuleOrder(req.RuleOrder, zClient); err != nil {
		return err
	}
	if !ValidateConditions(req.Conditions, zClient) {
		return fmt.Errorf("couldn't validate the zpa policy rule (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}
	policysetcontroller, _, err := zClient.policysetcontroller.Create(req)
	if err != nil {
		return err
	}
	d.SetId(policysetcontroller.ID)
	order, ok := d.GetOk("rule_order")
	if ok {
		reorder(order, policysetcontroller.PolicySetID, "ACCESS_POLICY", policysetcontroller.ID, zClient)
	}
	return resourcePolicyAccessRead(d, m)
}

func resourcePolicyAccessRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	globalPolicySet, _, err := zClient.policysetcontroller.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return err
	}
	log.Printf("[INFO] Getting Policy Set Rule: globalPolicySet:%s id: %s\n", globalPolicySet.ID, d.Id())
	resp, _, err := zClient.policysetcontroller.GetPolicyRule(globalPolicySet.ID, d.Id())
	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing policy rule %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Got Policy Set Rule:\n%+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("description", resp.Description)
	_ = d.Set("name", resp.Name)
	_ = d.Set("action", resp.Action)
	_ = d.Set("action_id", resp.ActionID)
	_ = d.Set("custom_msg", resp.CustomMsg)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("operator", resp.Operator)
	_ = d.Set("policy_set_id", resp.PolicySetID)
	_ = d.Set("policy_type", resp.PolicyType)
	_ = d.Set("priority", resp.Priority)
	_ = d.Set("rule_order", resp.RuleOrder)
	_ = d.Set("lss_default_rule", resp.LSSDefaultRule)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))
	_ = d.Set("app_server_groups", flattenPolicyRuleServerGroups(resp.AppServerGroups))
	_ = d.Set("app_connector_groups", flattenPolicyRuleAppConnectorGroups(resp.AppConnectorGroups))

	return nil
}

func resourcePolicyAccessUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicySet, _, err := zClient.policysetcontroller.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return err
	}
	ruleID := d.Id()
	log.Printf("[INFO] Updating policy rule ID: %v\n", ruleID)
	req, err := expandCreatePolicyRule(d)
	if err != nil {
		return err
	}
	if err := validateAccessPolicyRuleOrder(req.RuleOrder, zClient); err != nil {
		return err
	}
	if !ValidateConditions(req.Conditions, zClient) {
		return fmt.Errorf("couldn't validate the zpa policy rule (%s) operands, please make sure you are using valid inputs for APP type, LHS & RHS", req.Name)
	}
	if _, err := zClient.policysetcontroller.Update(globalPolicySet.ID, ruleID, req); err != nil {
		return err
	}
	if d.HasChange("rule_order") {
		order, ok := d.GetOk("rule_order")
		if ok {
			reorder(order, globalPolicySet.ID, "ACCESS_POLICY", ruleID, zClient)
		}
	}
	return resourcePolicyAccessRead(d, m)
}

func resourcePolicyAccessDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	globalPolicySet, _, err := zClient.policysetcontroller.GetByPolicyType("ACCESS_POLICY")
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	if _, err := zClient.policysetcontroller.Delete(globalPolicySet.ID, d.Id()); err != nil {
		return err
	}

	return nil

}

func expandCreatePolicyRule(d *schema.ResourceData) (*policysetcontroller.PolicyRule, error) {
	policySetID, ok := d.Get("policy_set_id").(string)
	if !ok {
		log.Printf("[ERROR] policy_set_id is not set\n")
		return nil, fmt.Errorf("policy_set_id is not set")
	}
	log.Printf("[INFO] action_id:%v\n", d.Get("action_id"))
	conditions, err := ExpandPolicyConditions(d)
	if err != nil {
		return nil, err
	}
	return &policysetcontroller.PolicyRule{
		ID:                 d.Get("id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Action:             d.Get("action").(string),
		ActionID:           d.Get("action_id").(string),
		BypassDefaultRule:  d.Get("bypass_default_rule").(bool),
		CustomMsg:          d.Get("custom_msg").(string),
		DefaultRule:        d.Get("default_rule").(bool),
		Operator:           d.Get("operator").(string),
		PolicySetID:        policySetID,
		PolicyType:         d.Get("policy_type").(string),
		Priority:           d.Get("priority").(string),
		RuleOrder:          d.Get("rule_order").(string),
		LSSDefaultRule:     d.Get("lss_default_rule").(bool),
		Conditions:         conditions,
		AppServerGroups:    expandPolicySetControllerAppServerGroups(d),
		AppConnectorGroups: expandPolicysetControllerAppConnectorGroups(d),
	}, nil
}

func expandPolicySetControllerAppServerGroups(d *schema.ResourceData) []policysetcontroller.AppServerGroups {
	appServerGroupsInterface, ok := d.GetOk("app_server_groups")
	if ok {
		appServer := appServerGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app server groups data: %+v\n", appServer)
		var appServerGroups []policysetcontroller.AppServerGroups
		for _, appServerGroup := range appServer.List() {
			appServerGroup, _ := appServerGroup.(map[string]interface{})
			if appServerGroup != nil {
				for _, id := range appServerGroup["id"].(*schema.Set).List() {
					appServerGroups = append(appServerGroups, policysetcontroller.AppServerGroups{
						ID: id.(string),
					})
				}
			}
		}
		return appServerGroups
	}

	return []policysetcontroller.AppServerGroups{}
}

func expandPolicysetControllerAppConnectorGroups(d *schema.ResourceData) []policysetcontroller.AppConnectorGroups {
	appConnectorGroupsInterface, ok := d.GetOk("app_connector_groups")
	if ok {
		appConnector := appConnectorGroupsInterface.(*schema.Set)
		log.Printf("[INFO] app connector groups data: %+v\n", appConnector)
		var appConnectorGroups []policysetcontroller.AppConnectorGroups
		for _, appConnectorGroup := range appConnector.List() {
			appConnectorGroup, _ := appConnectorGroup.(map[string]interface{})
			if appConnectorGroup != nil {
				for _, id := range appConnectorGroup["id"].(*schema.Set).List() {
					appConnectorGroups = append(appConnectorGroups, policysetcontroller.AppConnectorGroups{
						ID: id.(string),
					})
				}

			}
		}
		return appConnectorGroups
	}

	return []policysetcontroller.AppConnectorGroups{}
}

func flattenPolicyRuleServerGroups(appServerGroup []policysetcontroller.AppServerGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(appServerGroup))
	for i, serverGroup := range appServerGroup {
		ids[i] = serverGroup.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenPolicyRuleAppConnectorGroups(appConnectorGroups []policysetcontroller.AppConnectorGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]string, len(appConnectorGroups))
	for i, group := range appConnectorGroups {
		ids[i] = group.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}
