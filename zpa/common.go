package zpa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
)

var policySets = map[string]policysetcontroller.PolicySet{}
var policySetsMutex sync.Mutex

// Common values shared between Service Edge Groups and App Connector Groups
var versionProfileNameIDMapping map[string]string = map[string]string{
	"Default":          "0",
	"Previous Default": "1",
	"New Release":      "2",
}

// Common validation shared between Service Edge Groups and App Connector Groups
func validateAndSetProfileNameID(d *schema.ResourceData) error {
	_, versionProfileIDSet := d.GetOk("version_profile_id")
	versionProfileName, versionProfileNameSet := d.GetOk("version_profile_name")
	if versionProfileNameSet && d.HasChange("version_profile_name") {
		if id, ok := versionProfileNameIDMapping[versionProfileName.(string)]; ok {
			d.Set("version_profile_id", id)
		}
		return nil
	}
	if !versionProfileNameSet && !versionProfileIDSet {
		return errors.New("one of version_profile_id or version_profile_name must be set")
	}
	return nil
}

func ValidateConditions(conditions []policysetcontroller.Conditions, zClient *Client, microtenantID string) error {
	for _, condition := range conditions {
		if err := validateOperands(condition.Operands, zClient, microtenantID); err != nil {
			return err
		}
	}
	return nil
}
func validateOperands(operands []policysetcontroller.Operands, zClient *Client, microtenantID string) error {
	for _, operand := range operands {
		if err := validateOperand(operand, zClient, microtenantID); err != nil {
			return err
		}
	}
	return nil
}

