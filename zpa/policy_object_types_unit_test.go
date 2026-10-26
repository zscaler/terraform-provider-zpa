package zpa

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

// conditionObjectTypes returns the object types a policy rule resource accepts
// in conditions.operands.object_type, from its schema's validation.
func conditionObjectTypes(t *testing.T, r *schema.Resource) []string {
	t.Helper()
	operands := r.Schema["conditions"].Elem.(*schema.Resource).Schema["operands"].Elem.(*schema.Resource)
	objectType := operands.Schema["object_type"]
	var accepted []string
	for _, candidate := range []string{
		"APP", "APP_GROUP", "LOCATION", "IDP", "SAML", "SCIM", "SCIM_GROUP", "CLIENT_TYPE", "POSTURE", "TRUSTED_NETWORK",
		"BRANCH_CONNECTOR_GROUP", "EDGE_CONNECTOR_GROUP", "MACHINE_GRP", "COUNTRY_CODE", "PLATFORM", "RISK_FACTOR_TYPE",
		"CHROME_ENTERPRISE", "CHROME_POSTURE_PROFILE", "WORKLOAD_TAG_GROUP", "ENTITY_TAG_GRP", "ENTITY_TAG_VAL",
		"CONSOLE", "USER_PORTAL", "PRIVILEGE_PORTAL",
	} {
		var errs []error
		if objectType.ValidateFunc != nil {
			_, errs = objectType.ValidateFunc(candidate, "object_type")
		} else if objectType.ValidateDiagFunc != nil {
			if diags := objectType.ValidateDiagFunc(candidate, nil); diags.HasError() {
				errs = append(errs, diags[0].Validate())
			}
		} else {
			t.Fatalf("object_type has no validation")
		}
		if len(errs) == 0 {
			accepted = append(accepted, candidate)
		}
	}
	sort.Strings(accepted)
	return accepted
}

// A v1 policy rule resource and its v2 counterpart manage the same rules through
// different API shapes, so they accept the same condition object types.
func TestPolicyRules_V1AndV2AcceptTheSameObjectTypes(t *testing.T) {
	pairs := []struct {
		name   string
		v1, v2 func() *schema.Resource
	}{
		{"access", resourcePolicyAccessRule, resourcePolicyAccessRuleV2},
		{"forwarding", resourcePolicyForwardingRule, resourcePolicyForwardingRuleV2},
		{"inspection", resourcePolicyInspectionRule, resourcePolicyInspectionRuleV2},
		{"isolation", resourcePolicyIsolationRule, resourcePolicyIsolationRuleV2},
		{"timeout", resourcePolicyTimeoutRule, resourcePolicyTimeoutRuleV2},
	}
	for _, p := range pairs {
		v1, v2 := conditionObjectTypes(t, p.v1()), conditionObjectTypes(t, p.v2())
		if !reflect.DeepEqual(v1, v2) {
			t.Errorf("%s: v1 accepts %v\n   v2 accepts %v", p.name, v1, v2)
		}
	}
}

// Every object type a v1 resource accepts in its schema must also pass the
// operand validation its create and update run, or the apply fails.
func TestPolicyRules_V1ValidationKnowsEveryObjectType(t *testing.T) {
	for _, r := range []func() *schema.Resource{
		resourcePolicyAccessRule, resourcePolicyForwardingRule, resourcePolicyInspectionRule, resourcePolicyIsolationRule, resourcePolicyTimeoutRule,
	} {
		for _, objectType := range conditionObjectTypes(t, r()) {
			err := validateOperand(context.Background(), policysetcontroller.Operands{ObjectType: objectType}, nil, "")
			if err != nil && strings.Contains(err.Error(), "invalid operand object type") {
				t.Errorf("object type %s is accepted by the schema but rejected by validateOperand", objectType)
			}
		}
	}
}

// ID-based object types without a lookup by ID need lhs "id" and a non-empty rhs.
func TestValidateOperand_IDBasedTypesWithoutLookup(t *testing.T) {
	for _, objectType := range []string{"LOCATION", "BRANCH_CONNECTOR_GROUP", "WORKLOAD_TAG_GROUP", "ENTITY_TAG_GRP", "ENTITY_TAG_VAL", "CHROME_POSTURE_PROFILE"} {
		ok := policysetcontroller.Operands{ObjectType: objectType, LHS: "id", RHS: "216196257331282583"}
		if err := validateOperand(context.Background(), ok, nil, ""); err != nil {
			t.Errorf("%s with lhs id and an rhs: unexpected error %v", objectType, err)
		}
		badLHS := policysetcontroller.Operands{ObjectType: objectType, LHS: "name", RHS: "x"}
		if err := validateOperand(context.Background(), badLHS, nil, ""); err == nil || !strings.Contains(err.Error(), "LHS must be") {
			t.Errorf("%s with lhs name: want an LHS error, got %v", objectType, err)
		}
		noRHS := policysetcontroller.Operands{ObjectType: objectType, LHS: "id"}
		if err := validateOperand(context.Background(), noRHS, nil, ""); err == nil || !strings.Contains(err.Error(), "RHS must be") {
			t.Errorf("%s without rhs: want an RHS error, got %v", objectType, err)
		}
	}
}

// The v2 read converts the v1-shaped operands of the API response: every
// object type the v2 access rule accepts must come back, as values for
// ID-based types and as entry values for lhs/rhs types.
func TestConvertV1ResponseToV2Request_KeepsEveryObjectType(t *testing.T) {
	idBased := []string{"APP", "APP_GROUP", "LOCATION", "IDP", "CLIENT_TYPE", "BRANCH_CONNECTOR_GROUP", "EDGE_CONNECTOR_GROUP", "MACHINE_GRP", "CHROME_POSTURE_PROFILE", "WORKLOAD_TAG_GROUP", "ENTITY_TAG_GRP", "ENTITY_TAG_VAL"}
	lhsRHS := []string{"SAML", "SCIM", "SCIM_GROUP", "POSTURE", "TRUSTED_NETWORK", "COUNTRY_CODE", "PLATFORM", "RISK_FACTOR_TYPE", "CHROME_ENTERPRISE"}

	var operands []policysetcontrollerv2.PolicyRuleResourceOperands
	for _, objectType := range idBased {
		operands = append(operands, policysetcontrollerv2.PolicyRuleResourceOperands{ObjectType: objectType, LHS: "id", RHS: "1"})
	}
	for _, objectType := range lhsRHS {
		operands = append(operands, policysetcontrollerv2.PolicyRuleResourceOperands{ObjectType: objectType, LHS: "l", RHS: "r"})
	}
	got := ConvertV1ResponseToV2Request(policysetcontrollerv2.PolicyRuleResource{
		Conditions: []policysetcontrollerv2.PolicyRuleResourceConditions{{Operator: "OR", Operands: operands}},
	})

	values := map[string][]string{}
	entries := map[string]int{}
	for _, op := range got.Conditions[0].Operands {
		if len(op.Values) > 0 {
			values[op.ObjectType] = op.Values
		}
		if len(op.EntryValuesLHSRHS) > 0 {
			entries[op.ObjectType] = len(op.EntryValuesLHSRHS)
		}
	}
	for _, objectType := range idBased {
		if !reflect.DeepEqual(values[objectType], []string{"1"}) {
			t.Errorf("%s: want values [1], got %v (entry values: %d)", objectType, values[objectType], entries[objectType])
		}
	}
	for _, objectType := range lhsRHS {
		if entries[objectType] != 1 {
			t.Errorf("%s: want 1 entry value, got %d (values: %v)", objectType, entries[objectType], values[objectType])
		}
	}
}
