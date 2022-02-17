package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/common"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/policysetrule"
)

func ValidateConditions(conditions []policysetrule.Conditions, zClient *Client) bool {
	for _, condition := range conditions {
		if !validateOperands(condition.Operands, zClient) {
			return false
		}
	}
	return true
}
func validateOperands(operands []policysetrule.Operands, zClient *Client) bool {
	for _, operand := range operands {
		if !validateOperand(operand, zClient) {
			return false
		}
	}
	return true
}
func validateOperand(operand policysetrule.Operands, zClient *Client) bool {
	switch operand.ObjectType {
	case "APP":
		return customValidate(operand, []string{"id"}, "application segment ID", Getter(func(id string) error {
			_, _, err := zClient.applicationsegment.Get(id)
			return err
		}))
	case "APP_GROUP":
		return customValidate(operand, []string{"id"}, "Segment Group ID", Getter(func(id string) error {
			_, _, err := zClient.segmentgroup.Get(id)
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
		return customValidate(operand, []string{"id"}, "'zpn_client_type_zapp' or 'zpn_client_type_exporter' or 'zpn_client_type_ip_anchoring' or 'zpn_client_type_browser_isolation' or 'zpn_client_type_machine_tunnel' or 'zpn_client_type_edge_connector'", Getter(func(id string) error {
			if id != "zpn_client_type_zapp" && id != "zpn_client_type_exporter" && id != "zpn_client_type_ip_anchoring" && id != "zpn_client_type_browser_isolation" && id != "zpn_client_type_machine_tunnel" && id != "zpn_client_type_edge_connector" {
				return fmt.Errorf("RHS values must be 'zpn_client_type_zapp' or 'zpn_client_type_exporter' or 'zpn_client_type_ip_anchoring' or 'zpn_client_type_browser_isolation' or 'zpn_client_type_machine_tunnel' or 'zpn_client_type_edge_connector' when object type is CLIENT_TYPE")
			}
			return nil
		}))
	case "MACHINE_GRP":
		return customValidate(operand, []string{"id"}, "machine group ID", Getter(func(id string) error {
			_, _, err := zClient.machinegroup.Get(id)
			return err
		}))
	case "POSTURE":
		if operand.LHS == "" {
			lhsWarn(operand.ObjectType, "valid posture network ID", operand.LHS, nil)
			return false
		}
		_, _, err := zClient.postureprofile.GetByPostureUDID(operand.LHS)
		if err != nil {
			lhsWarn(operand.ObjectType, "valid posture network ID", operand.LHS, err)
			return false
		}
		if !contains([]string{"true", "false"}, operand.RHS) {
			rhsWarn(operand.ObjectType, "\"true\"/\"false\"", operand.RHS, nil)
			return false
		}
		return true
	case "TRUSTED_NETWORK":
		if operand.LHS == "" {
			lhsWarn(operand.ObjectType, "valid trusted network ID", operand.LHS, nil)
			return false
		}
		_, _, err := zClient.trustednetwork.GetByNetID(operand.LHS)
		if err != nil {
			lhsWarn(operand.ObjectType, "valid trusted network ID", operand.LHS, err)
			return false
		}
		if operand.RHS != "true" {
			rhsWarn(operand.ObjectType, "\"true\"", operand.RHS, nil)
			return false
		}
		return true
	case "SAML":
		if operand.LHS == "" {
			lhsWarn(operand.ObjectType, "valid SAML Attribute ID", operand.LHS, nil)
			return false
		}
		_, _, err := zClient.samlattribute.Get(operand.LHS)
		if err != nil {
			lhsWarn(operand.ObjectType, "valid SAML Attribute ID", operand.LHS, err)
			return false
		}
		if operand.RHS == "" {
			rhsWarn(operand.ObjectType, "SAML Attribute Value", operand.RHS, nil)
			return false
		}
		return true
	case "SCIM":
		if operand.LHS == "" {
			lhsWarn(operand.ObjectType, "valid SCIM Attribute ID", operand.LHS, nil)
			return false
		}
		_, _, err := zClient.scimattributeheader.Get(operand.LHS)
		if err != nil {
			lhsWarn(operand.ObjectType, "valid SCIM Attribute ID", operand.LHS, err)
			return false
		}
		if operand.RHS == "" {
			rhsWarn(operand.ObjectType, "SCIM Attribute Value", operand.RHS, nil)
			return false
		}
		return true
	case "SCIM_GROUP":
		if operand.LHS == "" {
			lhsWarn(operand.ObjectType, "valid IDP Controller ID", operand.LHS, nil)
			return false
		}
		_, _, err := zClient.idpcontroller.Get(operand.LHS)
		if err != nil {
			lhsWarn(operand.ObjectType, "valid IDP Controller ID", operand.LHS, err)
			return false
		}
		if operand.RHS == "" {
			rhsWarn(operand.ObjectType, "SCIM Group ID", operand.RHS, nil)
			return false
		}
		_, _, err = zClient.scimgroup.Get(operand.RHS)
		if err != nil {
			rhsWarn(operand.ObjectType, "SCIM Group ID", operand.RHS, err)
			return false
		}
		return true
	default:
		log.Printf("[WARN] invalid operand object type %s\n", operand.ObjectType)
		return false
	}
}

type Getter func(id string) error

func (g Getter) Get(id string) error {
	return g(id)
}
func customValidate(operand policysetrule.Operands, expectedLHS []string, expectedRHS string, clientRHS Getter) bool {
	if operand.LHS == "" || !contains(expectedLHS, operand.LHS) {
		lhsWarn(operand.ObjectType, expectedLHS, operand.LHS, nil)
		return false
	}
	if operand.RHS == "" {
		rhsWarn(operand.ObjectType, expectedRHS, operand.RHS, nil)
		return false
	}
	err := clientRHS.Get(operand.RHS)
	if err != nil {
		rhsWarn(operand.ObjectType, expectedRHS, operand.RHS, err)
		return false
	}
	return true
}
func rhsWarn(objType, expected, rhs interface{}, err error) {
	log.Printf("[WARN] when operand object type is %v RHS must be %#v, value is \"%v\", %v\n", objType, expected, rhs, err)
}
func lhsWarn(objType, expected, lhs interface{}, err error) {
	log.Printf("[WARN] when operand object type is %v LHS must be %#v value is \"%v\", %v\n", objType, expected, lhs, err)
}

func reorder(orderI interface{}, policySetID, id string, zClient *Client) {
	defer reorderAll(policySetID, zClient)
	if orderI == nil {
		log.Printf("[WARN] Invalid order for policy set %s: %v\n", id, orderI)
		return
	}
	order, ok := orderI.(string)
	if !ok || order == "" {
		log.Printf("[WARN] Invalid order for policy set %s: %v\n", id, order)
		return
	}
	orderInt, err := strconv.Atoi(order)
	if err != nil || orderInt < 0 {
		log.Printf("[ERROR] couldn't reorder the policy set, the order may not have taken place:%v %v\n", orderInt, err)
		return
	}
	rules.Lock()
	rules.orders[id] = orderInt
	rules.Unlock()
}

// we keep calling reordering endpoint to reorder all rules after new rule was added
// because the reorder endpoint shifts all order up to replac the new order.
func reorderAll(policySetID string, zClient *Client) {
	rules.Lock()
	defer rules.Unlock()
	count, _, _ := zClient.policytype.RulesCount()
	for k, v := range rules.orders {
		if v <= count {
			_, err := zClient.policysetrule.Reorder(policySetID, k, v)
			if err != nil {
				log.Printf("[ERROR] couldn't reorder the policy set, the order may not have taken place: %v\n", err)
			}
		}
	}
}

func GetPolicyConditionsSchema(objectTypes []string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
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
				"operands": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
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
								Description: "This denotes the value for the given object type. Its value depends upon the key.",
							},
							"rhs_list": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "This denotes a list of values for the given object type. The value depend upon the key. If rhs is defined this list will be ignored",
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

func ExpandPolicyConditions(d *schema.ResourceData) ([]policysetrule.Conditions, error) {
	conditionInterface, ok := d.GetOk("conditions")
	if ok {
		conditions := conditionInterface.([]interface{})
		log.Printf("[INFO] conditions data: %+v\n", conditions)
		var conditionSets []policysetrule.Conditions
		for _, condition := range conditions {
			conditionSet, _ := condition.(map[string]interface{})
			if conditionSet != nil {
				operands, err := expandOperandsList(conditionSet["operands"])
				if err != nil {
					return nil, err
				}
				conditionSets = append(conditionSets, policysetrule.Conditions{
					ID:       conditionSet["id"].(string),
					Negated:  conditionSet["negated"].(bool),
					Operator: conditionSet["operator"].(string),
					Operands: operands,
				})
			}
		}
		return conditionSets, nil
	}

	return []policysetrule.Conditions{}, nil
}

func expandOperandsList(ops interface{}) ([]policysetrule.Operands, error) {
	if ops != nil {
		operands := ops.([]interface{})
		log.Printf("[INFO] operands data: %+v\n", operands)
		var operandsSets []policysetrule.Operands
		for _, operand := range operands {
			operandSet, _ := operand.(map[string]interface{})
			id, _ := operandSet["id"].(string)
			IdpID, _ := operandSet["idp_id"].(string)
			rhs, ok := operandSet["rhs"].(string)
			op := policysetrule.Operands{
				ID:         id,
				IdpID:      IdpID,
				LHS:        operandSet["lhs"].(string),
				ObjectType: operandSet["object_type"].(string),
				RHS:        rhs,
				Name:       operandSet["name"].(string),
			}
			if ok && rhs != "" {
				if operandSet != nil {
					operandsSets = append(operandsSets, op)
				}
			} else {
				// try rhs_list
				rhsList, ok := operandSet["rhs_list"].([]interface{})
				if ok && len(rhsList) > 0 {
					for _, e := range rhsList {
						op.RHS, _ = e.(string)
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
	return []policysetrule.Operands{}, nil
}

func flattenPolicyConditions(conditions []policysetrule.Conditions) []interface{} {
	ruleConditions := make([]interface{}, len(conditions))
	for i, ruleConditionItems := range conditions {
		ruleConditions[i] = map[string]interface{}{
			"id":       ruleConditionItems.ID,
			"negated":  ruleConditionItems.Negated,
			"operator": ruleConditionItems.Operator,
			"operands": flattenPolicyRuleOperands(ruleConditionItems.Operands),
		}
	}

	return ruleConditions
}

func flattenPolicyRuleOperands(conditionOperand []policysetrule.Operands) []interface{} {
	conditionOperands := make([]interface{}, len(conditionOperand))
	for i, operandItems := range conditionOperand {
		conditionOperands[i] = map[string]interface{}{
			"id":          operandItems.ID,
			"idp_id":      operandItems.IdpID,
			"lhs":         operandItems.LHS,
			"object_type": operandItems.ObjectType,
			"rhs":         operandItems.RHS,
			"name":        operandItems.Name,
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
			Required: true,
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
		},
		"reauth_timeout": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"rule_order": {
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
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
				"to": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
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

func resourceAppSegmentPortRange(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
				"to": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
			},
		},
	}
}

// func flattenAppSegmentPortRange(ports []applicationsegment.AppSegmentPortRange) []interface{} {
// 	portsObj := make([]interface{}, len(ports))
// 	for i, val := range ports {
// 		portsObj[i] = map[string]interface{}{
// 			"from": val.From,
// 			"to":   val.To,
// 		}
// 	}
// 	return portsObj
// }

// func expandAppSegmentPortRange(d *schema.ResourceData, key string) []applicationsegment.AppSegmentPortRange {
// 	var ports []applicationsegment.AppSegmentPortRange
// 	if portsInterface, ok := d.GetOk(key); ok {
// 		portSet, ok := portsInterface.(*schema.Set)
// 		if !ok {
// 			log.Printf("[ERROR] conversion failed, destUdpPortsInterface")
// 			return ports
// 		}
// 		ports = make([]applicationsegment.AppSegmentPortRange, len(portSet.List()))
// 		for i, val := range portSet.List() {
// 			portItem := val.(map[string]interface{})
// 			ports[i] = applicationsegment.AppSegmentPortRange{
// 				From: portItem["from"].(string),
// 				To:   portItem["to"].(string),
// 			}
// 		}
// 	}
// 	return ports
// }

func importPolicyStateFunc(types []string) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		zClient := m.(*Client)
		id := d.Id()
		_, parseIDErr := strconv.ParseInt(id, 10, 64)
		if parseIDErr == nil {
			// assume if the passed value is an int
			d.Set("id", id)
		} else {
			resp, _, err := zClient.policysetrule.GetByNameAndTypes(types, id)
			if err == nil {
				d.SetId(resp.ID)
				d.Set("id", resp.ID)
			} else {
				return []*schema.ResourceData{d}, err
			}

		}
		return []*schema.ResourceData{d}, nil
	}
}