func validateOperand(operand policysetcontroller.Operands, zClient *Client, microtenantID string) error {
	switch operand.ObjectType {
	case "APP":
		return customValidate(operand, []string{"id"}, "application segment ID", Getter(func(id string) error {
			_, _, err := zClient.applicationsegment.WithMicroTenant(microtenantID).Get(id)
			return err
		}))
	case "APP_GROUP":
		return customValidate(operand, []string{"id"}, "Segment Group ID", Getter(func(id string) error {
			_, _, err := zClient.segmentgroup.WithMicroTenant(microtenantID).Get(id)
			return err
		}))

	case "IDP":
		return customValidate(operand, []string{"id"}, "IDP ID", Getter(func(id string) error {
			_, _, err := zClient.idpcontroller.Get(id)
			return err
		}))
	case "EDGE_CONNECTOR_GROUP":
		return customValidate(operand, []string{"id"}, "cloud connector group ID", Getter(func(id string) error {
			_, _, err := zClient.cloudconnectorgroup.Get(id)
			return err
		}))
	case "CLIENT_TYPE":
		return customValidate(operand, []string{"id"}, "'zpn_client_type_zapp' or 'zpn_client_type_exporter' or 'zpn_client_type_ip_anchoring' or 'zpn_client_type_browser_isolation' or 'zpn_client_type_machine_tunnel' or 'zpn_client_type_edge_connector' or 'zpn_client_type_exporter_noauth' or 'zpn_client_type_slogger' or 'zpn_client_type_branch_connector'", Getter(func(id string) error {
			if id != "zpn_client_type_zapp" && id != "zpn_client_type_exporter" && id != "zpn_client_type_exporter_noauth" && id != "zpn_client_type_ip_anchoring" && id != "zpn_client_type_browser_isolation" && id != "zpn_client_type_machine_tunnel" && id != "zpn_client_type_edge_connector" && id != "zpn_client_type_slogger" && id != "zpn_client_type_branch_connector" && id != "zpn_client_type_zapp_partner" {
				return fmt.Errorf("RHS values must be 'zpn_client_type_zapp' or 'zpn_client_type_exporter' or 'zpn_client_type_exporter_noauth' or 'zpn_client_type_ip_anchoring' or 'zpn_client_type_browser_isolation' or 'zpn_client_type_machine_tunnel' or 'zpn_client_type_edge_connector' or 'zpn_client_type_slogger' or 'zpn_client_type_branch_connector' or 'zpn_client_type_zapp_partner 'when object type is CLIENT_TYPE")
			}
			return nil
		}))
	case "MACHINE_GRP":
		return customValidate(operand, []string{"id"}, "machine group ID", Getter(func(id string) error {
			_, _, err := zClient.machinegroup.WithMicroTenant(microtenantID).Get(id)
			return err
		}))
	case "POSTURE":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid posture network ID", operand.LHS, nil)
		}
		_, _, err := zClient.postureprofile.GetByPostureUDID(operand.LHS)
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid posture network ID", operand.LHS, err)
		}
		if !contains([]string{"true", "false"}, operand.RHS) {
			return rhsWarn(operand.ObjectType, "\"true\"/\"false\"", operand.RHS, nil)
		}
		return nil
	case "TRUSTED_NETWORK":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid trusted network ID", operand.LHS, nil)
		}
		_, _, err := zClient.trustednetwork.GetByNetID(operand.LHS)
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid trusted network ID", operand.LHS, err)
		}
		if operand.RHS != "true" {
			return rhsWarn(operand.ObjectType, "\"true\"", operand.RHS, nil)
		}
		return nil
	case "PLATFORM":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid platform ID", operand.LHS, nil)
		}
		_, _, err := zClient.platforms.GetAllPlatforms()
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid platform ID", operand.LHS, err)
		}
		if operand.RHS != "true" {
			return rhsWarn(operand.ObjectType, "\"true\"", operand.RHS, nil)
		}
		return nil
	case "SAML":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid SAML Attribute ID", operand.LHS, nil)
		}
		_, _, err := zClient.samlattribute.Get(operand.LHS)
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid SAML Attribute ID", operand.LHS, err)
		}
		if operand.RHS == "" {
			return rhsWarn(operand.ObjectType, "SAML Attribute Value", operand.RHS, nil)
		}
		return nil
	case "SCIM":
		if operand.IdpID == "" {
			return fmt.Errorf("[WARN] when operand object type is %v Idp ID must be set", operand.ObjectType)
		}
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid SCIM Attribute ID", operand.LHS, nil)
		}
		scim, _, err := zClient.scimattributeheader.Get(operand.IdpID, operand.LHS)
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid SCIM Attribute ID", operand.LHS, err)
		}
		if operand.RHS == "" {
			return rhsWarn(operand.ObjectType, "SCIM Attribute Value", operand.RHS, nil)
		}
		values, _ := zClient.scimattributeheader.SearchValues(scim.IdpID, scim.ID, operand.RHS)
		if len(values) == 0 {
			return rhsWarn(operand.ObjectType, fmt.Sprintf("valid SCIM Attribute Value (%s)", values), operand.RHS, nil)

		}
		return nil
	case "SCIM_GROUP":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid IDP Controller ID", operand.LHS, nil)
		}
		_, _, err := zClient.idpcontroller.Get(operand.LHS)
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid IDP Controller ID", operand.LHS, err)
		}
		if operand.RHS == "" {
			return rhsWarn(operand.ObjectType, "SCIM Group ID", operand.RHS, nil)
		}
		_, _, err = zClient.scimgroup.Get(operand.RHS)
		if err != nil {
			return rhsWarn(operand.ObjectType, "SCIM Group ID", operand.RHS, err)
		}
		return nil
	case "COUNTRY_CODE":
		if operand.LHS == "" || !isValidAlpha2(operand.LHS) {
			return lhsWarn(operand.ObjectType, "valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", operand.LHS, nil)
		}
		return nil
	default:
		return fmt.Errorf("[WARN] invalid operand object type %s", operand.ObjectType)
	}
}

type Getter func(id string) error

func (g Getter) Get(id string) error {
	return g(id)
}
func customValidate(operand policysetcontroller.Operands, expectedLHS []string, expectedRHS string, clientRHS Getter) error {
	if operand.LHS == "" || !contains(expectedLHS, operand.LHS) {
		return lhsWarn(operand.ObjectType, expectedLHS, operand.LHS, nil)
	}
	if operand.RHS == "" {
		return rhsWarn(operand.ObjectType, expectedRHS, operand.RHS, nil)
	}
	err := clientRHS.Get(operand.RHS)
	if err != nil {
		return rhsWarn(operand.ObjectType, expectedRHS, operand.RHS, err)
	}
	return nil
}
func rhsWarn(objType, expected, rhs interface{}, err error) error {
	return fmt.Errorf("[WARN] when operand object type is %v RHS must be %#v, value is \"%v\", %v", objType, expected, rhs, err)
}
func lhsWarn(objType, expected, lhs interface{}, err error) error {
	return fmt.Errorf("[WARN] when operand object type is %v LHS must be %#v value is \"%v\", %v", objType, expected, lhs, err)
}

