package zpa

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

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
						"REQUIRE_APPROVAL",
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
					"APP",
					"APP_GROUP",
					"LOCATION",
					"IDP",
					"SAML",
					"SCIM",
					"SCIM_GROUP",
					"CLIENT_TYPE",
					"POSTURE",
					"TRUSTED_NETWORK",
					"BRANCH_CONNECTOR_GROUP",
					"EDGE_CONNECTOR_GROUP",
					"MACHINE_GRP",
					"COUNTRY_CODE",
					"PLATFORM",
					"RISK_FACTOR_TYPE",
				}),
			},
		),
	}
}

func resourcePolicyAccessCreate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(zClient, "ACCESS_POLICY", microTenantID)
		if err != nil {
			return err
		}
	}

	req, err := expandCreatePolicyRule(d, policySetID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating ZPA policy access rule with request\n%+v\n", req)

	if err := ValidateConditions(req.Conditions, zClient, microTenantID); err != nil {
		return err
	}

	resp, _, err := policysetcontroller.CreateRule(service, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePolicyAccessRead(d, meta)
}

func resourcePolicyAccessRead(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	microTenantID := GetString(d.Get("microtenant_id"))

	policySetID, err := fetchPolicySetIDByType(zClient, "ACCESS_POLICY", microTenantID)
	if err != nil {
		return err
	}

	service := zClient.PolicySetController
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	log.Printf("[INFO] Getting Policy Set Rule: policySetID:%s id: %s\n", policySetID, d.Id())
	resp, respErr, err := policysetcontroller.GetPolicyRule(service, policySetID, d.Id())
	if err != nil {
		if respErr != nil && (respErr.StatusCode == 404 || respErr.StatusCode == http.StatusNotFound) {
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
	_ = d.Set("policy_set_id", policySetID)
	_ = d.Set("policy_type", "ACCESS_POLICY")
	_ = d.Set("priority", resp.Priority)
	_ = d.Set("lss_default_rule", resp.LSSDefaultRule)
	_ = d.Set("microtenant_id", microTenantID)
	_ = d.Set("conditions", flattenPolicyConditions(resp.Conditions))
	_ = d.Set("app_server_groups", flattenCommonAppServerGroups(resp.AppServerGroups))
	_ = d.Set("app_connector_groups", flattenCommonAppConnectorGroups(resp.AppConnectorGroups))

	return nil
}

func resourcePolicyAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(zClient, "ACCESS_POLICY", microTenantID)
		if err != nil {
			return err
		}
	}

	ruleID := d.Id()
	req, err := expandCreatePolicyRule(d, policySetID)
	if err != nil {
		return err
	}

	if err := ValidateConditions(req.Conditions, zClient, microTenantID); err != nil {
		return err
	}

	if _, err := policysetcontroller.UpdateRule(service, policySetID, ruleID, req); err != nil {
		return err
	}

	return resourcePolicyAccessRead(d, meta)
}

func resourcePolicyAccessDelete(d *schema.ResourceData, meta interface{}) error {
	zClient := meta.(*Client)
	service := zClient.PolicySetController

	microTenantID := GetString(d.Get("microtenant_id"))
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	var policySetID string
	var err error

	if v, ok := d.GetOk("policy_set_id"); ok {
		policySetID = v.(string)
	} else {
		policySetID, err = fetchPolicySetIDByType(zClient, "ACCESS_POLICY", microTenantID)
		if err != nil {
			return err
		}
	}

	log.Printf("[INFO] Deleting policy set rule with id %v\n", d.Id())

	if _, err := policysetcontroller.Delete(service, policySetID, d.Id()); err != nil {
		return err
	}

	return nil
}

func expandCreatePolicyRule(d *schema.ResourceData, policySetID string) (*policysetcontroller.PolicyRule, error) {

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
		MicroTenantID:      d.Get("microtenant_id").(string),
		LSSDefaultRule:     d.Get("lss_default_rule").(bool),
		Conditions:         conditions,
		AppServerGroups:    expandCommonServerGroups(d),
		AppConnectorGroups: expandCommonAppConnectorGroups(d),
	}, nil
}

/*
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

func flattenPolicyRuleServerGroups(appServerGroup []servergroup.AppServerGroups) []interface{} {
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
*/
