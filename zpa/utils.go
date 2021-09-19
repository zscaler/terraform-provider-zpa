package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zpa/gozscaler/policysetrule"
)

func SetToStringSlice(d *schema.Set) []string {
	list := d.List()
	return ListToStringSlice(list)
}

func ListToStringSlice(v []interface{}) []string {
	if len(v) == 0 {
		return []string{}
	}

	ans := make([]string, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case nil:
			ans[i] = ""
		case string:
			ans[i] = x
		}
	}

	return ans
}

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
	case "CLOUD_CONNECTOR_GROUP":
		return customValidate(operand, []string{"id"}, "cloud connector group ID", Getter(func(id string) error {
			_, _, err := zClient.cloudconnectorgroup.Get(id)
			return err
		}))
	case "CLIENT_TYPE":
		return customValidate(operand, []string{"id"}, "'zpn_client_type_zapp' or 'zpn_client_type_exporter'", Getter(func(id string) error {
			if id != "zpn_client_type_zapp" && id != "zpn_client_type_exporter" {
				return fmt.Errorf("RHS values must be 'zpn_client_type_zapp' or 'zpn_client_type_exporter' wehn object type is CLIENT_TYPE")
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

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
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
	count, _, _ := zClient.policysetglobal.RulesCount()
	for k, v := range rules.orders {
		if v <= count {
			_, err := zClient.policysetrule.Reorder(policySetID, k, v)
			if err != nil {
				log.Printf("[ERROR] couldn't reorder the policy set, the order may not have taken place: %v\n", err)
			}
		}
	}
}