func GetPolicyConditionsSchema(objectTypes []string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "This is for proviidng the set of conditions for the policy.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"negated": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"operator": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"AND",
						"OR",
					}, false),
				},
				"microtenant_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"operands": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "This signifies the various policy criteria.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"idp_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
							"name": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
							"lhs": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "This signifies the key for the object type. String ID example: id ",
							},
							"rhs": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								Description: "This denotes the value for the given object type. Its value depends upon the key.",
							},
							"microtenant_id": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								Description: "This denotes the value for the given object type. Its value depends upon the key.",
							},
							"rhs_list": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "This denotes a list of values for the given object type. The value depend upon the key. If rhs is defined this list will be ignored",
								Computed:    true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"object_type": {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "  This is for specifying the policy critiera.",
								ValidateFunc: validation.StringInSlice(objectTypes, false),
							},
						},
					},
				},
			},
		},
	}
}

func ExpandPolicyConditions(d *schema.ResourceData) ([]policysetcontroller.Conditions, error) {
	conditionInterface, ok := d.GetOk("conditions")
	if ok {
		conditions := conditionInterface.([]interface{})
		log.Printf("[INFO] conditions data: %+v\n", conditions)
		var conditionSets []policysetcontroller.Conditions
		for _, condition := range conditions {
			conditionSet, _ := condition.(map[string]interface{})
			if conditionSet != nil {
				operands, err := expandOperandsList(conditionSet["operands"])
				if err != nil {
					return nil, err
				}
				conditionSets = append(conditionSets, policysetcontroller.Conditions{
					ID:            conditionSet["id"].(string),
					Negated:       conditionSet["negated"].(bool),
					Operator:      conditionSet["operator"].(string),
					MicroTenantID: conditionSet["microtenant_id"].(string),
					Operands:      operands,
				})
			}
		}
		return conditionSets, nil
	}

	return []policysetcontroller.Conditions{}, nil
}

func expandOperandsList(ops interface{}) ([]policysetcontroller.Operands, error) {
	if ops != nil {
		operands := ops.([]interface{})
		log.Printf("[INFO] operands data: %+v\n", operands)
		var operandsSets []policysetcontroller.Operands
		for _, operand := range operands {
			operandSet, _ := operand.(map[string]interface{})
			id, _ := operandSet["id"].(string)
			IdpID, _ := operandSet["idp_id"].(string)
			rhs, ok := operandSet["rhs"].(string)
			op := policysetcontroller.Operands{
				ID:         id,
				Name:       operandSet["name"].(string),
				LHS:        operandSet["lhs"].(string),
				ObjectType: operandSet["object_type"].(string),
				IdpID:      IdpID,
				RHS:        rhs,
			}
			if ok && rhs != "" {
				if operandSet != nil {
					operandsSets = append(operandsSets, op)
				}
			} else {
				// try rhs_list
				rhsList := SetToStringSlice(operandSet["rhs_list"].(*schema.Set))
				if ok && len(rhsList) > 0 {
					for _, e := range rhsList {
						op.RHS = e
						operandsSets = append(operandsSets, op)
					}
				} else {
					log.Printf("[ERROR] No RHS is provided\n")
					return nil, fmt.Errorf("no RHS is provided")
				}

			}
		}

		return operandsSets, nil
	}
	return []policysetcontroller.Operands{}, nil
}

func flattenPolicyConditions(conditions []policysetcontroller.Conditions) []interface{} {
	ruleConditions := make([]interface{}, len(conditions))
	for i, ruleConditionItems := range conditions {
		ruleConditions[i] = map[string]interface{}{
			"id":             ruleConditionItems.ID,
			"negated":        ruleConditionItems.Negated,
			"operator":       ruleConditionItems.Operator,
			"microtenant_id": ruleConditionItems.MicroTenantID,
			"operands":       flattenPolicyRuleOperands(ruleConditionItems.Operands),
		}
	}

	return ruleConditions
}

func flattenPolicyRuleOperands(conditionOperand []policysetcontroller.Operands) []interface{} {
	conditionOperands := make([]interface{}, len(conditionOperand))
	for i, operandItems := range conditionOperand {
		conditionOperands[i] = map[string]interface{}{
			"id":             operandItems.ID,
			"idp_id":         operandItems.IdpID,
			"lhs":            operandItems.LHS,
			"object_type":    operandItems.ObjectType,
			"rhs":            operandItems.RHS,
			"name":           operandItems.Name,
			"microtenant_id": operandItems.MicroTenantID,
		}
	}

	return conditionOperands
}

func CommonPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "This field defines the description of the server.",
		},
		"bypass_default_rule": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"custom_msg": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "This is for providing a customer message for the user.",
		},
		"default_rule": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This is for providing a customer message for the user.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "This is the description of the access policy.",
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "This is the name of the policy.",
		},
		"operator": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				"AND",
				"OR",
			}, false),
		},
		"policy_set_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"policy_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"priority": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"reauth_default_rule": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"reauth_idle_timeout": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"reauth_timeout": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"zpn_isolation_profile_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"zpn_cbi_profile_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"zpn_inspection_profile_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"rule_order": {
			Type:       schema.TypeString,
			Optional:   true,
			Computed:   true,
			Deprecated: "The `rule_order` field is now deprecated for all zpa access policy resources in favor of the resource `zpa_policy_access_rule_reorder`",
		},
		"microtenant_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"lss_default_rule": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func resourceNetworkPortsSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"to": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func flattenNetworkPorts(ports []common.NetworkPorts) []interface{} {
	portsObj := make([]interface{}, len(ports))
	for i, val := range ports {
		portsObj[i] = map[string]interface{}{
			"from": val.From,
			"to":   val.To,
		}
	}
	return portsObj
}

/*
func expandNetwokPorts(d *schema.ResourceData, key string) []common.NetworkPorts {
	var ports []common.NetworkPorts
	if portsInterface, ok := d.GetOk(key); ok {
		portSet, ok := portsInterface.(*schema.Set)
		if !ok {
			log.Printf("[ERROR] conversion failed, destUdpPortsInterface")
			return ports
		}
		ports = make([]common.NetworkPorts, len(portSet.List()))
		for i, val := range portSet.List() {
			portItem := val.(map[string]interface{})
			ports[i] = common.NetworkPorts{
				From: portItem["from"].(string),
				To:   portItem["to"].(string),
			}
		}
	}
	return ports
}
*/

func resourceAppSegmentPortRange(desc string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		// Activate the "Attributes as Blocks" processing mode to permit dynamic declaration of no ports
		ConfigMode:  schema.SchemaConfigModeAttr,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.NoZeroValues,
				},
				"to": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.NoZeroValues,
				},
			},
		},
	}
}

func importPolicyStateContextFunc(types []string) schema.StateContextFunc {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		service := m.(*Client).policysetcontroller.WithMicroTenant(GetString(d.Get("microtenant_id")))

		id := d.Id()
		_, parseIDErr := strconv.ParseInt(id, 10, 64)
		if parseIDErr == nil {
			// assume if the passed value is an int
			_ = d.Set("id", id)
		} else {
			resp, _, err := service.GetByNameAndTypes(types, id)
			if err == nil {
				d.SetId(resp.ID)
				_ = d.Set("id", resp.ID)
			} else {
				return []*schema.ResourceData{d}, err
			}

		}
		return []*schema.ResourceData{d}, nil
	}
}

func dataInspectionRulesSchema() *schema.Schema {
	return &schema.Schema{
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
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenInspectionRules(rule []common.Rules) []interface{} {
	rules := make([]interface{}, len(rule))
	for i, rule := range rule {
		rules[i] = map[string]interface{}{
			"conditions": flattenInspectionRulesConditions(rule),
			"names":      rule.Names,
			"type":       rule.Type,
		}
	}

	return rules
}

func flattenInspectionRulesConditions(condition common.Rules) []interface{} {
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

func GetString(v interface{}) string {
	if v == nil {
		return ""
	}
	str, ok := v.(string)
	if ok {
		return str
	}
	return fmt.Sprintf("%v", v)
}

var AllowedPolicyTypes = map[string]struct{}{
	"ACCESS_POLICY":                        {},
	"GLOBAL_POLICY":                        {},
	"TIMEOUT_POLICY":                       {},
	"REAUTH_POLICY":                        {},
	"CLIENT_FORWARDING_POLICY":             {},
	"BYPASS_POLICY":                        {},
	"ISOLATION_POLICY":                     {},
	"INSPECTION_POLICY":                    {},
	"CREDENTIAL_POLICY":                    {},
	"CAPABILITIES_POLICY":                  {},
	"CLIENTLESS_SESSION_PROTECTION_POLICY": {},
}

func GetGlobalPolicySetByPolicyType(policysetcontroller policysetcontroller.Service, policyType string) (*policysetcontroller.PolicySet, error) {
	// Check if the provided policy type is allowed
	if _, ok := AllowedPolicyTypes[policyType]; !ok {
		return nil, fmt.Errorf("invalid policy type: %s", policyType)
	}

	policySetsMutex.Lock()
	defer policySetsMutex.Unlock()

	if p, ok := policySets[policyType]; ok {
		return &p, nil
	}
	globalPolicySet, _, err := policysetcontroller.GetByPolicyType(policyType)
	if err != nil {
		return nil, err
	}
	policySets[policyType] = *globalPolicySet
	return globalPolicySet, nil
}
