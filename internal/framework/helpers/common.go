package helpers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	clientpkg "github.com/SecurityGeekIO/terraform-provider-zpa/v4/internal/framework/client"
	fwstringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwrschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	fwplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	fwstringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	fwvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector_group"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/machinegroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/platforms"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/postureprofile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimattributeheader"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
)

type Client struct {
	*clientpkg.Client
	mu               sync.RWMutex
	policySetIDCache map[string]string
}

func NewHelperClient(client *clientpkg.Client) *Client {
	return &Client{
		Client:           client,
		policySetIDCache: make(map[string]string),
	}
}

var (
	policySets      = map[string]policysetcontroller.PolicySet{}
	policySetsMutex sync.Mutex
)

func validateAndSetProfileNameID(ctx context.Context, d *schema.ResourceData, service *zscaler.Service) error {
	// Ensure override_version_profile is true, otherwise we skip this logic
	overrideVersionProfile := d.Get("override_version_profile").(bool)
	if !overrideVersionProfile {
		return nil // Skip processing if override_version_profile is false
	}

	// If version_profile_id is already set, return (user explicitly provided it)
	if vpid, ok := d.GetOk("version_profile_id"); ok && vpid.(string) != "" {
		return nil
	}

	// Retrieve version_profile_name
	versionProfileName, ok := d.GetOk("version_profile_name")
	if !ok || versionProfileName.(string) == "" {
		return fmt.Errorf("version_profile_name must be provided when override_version_profile is enabled")
	}

	// Lookup version_profile_id using version_profile_name
	resp, _, err := customerversionprofile.GetByName(ctx, service, versionProfileName.(string))
	if err != nil {
		return fmt.Errorf("failed to find version profile with name %s: %v", versionProfileName, err)
	}

	// Set version_profile_id based on the lookup result
	_ = d.Set("version_profile_id", resp.ID)
	log.Printf("[INFO] Automatically resolved version_profile_id: %s for version_profile_name: %s", resp.ID, versionProfileName)

	return nil
}

func ValidateConditions(ctx context.Context, conditions []policysetcontroller.Conditions, zClient *Client, microtenantID string) error {
	for _, condition := range conditions {
		if err := validateOperands(ctx, condition.Operands, zClient, microtenantID); err != nil {
			return err
		}
	}
	return nil
}

func validateOperands(ctx context.Context, operands []policysetcontroller.Operands, zClient *Client, microtenantID string) error {
	for _, operand := range operands {
		if err := validateOperand(ctx, operand, zClient, microtenantID); err != nil {
			return err
		}
	}
	return nil
}

func validateOperand(ctx context.Context, operand policysetcontroller.Operands, zClient *Client, microtenantID string) error {
	switch operand.ObjectType {
	case "APP":
		return customValidate(operand, []string{"id"}, "application segment ID", Getter(func(id string) error {
			_, _, err := applicationsegment.Get(ctx, zClient.Service.WithMicroTenant(microtenantID), id)
			return err
		}))
	case "APP_GROUP":
		return customValidate(operand, []string{"id"}, "Segment Group ID", Getter(func(id string) error {
			_, _, err := segmentgroup.Get(ctx, zClient.Service.WithMicroTenant(microtenantID), id)
			return err
		}))

	case "IDP":
		return customValidate(operand, []string{"id"}, "IDP ID", Getter(func(id string) error {
			_, _, err := idpcontroller.Get(ctx, zClient.Service, id)
			return err
		}))
	case "EDGE_CONNECTOR_GROUP":
		return customValidate(operand, []string{"id"}, "cloud connector group ID", Getter(func(id string) error {
			_, _, err := cloud_connector_group.Get(ctx, zClient.Service, id)
			return err
		}))
	case "CLIENT_TYPE":
		return customValidate(operand, []string{"id"}, "'zpn_client_type_zapp' or 'zpn_client_type_exporter' or 'zpn_client_type_exporter_noauth' or 'zpn_client_type_ip_anchoring' or 'zpn_client_type_browser_isolation' or 'zpn_client_type_machine_tunnel' or 'zpn_client_type_edge_connector' or 'zpn_client_type_slogger' or 'zpn_client_type_branch_connector' or 'zpn_client_type_zapp_partner' or 'zpn_client_type_vdi'", Getter(func(id string) error {
			if id != "zpn_client_type_zapp" &&
				id != "zpn_client_type_exporter" &&
				id != "zpn_client_type_exporter_noauth" &&
				id != "zpn_client_type_ip_anchoring" &&
				id != "zpn_client_type_browser_isolation" &&
				id != "zpn_client_type_machine_tunnel" &&
				id != "zpn_client_type_edge_connector" &&
				id != "zpn_client_type_slogger" &&
				id != "zpn_client_type_branch_connector" &&
				id != "zpn_client_type_zapp_partner" &&
				id != "zpn_client_type_vdi" {
				return fmt.Errorf("RHS values must be one of 'zpn_client_type_zapp', 'zpn_client_type_exporter', 'zpn_client_type_exporter_noauth', 'zpn_client_type_ip_anchoring', 'zpn_client_type_browser_isolation', 'zpn_client_type_machine_tunnel', 'zpn_client_type_edge_connector', 'zpn_client_type_slogger', 'zpn_client_type_branch_connector', 'zpn_client_type_zapp_partner', or 'zpn_client_type_vdi' when object type is CLIENT_TYPE")
			}
			return nil
		}))

	case "MACHINE_GRP":
		return customValidate(operand, []string{"id"}, "machine group ID", Getter(func(id string) error {
			_, _, err := machinegroup.Get(ctx, zClient.Service.WithMicroTenant(microtenantID), id)
			return err
		}))
	case "POSTURE":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid posture network ID", operand.LHS, nil)
		}
		_, _, err := postureprofile.GetByPostureUDID(ctx, zClient.Service, operand.LHS)
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
		_, _, err := trustednetwork.GetByNetID(ctx, zClient.Service, operand.LHS)
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
		_, _, err := platforms.GetAllPlatforms(ctx, zClient.Service)
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
		_, _, err := samlattribute.Get(ctx, zClient.Service, operand.LHS)
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
		scim, _, err := scimattributeheader.Get(ctx, zClient.Service, operand.IdpID, operand.LHS)
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid SCIM Attribute ID", operand.LHS, err)
		}
		if operand.RHS == "" {
			return rhsWarn(operand.ObjectType, "SCIM Attribute Value", operand.RHS, nil)
		}
		values, _ := scimattributeheader.SearchValues(ctx, zClient.Service, scim.IdpID, scim.ID, operand.RHS)
		if len(values) == 0 {
			return rhsWarn(operand.ObjectType, fmt.Sprintf("valid SCIM Attribute Value (%s)", values), operand.RHS, nil)
		}
		return nil
	case "SCIM_GROUP":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "valid IDP Controller ID", operand.LHS, nil)
		}
		_, _, err := idpcontroller.Get(ctx, zClient.Service, operand.LHS)
		if err != nil {
			return lhsWarn(operand.ObjectType, "valid IDP Controller ID", operand.LHS, err)
		}
		if operand.RHS == "" {
			return rhsWarn(operand.ObjectType, "SCIM Group ID", operand.RHS, nil)
		}

		// Call the Get function with ScimGroup as the service parameter
		_, _, err = scimgroup.Get(ctx, zClient.Service, operand.RHS)
		if err != nil {
			return rhsWarn(operand.ObjectType, "SCIM Group ID", operand.RHS, err)
		}
		return nil
	case "COUNTRY_CODE":
		if operand.LHS == "" || !isValidAlpha2(operand.LHS) {
			return lhsWarn(operand.ObjectType, "valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", operand.LHS, nil)
		}
		return nil
	case "RISK_FACTOR_TYPE":
		// Check if lhs is "ZIA"
		if operand.LHS != "ZIA" {
			return lhsWarn(operand.ObjectType, "\"ZIA\"", operand.LHS, nil)
		}
		// Validate rhs for RISK_FACTOR_TYPE
		validRHSValues := []string{"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"}
		if !contains(validRHSValues, operand.RHS) {
			return rhsWarn(operand.ObjectType, "\"UNKNOWN\", \"LOW\", \"MEDIUM\", \"HIGH\", \"CRITICAL\"", operand.RHS, nil)
		}
		return nil
	case "CHROME_ENTERPRISE":
		if operand.LHS == "" {
			return lhsWarn(operand.ObjectType, "managed", operand.LHS, nil)
		}
		if operand.LHS != "managed" {
			return lhsWarn(operand.ObjectType, "managed", operand.LHS, nil)
		}
		if operand.RHS == "" {
			return rhsWarn(operand.ObjectType, "true/false", operand.RHS, nil)
		}
		if operand.RHS != "true" && operand.RHS != "false" {
			return rhsWarn(operand.ObjectType, "true/false", operand.RHS, nil)
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
			Computed: true,
		},
		"custom_msg": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "This is for providing a customer message for the user.",
		},
		"default_rule": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
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
			Computed: true,
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

func resourceAppSegmentPortRange(desc string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
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
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		zClient := meta.(*Client)
		service := zClient.Service

		microTenantID := GetString(d.Get("microtenant_id"))
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		id := d.Id()
		_, parseIDErr := strconv.ParseInt(id, 10, 64)
		if parseIDErr == nil {
			// assume if the passed value is an int
			_ = d.Set("id", id)
		} else {
			resp, _, err := policysetcontroller.GetByNameAndTypes(ctx, service, types, id)
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

func importPolicyStateContextFuncV2(types []string) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		zClient := meta.(*Client)
		service := zClient.Service

		microTenantID := GetString(d.Get("microtenant_id"))
		if microTenantID != "" {
			service = service.WithMicroTenant(microTenantID)
		}

		id := d.Id()
		_, parseIDErr := strconv.ParseInt(id, 10, 64)
		if parseIDErr == nil {
			// assume if the passed value is an int
			_ = d.Set("id", id)
		} else {
			resp, _, err := policysetcontrollerv2.GetByNameAndTypes(ctx, service, types, id)
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
	"REDIRECTION_POLICY":                   {},
}

func GetGlobalPolicySetByPolicyType(ctx context.Context, zClient *Client, policyType string) (*policysetcontroller.PolicySet, error) {
	// Check if the provided policy type is allowed
	if _, ok := AllowedPolicyTypes[policyType]; !ok {
		return nil, fmt.Errorf("invalid policy type: %s", policyType)
	}

	policySetsMutex.Lock()
	defer policySetsMutex.Unlock()

	if p, ok := policySets[policyType]; ok {
		return &p, nil
	}

	service := zClient.Service
	globalPolicySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, policyType)
	if err != nil {
		return nil, err
	}
	policySets[policyType] = *globalPolicySet
	return globalPolicySet, nil
}

// PolicyConditionsV2Block returns a reusable nested block definition that matches the shape of
// the access policy V2 conditions used throughout the provider's Plugin Framework resources.
func PolicyConditionsV2Block(objectTypes []string) fwrschema.Block {
	return fwrschema.SetNestedBlock{
		NestedObject: fwrschema.NestedBlockObject{
			Attributes: map[string]fwrschema.Attribute{
				"id": fwrschema.StringAttribute{
					Computed: true,
					PlanModifiers: []fwplanmodifier.String{
						fwstringplanmodifier.UseStateForUnknown(),
					},
				},
				"operator": fwrschema.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []fwvalidator.String{
						fwstringvalidator.OneOf("AND", "OR"),
					},
				},
			},
			Blocks: map[string]fwrschema.Block{
				"operands": fwrschema.SetNestedBlock{
					NestedObject: fwrschema.NestedBlockObject{
						Attributes: map[string]fwrschema.Attribute{
							"values": fwrschema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
							"object_type": fwrschema.StringAttribute{
								Optional: true,
								Validators: []fwvalidator.String{
									fwstringvalidator.OneOf(objectTypes...),
								},
							},
						},
						Blocks: map[string]fwrschema.Block{
							"entry_values": fwrschema.SetNestedBlock{
								NestedObject: fwrschema.NestedBlockObject{
									Attributes: map[string]fwrschema.Attribute{
										"rhs": fwrschema.StringAttribute{Optional: true},
										"lhs": fwrschema.StringAttribute{Optional: true},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// ######################################################################################################################
// ################################ ZPA ACCESS POLICY V2 COMMON CONDITIONS FUNCTIONS ####################################
// ######################################################################################################################
func ExpandPolicyConditionsV2(d *schema.ResourceData) ([]policysetcontrollerv2.PolicyRuleResourceConditions, error) {
	conditionInterface, ok := d.GetOk("conditions")
	if ok {
		// Assert conditionInterface as a *schema.Set
		conditionsSet := conditionInterface.(*schema.Set)
		log.Printf("[INFO] conditions data: %+v\n", conditionsSet.List())

		var conditionSets []policysetcontrollerv2.PolicyRuleResourceConditions
		for _, condition := range conditionsSet.List() { // Use Set.List() to get []interface{}
			conditionSet, _ := condition.(map[string]interface{})
			if conditionSet != nil {
				operands, err := expandOperandsListV2(conditionSet["operands"])
				if err != nil {
					return nil, err
				}
				conditionSets = append(conditionSets, policysetcontrollerv2.PolicyRuleResourceConditions{
					ID:       conditionSet["id"].(string),
					Operator: conditionSet["operator"].(string),
					Operands: operands,
				})
			}
		}
		return conditionSets, nil
	}
	return []policysetcontrollerv2.PolicyRuleResourceConditions{}, nil
}

func expandOperandsListV2(ops interface{}) ([]policysetcontrollerv2.PolicyRuleResourceOperands, error) {
	if ops != nil {
		// Assert ops as a *schema.Set
		operandsSet := ops.(*schema.Set)
		log.Printf("[INFO] operands data: %+v\n", operandsSet.List())

		var operandsSets []policysetcontrollerv2.PolicyRuleResourceOperands
		for _, operand := range operandsSet.List() { // Use Set.List() to get []interface{}
			operandSet, _ := operand.(map[string]interface{})
			id, _ := operandSet["id"].(string)
			IdpID, _ := operandSet["idp_id"].(string)

			// Extracting Values from TypeSet
			var values []string
			if valuesInterface, valuesOk := operandSet["values"].(*schema.Set); valuesOk && valuesInterface != nil {
				for _, v := range valuesInterface.List() {
					if strVal, ok := v.(string); ok {
						values = append(values, strVal)
					}
				}
			}

			// Extracting EntryValues from TypeSet
			var entryValues []policysetcontrollerv2.OperandsResourceLHSRHSValue
			if entryValuesInterface, ok := operandSet["entry_values"].(*schema.Set); ok {
				for _, ev := range entryValuesInterface.List() {
					entryValueMap, _ := ev.(map[string]interface{})
					lhs, _ := entryValueMap["lhs"].(string)
					rhs, _ := entryValueMap["rhs"].(string)

					entryValues = append(entryValues, policysetcontrollerv2.OperandsResourceLHSRHSValue{
						LHS: lhs,
						RHS: rhs,
					})
				}
			}

			log.Printf("[DEBUG] Extracted values: %+v\n", values)
			log.Printf("[DEBUG] Extracted entryValues: %+v\n", entryValues)

			op := policysetcontrollerv2.PolicyRuleResourceOperands{
				ID:                id,
				ObjectType:        operandSet["object_type"].(string),
				IDPID:             IdpID,
				Values:            values,
				EntryValuesLHSRHS: entryValues,
			}

			operandsSets = append(operandsSets, op)
		}

		return operandsSets, nil
	}
	return []policysetcontrollerv2.PolicyRuleResourceOperands{}, nil
}

func flattenConditionsV2(conditions []policysetcontrollerv2.PolicyRuleResourceConditions) []interface{} {
	if conditions == nil {
		return nil
	}

	c := make([]interface{}, len(conditions)) // Simplified slice initialization
	for i, condition := range conditions {
		condMap := make(map[string]interface{})
		condMap["operator"] = condition.Operator
		condMap["operands"] = flattenOperandsV2(condition.Operands)
		c[i] = condMap
	}
	return c
}

func flattenOperandsV2(operands []policysetcontrollerv2.PolicyRuleResourceOperands) []interface{} {
	if operands == nil {
		return nil
	}

	o := make([]interface{}, len(operands))
	for i, operand := range operands {
		operandMap := make(map[string]interface{})
		operandMap["object_type"] = operand.ObjectType

		if len(operand.Values) > 0 {
			operandMap["values"] = operand.Values
		} else {
			operandMap["values"] = []interface{}{}
		}

		entryValues := make([]interface{}, len(operand.EntryValuesLHSRHS))
		for j, entryValue := range operand.EntryValuesLHSRHS {
			entryValues[j] = map[string]interface{}{
				"lhs": entryValue.LHS,
				"rhs": entryValue.RHS,
			}
		}

		if len(entryValues) > 0 {
			operandMap["entry_values"] = entryValues
		} else {
			operandMap["entry_values"] = []interface{}{}
		}

		o[i] = operandMap
	}
	return o
}

func validateObjectTypeUniqueness(d *schema.ResourceData) error {
	conditions, ok := d.GetOk("conditions")
	if !ok {
		return nil
	}

	conditionsSet := conditions.(*schema.Set).List()

	for _, condition := range conditionsSet {
		conditionMap := condition.(map[string]interface{})
		if operands, ok := conditionMap["operands"].(*schema.Set); ok {
			objectTypeSet := make(map[string]struct{})

			for _, operand := range operands.List() {
				operandMap := operand.(map[string]interface{})
				objectType := operandMap["object_type"].(string)

				// Check for duplicate object_type
				if _, exists := objectTypeSet[objectType]; exists {
					return fmt.Errorf("object_type '%s' can only be specified once in the operands block. Please aggregate all entry_values under the same object_type", objectType)
				}
				objectTypeSet[objectType] = struct{}{}
			}
		}
	}

	return nil
}

// ValidatePolicyRuleConditions ensures that the necessary values are provided for specific object types.
func ValidatePolicyRuleConditions(d *schema.ResourceData) error {
	conditions, ok := d.GetOk("conditions")
	if !ok {
		// If conditions are not provided, there's nothing to validate
		return nil
	}

	validClientTypes := []string{
		"zpn_client_type_exporter",
		"zpn_client_type_exporter_noauth",
		"zpn_client_type_machine_tunnel",
		"zpn_client_type_edge_connector",
		"zpn_client_type_zia_inspection",
		"zpn_client_type_vdi",
		"zpn_client_type_zapp",
		"zpn_client_type_slogger",
		"zpn_client_type_browser_isolation",
		"zpn_client_type_ip_anchoring",
		"zpn_client_type_zapp_partner",
		"zpn_client_type_branch_connector",
	}

	validPlatformTypes := []string{"mac", "linux", "ios", "windows", "android"}

	validRiskScore := []string{"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"}

	// Adjust to handle *schema.Set instead of []interface{}
	conditionsSet := conditions.(*schema.Set)
	for _, condition := range conditionsSet.List() {
		conditionMap := condition.(map[string]interface{})
		operandsSet, ok := conditionMap["operands"].(*schema.Set)
		if !ok {
			// No operands to validate
			continue
		}

		for _, operand := range operandsSet.List() {
			operandMap := operand.(map[string]interface{})
			objectType := operandMap["object_type"].(string)
			valuesSet, valuesPresent := operandMap["values"].(*schema.Set)

			switch objectType {
			case "APP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("an Application Segment ID must be provided when object_type = APP")
				}
			case "APP_GROUP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Segment Group ID must be provided when object_type = APP_GROUP")
				}
			case "MACHINE_GRP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Machine Group ID must be provided when object_type = MACHINE_GRP")
				}
			case "LOCATION":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Location ID must be provided when object_type = LOCATION")
				}
			case "EDGE_CONNECTOR_GROUP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Edge Connector Group ID must be provided when object_type = EDGE_CONNECTOR_GROUP")
				}
			case "BRANCH_CONNECTOR_GROUP":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Branch Connector Group ID must be provided when object_type = BRANCH_CONNECTOR_GROUP")
				}
			case "USER_PORTAL":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a User Portal ID must be provided when object_type = USER_PORTAL")
				}
			case "CLIENT_TYPE":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("please provide one of the valid Client Types: %v", validClientTypes)
				}
				for _, v := range valuesSet.List() {
					value := v.(string)
					if !contains(validClientTypes, value) {
						return fmt.Errorf("invalid Client Type '%s'. Please provide one of the valid Client Types: %v", value, validClientTypes)
					}
				}
			case "PLATFORM":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("please provide one of the valid platform types: %v", validPlatformTypes)
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					// rhs, rhsOk := evMap["rhs"].(string)
					if !lhsOk || !contains(validPlatformTypes, lhs) {
						return fmt.Errorf("please provide one of the valid platform types: %v", validPlatformTypes)
					}
					// if !rhsOk || rhs != "true" {
					// 	return fmt.Errorf("rhs value must be 'true' for PLATFORM object_type")
					// }
				}
			case "RISK_FACTOR_TYPE":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("please provide valid risk factor values: %v", validRiskScore)
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)

					// Ensure lhs is "ZIA"
					if !lhsOk || lhs != "ZIA" {
						return fmt.Errorf("LHS must be 'ZIA' for RISK_FACTOR_TYPE")
					}

					// Ensure rhs is one of the valid risk scores
					validRHSValues := []string{"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"}
					if !rhsOk || !contains(validRHSValues, rhs) {
						return fmt.Errorf("RHS must be one of 'UNKNOWN', 'LOW', 'MEDIUM', 'HIGH', 'CRITICAL' for RISK_FACTOR_TYPE")
					}
				}
			case "POSTURE":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("please provide a valid Posture UDID")
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)
					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid Posture UDID and cannot be empty for POSTURE object_type")
					}
					if !rhsOk || (rhs != "true" && rhs != "false") {
						return fmt.Errorf("rhs value must be 'true' or 'false' for POSTURE object_type")
					}
				}
			case "TRUSTED_NETWORK":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("please provide a valid Network ID")
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)

					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid Network ID and cannot be empty for TRUSTED_NETWORK object_type")
					}
					if !rhsOk || (rhs != "true" && rhs != "false") {
						return fmt.Errorf("rhs value must be 'true' or 'false' for TRUSTED_NETWORK object_type")
					}
				}
			case "COUNTRY_CODE":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("please provide a valid country code in 'entry_values'")
				}

				var invalidCodes []string
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)

					// Validate 'lhs' as a country code
					if lhsOk {
						_, errors := validateCountryCode(lhs, "lhs")
						if len(errors) > 0 {
							// Collect invalid country codes instead of returning immediately
							invalidCodes = append(invalidCodes, lhs)
						}
					} else {
						return fmt.Errorf("a valid ISO-3166 Alpha-2 country code must be provided in 'lhs'")
					}

					// Ensure 'rhs' is "true"
					if !rhsOk || rhs != "true" {
						return fmt.Errorf("rhs value must be 'true' for COUNTRY_CODE object_type")
					}
				}

				// If there are any invalid country codes, return an aggregated error message
				if len(invalidCodes) > 0 {
					return fmt.Errorf("'%s' is not a valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", strings.Join(invalidCodes, "', '"))
				}
			case "SAML":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("entry_values must be provided for SAML object_type")
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs := evMap["rhs"].(string) // Directly accessing the value, assuming zero value (empty string) is not valid

					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid SAML attribute ID and cannot be empty for SAML object_type")
					}
					if rhs == "" {
						return fmt.Errorf("RHS must be a valid string i.e email adddress, department, group etc and cannot be empty for SAML object_type")
					}
				}
			case "SCIM":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("entry_values must be provided for SCIM object_type")
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs := evMap["rhs"].(string) // Directly accessing the value, assuming zero value (empty string) is not valid

					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid IdP ID and cannot be empty for SCIM object_type")
					}
					if rhs == "" {
						return fmt.Errorf("RHS must be a valid string i.e email adddress, department, group etc and cannot be empty for SCIM object_type")
					}
				}
			case "SCIM_GROUP":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("entry_values must be provided for SCIM_GROUP object_type")
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs := evMap["rhs"].(string) // Directly accessing the value, assuming zero value (empty string) is not valid

					if !lhsOk || lhs == "" {
						return fmt.Errorf("LHS must be a valid IdP ID and cannot be empty for SCIM_GROUP object_type")
					}
					if rhs == "" {
						return fmt.Errorf("RHS must be a valid scim group ID and cannot be empty for SCIM_GROUP object_type")
					}
				}
			case "CHROME_POSTURE_PROFILE":
				if !valuesPresent || valuesSet.Len() == 0 {
					return fmt.Errorf("a Chrome Posture Profile ID must be provided when object_type = CHROME_POSTURE_PROFILE")
				}
			case "CHROME_ENTERPRISE":
				entryValuesSet, ok := operandMap["entry_values"].(*schema.Set)
				if !ok || entryValuesSet.Len() == 0 {
					return fmt.Errorf("entry_values must be provided for CHROME_ENTERPRISE object_type")
				}
				for _, ev := range entryValuesSet.List() {
					evMap := ev.(map[string]interface{})
					lhs, lhsOk := evMap["lhs"].(string)
					rhs, rhsOk := evMap["rhs"].(string)

					if !lhsOk || lhs != "managed" {
						return fmt.Errorf("LHS must be 'managed' for CHROME_ENTERPRISE object_type")
					}
					if !rhsOk || (rhs != "true" && rhs != "false") {
						return fmt.Errorf("rhs value must be 'true' or 'false' for CHROME_ENTERPRISE object_type")
					}
				}
			}
		}
	}
	return nil
}

// FetchPolicySetIDByType returns the policy set ID for the supplied type, caching results per micro-tenant.
func FetchPolicySetIDByType(ctx context.Context, zClient *Client, policyType string, microTenantID string) (string, error) {
	// Create cache key including microtenant ID for multi-tenant scenarios
	cacheKey := policyType
	if microTenantID != "" {
		cacheKey = policyType + ":" + microTenantID
	}

	// First check: read lock (fast path for cached values)
	zClient.mu.RLock()
	if policySetID, found := zClient.policySetIDCache[cacheKey]; found {
		zClient.mu.RUnlock()
		return policySetID, nil
	}
	zClient.mu.RUnlock()

	// Not in cache, acquire write lock to prevent race condition
	zClient.mu.Lock()

	// Double-check: another goroutine might have fetched while we waited for the lock
	if policySetID, found := zClient.policySetIDCache[cacheKey]; found {
		zClient.mu.Unlock()
		return policySetID, nil
	}

	// Still not in cache, we need to fetch from API (hold the lock during fetch)
	service := zClient.Service
	if microTenantID != "" {
		service = service.WithMicroTenant(microTenantID)
	}

	globalPolicySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, policyType)
	if err != nil {
		zClient.mu.Unlock()
		return "", fmt.Errorf("failed to fetch policy set ID for type '%s': %v", policyType, err)
	}

	if zClient.policySetIDCache == nil {
		zClient.policySetIDCache = make(map[string]string)
	}

	// Store in cache before releasing lock
	zClient.policySetIDCache[cacheKey] = globalPolicySet.ID
	result := globalPolicySet.ID
	zClient.mu.Unlock()

	return result, nil
}

// ConvertV1ResponseToV2Request converts a PolicyRuleResource (API v1 response) to a PolicyRule (API v2 request) with aggregated values.
func ConvertV1ResponseToV2Request(v1Response policysetcontrollerv2.PolicyRuleResource) policysetcontrollerv2.PolicyRule {
	v2Request := policysetcontrollerv2.PolicyRule{
		ID:                           v1Response.ID,
		Name:                         v1Response.Name,
		Description:                  v1Response.Description,
		Action:                       v1Response.Action,
		PolicySetID:                  v1Response.PolicySetID,
		Operator:                     v1Response.Operator,
		CustomMsg:                    v1Response.CustomMsg,
		MicroTenantID:                v1Response.MicroTenantID,
		ZpnIsolationProfileID:        v1Response.ZpnIsolationProfileID,
		ZpnInspectionProfileID:       v1Response.ZpnInspectionProfileID,
		ActionID:                     v1Response.ActionID,
		Disabled:                     v1Response.Disabled,
		ExtranetEnabled:              v1Response.ExtranetEnabled,
		CreationTime:                 v1Response.CreationTime,
		ModifiedBy:                   v1Response.ModifiedBy,
		ModifiedTime:                 v1Response.ModifiedTime,
		PolicyType:                   v1Response.PolicyType,
		Priority:                     v1Response.Priority,
		ReauthIdleTimeout:            v1Response.ReauthIdleTimeout,
		ReauthTimeout:                v1Response.ReauthTimeout,
		RuleOrder:                    v1Response.RuleOrder,
		ZpnInspectionProfileName:     v1Response.ZpnInspectionProfileName,
		MicroTenantName:              v1Response.MicroTenantName,
		AppServerGroups:              v1Response.AppServerGroups,
		AppConnectorGroups:           v1Response.AppConnectorGroups,
		ServiceEdgeGroups:            v1Response.ServiceEdgeGroups,
		Credential:                   v1Response.Credential,
		CredentialPool:               v1Response.CredentialPool,
		PrivilegedCapabilities:       v1Response.PrivilegedCapabilities,
		PrivilegedPortalCapabilities: v1Response.PrivilegedPortalCapabilities,

		ExtranetDTO: v1Response.ExtranetDTO,

		Conditions: make([]policysetcontrollerv2.PolicyRuleResourceConditions, 0),
	}

	for _, condition := range v1Response.Conditions {
		newCondition := policysetcontrollerv2.PolicyRuleResourceConditions{
			Operator: condition.Operator,
			Operands: make([]policysetcontrollerv2.PolicyRuleResourceOperands, 0),
		}

		// Use a map to aggregate RHS values by ObjectType
		operandMap := make(map[string][]string)
		entryValuesMap := make(map[string][]policysetcontrollerv2.OperandsResourceLHSRHSValue)

		for _, operand := range condition.Operands {
			switch operand.ObjectType {
			case "APP", "APP_GROUP", "CONSOLE", "CHROME_POSTURE_PROFILE", "MACHINE_GRP", "LOCATION", "BRANCH_CONNECTOR_GROUP", "EDGE_CONNECTOR_GROUP", "CLIENT_TYPE", "USER_PORTAL", "PRIVILEGE_PORTAL":
				operandMap[operand.ObjectType] = append(operandMap[operand.ObjectType], operand.RHS)
			case "PLATFORM", "POSTURE", "TRUSTED_NETWORK", "SAML", "SCIM", "SCIM_GROUP", "COUNTRY_CODE", "RISK_FACTOR_TYPE", "CHROME_ENTERPRISE":
				entryValuesMap[operand.ObjectType] = append(entryValuesMap[operand.ObjectType], policysetcontrollerv2.OperandsResourceLHSRHSValue{
					LHS: operand.LHS,
					RHS: operand.RHS,
				})
			}
		}

		// Sort the operandMap by objectType to ensure consistent ordering
		sortedObjectTypes := make([]string, 0, len(operandMap))
		for objectType := range operandMap {
			sortedObjectTypes = append(sortedObjectTypes, objectType)
		}
		sort.Strings(sortedObjectTypes)

		// Create operand blocks from the aggregated data in sorted order
		for _, objectType := range sortedObjectTypes {
			values := operandMap[objectType]
			sort.Strings(values) // Sort the values within the objectType for consistency
			newCondition.Operands = append(newCondition.Operands, policysetcontrollerv2.PolicyRuleResourceOperands{
				ObjectType: objectType,
				Values:     values,
			})
		}

		// Sort the entryValuesMap by objectType to ensure consistent ordering
		sortedEntryValueTypes := make([]string, 0, len(entryValuesMap))
		for objectType := range entryValuesMap {
			sortedEntryValueTypes = append(sortedEntryValueTypes, objectType)
		}
		sort.Strings(sortedEntryValueTypes)

		// Create entryValues blocks from the aggregated data in sorted order
		for _, objectType := range sortedEntryValueTypes {
			entryValues := entryValuesMap[objectType]
			// Sort entryValues within the objectType by LHS, then RHS for consistency
			sort.Slice(entryValues, func(i, j int) bool {
				if entryValues[i].LHS == entryValues[j].LHS {
					return entryValues[i].RHS < entryValues[j].RHS
				}
				return entryValues[i].LHS < entryValues[j].LHS
			})
			newCondition.Operands = append(newCondition.Operands, policysetcontrollerv2.PolicyRuleResourceOperands{
				ObjectType:        objectType,
				EntryValuesLHSRHS: entryValues,
			})
		}

		// Append the sorted and normalized condition to the v2Request
		v2Request.Conditions = append(v2Request.Conditions, newCondition)
	}

	return v2Request
}

// -----------------------------------------------------------------------------
// Terraform Plugin Framework policy helpers
// -----------------------------------------------------------------------------

type PolicyConditionModel struct {
	ID            types.String                  `tfsdk:"id"`
	Operator      types.String                  `tfsdk:"operator"`
	MicroTenantID types.String                  `tfsdk:"microtenant_id"`
	CreationTime  types.String                  `tfsdk:"creation_time"`
	ModifiedBy    types.String                  `tfsdk:"modified_by"`
	ModifiedTime  types.String                  `tfsdk:"modified_time"`
	Operands      []PolicyConditionOperandModel `tfsdk:"operands"`
}

type PolicyConditionOperandModel struct {
	ID            types.String `tfsdk:"id"`
	IDPID         types.String `tfsdk:"idp_id"`
	Name          types.String `tfsdk:"name"`
	LHS           types.String `tfsdk:"lhs"`
	RHS           types.String `tfsdk:"rhs"`
	MicroTenantID types.String `tfsdk:"microtenant_id"`
	CreationTime  types.String `tfsdk:"creation_time"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	RHSList       types.Set    `tfsdk:"rhs_list"`
	ObjectType    types.String `tfsdk:"object_type"`
}

func PolicyCommonSchemaAttributes() map[string]fwrschema.Attribute {
	return map[string]fwrschema.Attribute{
		"id": fwrschema.StringAttribute{
			Computed: true,
			PlanModifiers: []fwplanmodifier.String{
				fwstringplanmodifier.UseStateForUnknown(),
			},
		},
		"name":                fwrschema.StringAttribute{Required: true},
		"description":         fwrschema.StringAttribute{Optional: true},
		"action_id":           fwrschema.StringAttribute{Optional: true},
		"bypass_default_rule": fwrschema.BoolAttribute{Optional: true, Computed: true},
		"custom_msg":          fwrschema.StringAttribute{Optional: true, Computed: true},
		"default_rule":        fwrschema.BoolAttribute{Optional: true, Computed: true},
		"operator": fwrschema.StringAttribute{
			Optional: true,
			Computed: true,
			Validators: []fwvalidator.String{
				fwstringvalidator.OneOf("AND", "OR"),
			},
		},
		"policy_set_id":       fwrschema.StringAttribute{Optional: true, Computed: true},
		"policy_type":         fwrschema.StringAttribute{Optional: true, Computed: true},
		"priority":            fwrschema.StringAttribute{Optional: true, Computed: true},
		"reauth_default_rule": fwrschema.BoolAttribute{Optional: true, Computed: true},
		"reauth_idle_timeout": fwrschema.StringAttribute{
			Optional: true,
			PlanModifiers: []fwplanmodifier.String{
				fwstringplanmodifier.RequiresReplace(),
			},
		},
		"reauth_timeout": fwrschema.StringAttribute{
			Optional: true,
			PlanModifiers: []fwplanmodifier.String{
				fwstringplanmodifier.RequiresReplace(),
			},
		},
		"zpn_isolation_profile_id":  fwrschema.StringAttribute{Optional: true, Computed: true},
		"zpn_cbi_profile_id":        fwrschema.StringAttribute{Optional: true, Computed: true},
		"zpn_inspection_profile_id": fwrschema.StringAttribute{Optional: true, Computed: true},
		"rule_order": fwrschema.StringAttribute{
			Optional:           true,
			Computed:           true,
			DeprecationMessage: "The `rule_order` field is deprecated in favor of the `zpa_policy_access_rule_reorder` resource.",
		},
		"microtenant_id":   fwrschema.StringAttribute{Optional: true, Computed: true},
		"lss_default_rule": fwrschema.BoolAttribute{Optional: true},
	}
}

func PolicyConditionsSchema(objectTypes []string) fwrschema.ListNestedAttribute {
	return fwrschema.ListNestedAttribute{
		Optional: true,
		NestedObject: fwrschema.NestedAttributeObject{
			Attributes: map[string]fwrschema.Attribute{
				"id":             fwrschema.StringAttribute{Computed: true},
				"operator":       fwrschema.StringAttribute{Required: true, Validators: []fwvalidator.String{fwstringvalidator.OneOf("AND", "OR")}},
				"microtenant_id": fwrschema.StringAttribute{Optional: true, Computed: true},
				"creation_time":  fwrschema.StringAttribute{Computed: true},
				"modified_by":    fwrschema.StringAttribute{Computed: true},
				"modified_time":  fwrschema.StringAttribute{Computed: true},
				"operands": fwrschema.ListNestedAttribute{
					Optional: true,
					NestedObject: fwrschema.NestedAttributeObject{
						Attributes: map[string]fwrschema.Attribute{
							"id":             fwrschema.StringAttribute{Computed: true},
							"idp_id":         fwrschema.StringAttribute{Optional: true, Computed: true},
							"name":           fwrschema.StringAttribute{Optional: true, Computed: true},
							"lhs":            fwrschema.StringAttribute{Required: true},
							"rhs":            fwrschema.StringAttribute{Optional: true, Computed: true},
							"microtenant_id": fwrschema.StringAttribute{Optional: true, Computed: true},
							"creation_time":  fwrschema.StringAttribute{Computed: true},
							"modified_by":    fwrschema.StringAttribute{Computed: true},
							"modified_time":  fwrschema.StringAttribute{Computed: true},
							"rhs_list": fwrschema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
							"object_type": fwrschema.StringAttribute{
								Required: true,
								Validators: []fwvalidator.String{
									fwstringvalidator.OneOf(objectTypes...),
								},
							},
						},
					},
				},
			},
		},
	}
}

func PolicyConditionsBlock(objectTypes []string) fwrschema.ListNestedBlock {
	return fwrschema.ListNestedBlock{
		NestedObject: fwrschema.NestedBlockObject{
			Attributes: map[string]fwrschema.Attribute{
				"id":             fwrschema.StringAttribute{Computed: true},
				"operator":       fwrschema.StringAttribute{Required: true, Validators: []fwvalidator.String{fwstringvalidator.OneOf("AND", "OR")}},
				"microtenant_id": fwrschema.StringAttribute{Optional: true, Computed: true},
				"creation_time":  fwrschema.StringAttribute{Computed: true},
				"modified_by":    fwrschema.StringAttribute{Computed: true},
				"modified_time":  fwrschema.StringAttribute{Computed: true},
			},
			Blocks: map[string]fwrschema.Block{
				"operands": fwrschema.ListNestedBlock{
					NestedObject: fwrschema.NestedBlockObject{
						Attributes: map[string]fwrschema.Attribute{
							"id":             fwrschema.StringAttribute{Computed: true},
							"idp_id":         fwrschema.StringAttribute{Optional: true, Computed: true},
							"name":           fwrschema.StringAttribute{Optional: true, Computed: true},
							"lhs":            fwrschema.StringAttribute{Required: true},
							"rhs":            fwrschema.StringAttribute{Optional: true, Computed: true},
							"microtenant_id": fwrschema.StringAttribute{Optional: true, Computed: true},
							"creation_time":  fwrschema.StringAttribute{Computed: true},
							"modified_by":    fwrschema.StringAttribute{Computed: true},
							"modified_time":  fwrschema.StringAttribute{Computed: true},
							"rhs_list": fwrschema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
							"object_type": fwrschema.StringAttribute{
								Required: true,
								Validators: []fwvalidator.String{
									fwstringvalidator.OneOf(objectTypes...),
								},
							},
						},
					},
				},
			},
		},
	}
}

func PolicyConditionModelsToSDK(ctx context.Context, models []PolicyConditionModel) ([]policysetcontroller.Conditions, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(models) == 0 {
		return nil, diags
	}

	conditions := make([]policysetcontroller.Conditions, 0, len(models))
	for _, model := range models {
		op, opDiags := policyOperandsFromModels(ctx, model.Operands)
		diags.Append(opDiags...)

		conditions = append(conditions, policysetcontroller.Conditions{
			ID:            StringValue(model.ID),
			Operator:      StringValue(model.Operator),
			MicroTenantID: StringValue(model.MicroTenantID),
			CreationTime:  StringValue(model.CreationTime),
			ModifiedBy:    StringValue(model.ModifiedBy),
			ModifiedTime:  StringValue(model.ModifiedTime),
			Operands:      op,
		})
	}

	return conditions, diags
}

func policyOperandsFromModels(ctx context.Context, models []PolicyConditionOperandModel) ([]policysetcontroller.Operands, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(models) == 0 {
		return nil, diags
	}

	operands := make([]policysetcontroller.Operands, 0)
	for _, model := range models {
		base := policysetcontroller.Operands{
			ID:            StringValue(model.ID),
			IdpID:         StringValue(model.IDPID),
			Name:          StringValue(model.Name),
			LHS:           StringValue(model.LHS),
			ObjectType:    StringValue(model.ObjectType),
			RHS:           StringValue(model.RHS),
			MicroTenantID: StringValue(model.MicroTenantID),
			CreationTime:  StringValue(model.CreationTime),
			ModifiedBy:    StringValue(model.ModifiedBy),
			ModifiedTime:  StringValue(model.ModifiedTime),
		}

		rhsList, rhsDiags := SetValueToStringSlice(ctx, model.RHSList)
		diags.Append(rhsDiags...)

		if base.RHS != "" || len(rhsList) == 0 {
			operands = append(operands, base)
			continue
		}

		for _, value := range rhsList {
			clone := base
			clone.RHS = strings.TrimSpace(value)
			operands = append(operands, clone)
		}
	}

	return operands, diags
}

func PolicyConditionsToModels(ctx context.Context, conditions []policysetcontroller.Conditions) ([]PolicyConditionModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(conditions) == 0 {
		return nil, diags
	}

	result := make([]PolicyConditionModel, 0, len(conditions))
	for _, condition := range conditions {
		operands, opDiags := policyOperandsToModels(condition.Operands)
		diags.Append(opDiags...)

		result = append(result, PolicyConditionModel{
			ID:            StringValueOrNull(condition.ID),
			Operator:      StringValueOrNull(condition.Operator),
			MicroTenantID: StringValueOrNull(condition.MicroTenantID),
			CreationTime:  StringValueOrNull(condition.CreationTime),
			ModifiedBy:    StringValueOrNull(condition.ModifiedBy),
			ModifiedTime:  StringValueOrNull(condition.ModifiedTime),
			Operands:      operands,
		})
	}

	return result, diags
}

func policyOperandsToModels(operands []policysetcontroller.Operands) ([]PolicyConditionOperandModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(operands) == 0 {
		return nil, diags
	}

	result := make([]PolicyConditionOperandModel, 0, len(operands))
	for _, operand := range operands {
		result = append(result, PolicyConditionOperandModel{
			ID:            StringValueOrNull(operand.ID),
			IDPID:         StringValueOrNull(operand.IdpID),
			Name:          StringValueOrNull(operand.Name),
			LHS:           StringValueOrNull(operand.LHS),
			RHS:           StringValueOrNull(operand.RHS),
			MicroTenantID: StringValueOrNull(operand.MicroTenantID),
			CreationTime:  StringValueOrNull(operand.CreationTime),
			ModifiedBy:    StringValueOrNull(operand.ModifiedBy),
			ModifiedTime:  StringValueOrNull(operand.ModifiedTime),
			ObjectType:    StringValueOrNull(operand.ObjectType),
			RHSList:       types.SetNull(types.StringType),
		})
	}

	return result, diags
}

func FlattenAppConnectorGroups(ctx context.Context, groups []appconnectorgroup.AppConnectorGroup) (types.List, diag.Diagnostics) {
	if len(groups) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: appConnectorGroupAttrTypesValue}), diag.Diagnostics{}
	}

	ids := make([]string, 0, len(groups))
	for _, group := range groups {
		if strings.TrimSpace(group.ID) != "" {
			ids = append(ids, group.ID)
		}
	}

	idSet, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: appConnectorGroupAttrTypesValue}), diags
	}

	obj, objDiags := types.ObjectValue(appConnectorGroupAttrTypesValue, map[string]attr.Value{"id": idSet})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: appConnectorGroupAttrTypesValue}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func ExpandAppConnectorGroups(ctx context.Context, list types.List) ([]appconnectorgroup.AppConnectorGroup, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	type appConnectorGroupModel struct {
		ID types.Set `tfsdk:"id"`
	}

	var models []appConnectorGroupModel
	var diags diag.Diagnostics
	diags.Append(list.ElementsAs(ctx, &models, false)...) // preserve order
	if diags.HasError() {
		return nil, diags
	}

	result := make([]appconnectorgroup.AppConnectorGroup, 0)
	for _, model := range models {
		ids, idsDiags := SetValueToStringSlice(ctx, model.ID)
		diags.Append(idsDiags...)
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			result = append(result, appconnectorgroup.AppConnectorGroup{ID: id})
		}
	}

	return result, diags
}

// -----------------------------------------------------------------------------
// Terraform Plugin Framework service edge helpers
// -----------------------------------------------------------------------------

func FlattenPrivateBrokerVersionToList(ctx context.Context, version serviceedgecontroller.PrivateBrokerVersion) (types.List, diag.Diagnostics) {
	subModuleAttrTypes := privateBrokerSubModuleAttrTypes()
	attrTypes := privateBrokerAttrTypes(subModuleAttrTypes)

	if version.ID == "" && len(version.ZPNSubModuleUpgradeList) == 0 {
		var diags diag.Diagnostics
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	attrValues, diags := buildPrivateBrokerVersionValues(ctx, version, subModuleAttrTypes)

	objValue, objDiags := types.ObjectValue(attrTypes, attrValues)
	diags.Append(objDiags...)

	listValue, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue})
	diags.Append(listDiags...)

	return listValue, diags
}

func FlattenPrivateBrokerVersionToObject(ctx context.Context, version serviceedgecontroller.PrivateBrokerVersion) (attr.Value, diag.Diagnostics) {
	subModuleAttrTypes := privateBrokerSubModuleAttrTypes()
	attrTypes := privateBrokerAttrTypes(subModuleAttrTypes)

	if version.ID == "" && len(version.ZPNSubModuleUpgradeList) == 0 {
		var diags diag.Diagnostics
		return types.ObjectNull(attrTypes), diags
	}

	attrValues, diags := buildPrivateBrokerVersionValues(ctx, version, subModuleAttrTypes)

	objValue, objDiags := types.ObjectValue(attrTypes, attrValues)
	diags.Append(objDiags...)

	return objValue, diags
}

func buildPrivateBrokerVersionValues(ctx context.Context, version serviceedgecontroller.PrivateBrokerVersion, subModuleAttrTypes map[string]attr.Type) (map[string]attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	subModuleList := types.ListNull(types.ObjectType{AttrTypes: subModuleAttrTypes})
	if len(version.ZPNSubModuleUpgradeList) > 0 {
		values := make([]attr.Value, 0, len(version.ZPNSubModuleUpgradeList))
		for _, module := range version.ZPNSubModuleUpgradeList {
			obj, objDiags := types.ObjectValue(subModuleAttrTypes, map[string]attr.Value{
				"id":               types.StringValue(module.ID),
				"creation_time":    types.StringValue(module.CreationTime),
				"current_version":  types.StringValue(module.CurrentVersion),
				"entity_gid":       types.StringValue(module.EntityGid),
				"entity_type":      types.StringValue(module.EntityType),
				"expected_version": types.StringValue(module.ExpectedVersion),
				"modified_by":      types.StringValue(module.ModifiedBy),
				"modified_time":    types.StringValue(module.ModifiedTime),
				"previous_version": types.StringValue(module.PreviousVersion),
				"role":             types.StringValue(module.Role),
				"upgrade_status":   types.StringValue(module.UpgradeStatus),
				"upgrade_time":     types.StringValue(module.UpgradeTime),
			})
			diags.Append(objDiags...)
			values = append(values, obj)
		}

		listValue, listDiags := types.ListValue(types.ObjectType{AttrTypes: subModuleAttrTypes}, values)
		diags.Append(listDiags...)
		subModuleList = listValue
	}

	attrValues := map[string]attr.Value{
		"id":                     types.StringValue(version.ID),
		"application_start_time": types.StringValue(version.ApplicationStartTime),
		"broker_id":              types.StringValue(version.BrokerId),
		"creation_time":          types.StringValue(version.CreationTime),
		"current_version":        types.StringValue(version.CurrentVersion),
		"disable_auto_update":    types.BoolValue(version.DisableAutoUpdate),
		"last_connect_time":      types.StringValue(version.LastConnectTime),
		"last_disconnect_time":   types.StringValue(version.LastDisconnectTime),
		"last_upgraded_time":     types.StringValue(version.LastUpgradedTime),
		"lone_warrior":           types.BoolValue(version.LoneWarrior),
		"modified_by":            types.StringValue(version.ModifiedBy),
		"modified_time":          types.StringValue(version.ModifiedTime),
		"platform":               types.StringValue(version.Platform),
		"platform_detail":        types.StringValue(version.PlatformDetail),
		"previous_version":       types.StringValue(version.PreviousVersion),
		"service_edge_group_id":  types.StringValue(version.ServiceEdgeGroupID),
		"private_ip":             types.StringValue(version.PrivateIP),
		"public_ip":              types.StringValue(version.PublicIP),
		"restart_instructions":   types.StringValue(version.RestartInstructions),
		"restart_time_in_sec":    types.StringValue(version.RestartTimeInSec),
		"runtime_os":             types.StringValue(version.RuntimeOS),
		"sarge_version":          types.StringValue(version.SargeVersion),
		"system_start_time":      types.StringValue(version.SystemStartTime),
		"tunnel_id":              types.StringValue(version.TunnelId),
		"upgrade_attempt":        types.StringValue(version.UpgradeAttempt),
		"upgrade_status":         types.StringValue(version.UpgradeStatus),
		"upgrade_now_once":       types.BoolValue(version.UpgradeNowOnce),
		"zpn_sub_module_upgrade": subModuleList,
	}

	return attrValues, diags
}

func privateBrokerSubModuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"creation_time":    types.StringType,
		"current_version":  types.StringType,
		"entity_gid":       types.StringType,
		"entity_type":      types.StringType,
		"expected_version": types.StringType,
		"modified_by":      types.StringType,
		"modified_time":    types.StringType,
		"previous_version": types.StringType,
		"role":             types.StringType,
		"upgrade_status":   types.StringType,
		"upgrade_time":     types.StringType,
	}
}

func privateBrokerAttrTypes(subModuleAttrTypes map[string]attr.Type) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"application_start_time": types.StringType,
		"broker_id":              types.StringType,
		"creation_time":          types.StringType,
		"current_version":        types.StringType,
		"disable_auto_update":    types.BoolType,
		"last_connect_time":      types.StringType,
		"last_disconnect_time":   types.StringType,
		"last_upgraded_time":     types.StringType,
		"lone_warrior":           types.BoolType,
		"modified_by":            types.StringType,
		"modified_time":          types.StringType,
		"platform":               types.StringType,
		"platform_detail":        types.StringType,
		"previous_version":       types.StringType,
		"service_edge_group_id":  types.StringType,
		"private_ip":             types.StringType,
		"public_ip":              types.StringType,
		"restart_instructions":   types.StringType,
		"restart_time_in_sec":    types.StringType,
		"runtime_os":             types.StringType,
		"sarge_version":          types.StringType,
		"system_start_time":      types.StringType,
		"tunnel_id":              types.StringType,
		"upgrade_attempt":        types.StringType,
		"upgrade_status":         types.StringType,
		"upgrade_now_once":       types.BoolType,
		"zpn_sub_module_upgrade": types.ListType{ElemType: types.ObjectType{AttrTypes: subModuleAttrTypes}},
	}
}

var networkPortAttrTypes = map[string]attr.Type{
	"from": types.StringType,
	"to":   types.StringType,
}

var serverGroupAttrTypes = map[string]attr.Type{
	"id": types.SetType{ElemType: types.StringType},
}

var appConnectorGroupAttrTypesValue = map[string]attr.Type{
	"id": types.SetType{ElemType: types.StringType},
}

// -----------------------------------------------------------------------------
// Terraform Plugin Framework application segment helpers
// -----------------------------------------------------------------------------

func FlattenNetworkPorts(ctx context.Context, ports []common.NetworkPorts) (types.List, diag.Diagnostics) {
	if len(ports) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: networkPortAttrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(ports))
	var diags diag.Diagnostics
	for _, port := range ports {
		obj, objDiags := types.ObjectValue(networkPortAttrTypes, map[string]attr.Value{
			"from": types.StringValue(port.From),
			"to":   types.StringValue(port.To),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: networkPortAttrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}

func ExpandNetworkPorts(ctx context.Context, ports types.List) ([]common.NetworkPorts, diag.Diagnostics) {
	if ports.IsNull() || ports.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	type networkPortModel struct {
		From types.String `tfsdk:"from"`
		To   types.String `tfsdk:"to"`
	}

	var models []networkPortModel
	var diags diag.Diagnostics
	diags.Append(ports.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(models) == 0 {
		return nil, diags
	}

	result := make([]common.NetworkPorts, 0, len(models))
	for _, model := range models {
		from := strings.TrimSpace(model.From.ValueString())
		to := strings.TrimSpace(model.To.ValueString())
		if from == "" && to == "" {
			continue
		}
		result = append(result, common.NetworkPorts{From: from, To: to})
	}

	return result, diags
}

func StringSliceToList(ctx context.Context, values []string) (types.List, diag.Diagnostics) {
	if len(values) == 0 {
		return types.ListNull(types.StringType), diag.Diagnostics{}
	}
	return types.ListValueFrom(ctx, types.StringType, values)
}

func ListValueToStringSlice(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return []string{}, diag.Diagnostics{}
	}

	var result []string
	var diags diag.Diagnostics
	diags.Append(list.ElementsAs(ctx, &result, false)...) // keep order
	if diags.HasError() {
		return nil, diags
	}
	return result, diags
}

func FlattenServerGroups(ctx context.Context, groups []servergroup.ServerGroup) (types.List, diag.Diagnostics) {
	if len(groups) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: serverGroupAttrTypes}), diag.Diagnostics{}
	}

	ids := make([]string, 0, len(groups))
	for _, group := range groups {
		if strings.TrimSpace(group.ID) != "" {
			ids = append(ids, group.ID)
		}
	}

	idSet, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: serverGroupAttrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(serverGroupAttrTypes, map[string]attr.Value{
		"id": idSet,
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: serverGroupAttrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func ExpandServerGroups(ctx context.Context, list types.List) ([]servergroup.ServerGroup, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	type serverGroupModel struct {
		ID types.Set `tfsdk:"id"`
	}

	var models []serverGroupModel
	var diags diag.Diagnostics
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(models) == 0 {
		return nil, diags
	}

	results := make([]servergroup.ServerGroup, 0)
	for _, model := range models {
		if model.ID.IsNull() || model.ID.IsUnknown() {
			continue
		}

		var ids []string
		diags.Append(model.ID.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, id := range ids {
			trimmed := strings.TrimSpace(id)
			if trimmed != "" {
				results = append(results, servergroup.ServerGroup{ID: trimmed})
			}
		}
	}

	return results, diags
}

func ValidateAppPorts(selectConnectorCloseToApp bool, udpPorts []common.NetworkPorts, udpRanges []string) diag.Diagnostics {
	var diags diag.Diagnostics
	if selectConnectorCloseToApp {
		if len(udpPorts) > 0 || len(udpRanges) > 0 {
			diags.AddError("Invalid port configuration", "App Connector Closest to App supports only TCP applications")
		}
	}
	return diags
}

func SetValueToStringSlice(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	if set.IsNull() || set.IsUnknown() {
		return []string{}, diag.Diagnostics{}
	}

	var result []string
	var diags diag.Diagnostics
	diags.Append(set.ElementsAs(ctx, &result, false)...) // preserve order
	return result, diags
}

func StringSliceToIntSlice(values []string) ([]int, diag.Diagnostics) {
	ints := make([]int, 0, len(values))
	var diags diag.Diagnostics

	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			diags.AddError("Invalid application ID", "Application IDs must not be empty")
			continue
		}

		intValue, err := strconv.Atoi(trimmed)
		if err != nil {
			diags.AddError("Invalid application ID", fmt.Sprintf("Unable to convert application ID %q to integer: %v", trimmed, err))
			continue
		}
		ints = append(ints, intValue)
	}

	return ints, diags
}

func ExpandZPNERID(ctx context.Context, list types.List) (*common.ZPNERID, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	type zpnERModel struct {
		ID types.Set `tfsdk:"id"`
	}

	var models []zpnERModel
	var diags diag.Diagnostics
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	for _, model := range models {
		if model.ID.IsNull() || model.ID.IsUnknown() {
			continue
		}

		var ids []string
		diags.Append(model.ID.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, id := range ids {
			trimmed := strings.TrimSpace(id)
			if trimmed != "" {
				return &common.ZPNERID{ID: trimmed}, diags
			}
		}
	}

	return nil, diags
}

func FlattenZPNERID(ctx context.Context, value *common.ZPNERID) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id": types.SetType{ElemType: types.StringType},
	}

	if value == nil || strings.TrimSpace(value.ID) == "" {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	idSet, diags := types.SetValueFrom(ctx, types.StringType, []string{value.ID})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"id": idSet,
	})
	diags.Append(objDiags...)

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

func ExpandPRACommonApps(ctx context.Context, value types.List) ([]applicationsegmentpra.AppsConfig, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	type praCommonAppsModel struct {
		AppsConfig types.List `tfsdk:"apps_config"`
	}

	type praAppConfigModel struct {
		AppID               types.String `tfsdk:"app_id"`
		PRAAppID            types.String `tfsdk:"pra_app_id"`
		Name                types.String `tfsdk:"name"`
		Description         types.String `tfsdk:"description"`
		AppTypes            types.Set    `tfsdk:"app_types"`
		ApplicationPort     types.String `tfsdk:"application_port"`
		ApplicationProtocol types.String `tfsdk:"application_protocol"`
		ConnectionSecurity  types.String `tfsdk:"connection_security"`
		Domain              types.String `tfsdk:"domain"`
	}

	var containers []praCommonAppsModel
	var diags diag.Diagnostics
	diags.Append(value.ElementsAs(ctx, &containers, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(containers) == 0 {
		return nil, diags
	}

	var configs []applicationsegmentpra.AppsConfig
	for _, container := range containers {
		if container.AppsConfig.IsNull() || container.AppsConfig.IsUnknown() {
			continue
		}

		var items []praAppConfigModel
		diags.Append(container.AppsConfig.ElementsAs(ctx, &items, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, item := range items {
			appTypes, appTypeDiags := SetValueToStringSlice(ctx, item.AppTypes)
			diags.Append(appTypeDiags...)
			if diags.HasError() {
				return nil, diags
			}

			configs = append(configs, applicationsegmentpra.AppsConfig{
				AppID:               strings.TrimSpace(item.AppID.ValueString()),
				PRAAppID:            strings.TrimSpace(item.PRAAppID.ValueString()),
				Name:                strings.TrimSpace(item.Name.ValueString()),
				Description:         strings.TrimSpace(item.Description.ValueString()),
				AppTypes:            appTypes,
				ApplicationPort:     strings.TrimSpace(item.ApplicationPort.ValueString()),
				ApplicationProtocol: strings.TrimSpace(item.ApplicationProtocol.ValueString()),
				ConnectionSecurity:  strings.TrimSpace(item.ConnectionSecurity.ValueString()),
				Domain:              strings.TrimSpace(item.Domain.ValueString()),
			})
		}
	}

	return configs, diags
}

func FlattenPRACommonApps(ctx context.Context, dto applicationsegmentpra.CommonAppsDto) (types.List, diag.Diagnostics) {
	appAttrTypes := map[string]attr.Type{
		"app_id":               types.StringType,
		"pra_app_id":           types.StringType,
		"name":                 types.StringType,
		"description":          types.StringType,
		"app_types":            types.SetType{ElemType: types.StringType},
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"connection_security":  types.StringType,
		"domain":               types.StringType,
	}

	commonAttrTypes := map[string]attr.Type{
		"apps_config": types.ListType{ElemType: types.ObjectType{AttrTypes: appAttrTypes}},
	}

	if len(dto.AppsConfig) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: commonAttrTypes}), diag.Diagnostics{}
	}

	appConfigValues := make([]attr.Value, 0, len(dto.AppsConfig))
	var diags diag.Diagnostics
	for _, app := range dto.AppsConfig {
		appTypes, setDiags := types.SetValueFrom(ctx, types.StringType, app.AppTypes)
		diags.Append(setDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: commonAttrTypes}), diags
		}

		obj, objDiags := types.ObjectValue(appAttrTypes, map[string]attr.Value{
			"app_id":               StringValueOrNull(app.AppID),
			"pra_app_id":           StringValueOrNull(app.PRAAppID),
			"name":                 StringValueOrNull(app.Name),
			"description":          StringValueOrNull(app.Description),
			"app_types":            appTypes,
			"application_port":     StringValueOrNull(app.ApplicationPort),
			"application_protocol": StringValueOrNull(app.ApplicationProtocol),
			"connection_security":  StringValueOrNull(app.ConnectionSecurity),
			"domain":               StringValueOrNull(app.Domain),
		})
		diags.Append(objDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: commonAttrTypes}), diags
		}
		appConfigValues = append(appConfigValues, obj)
	}

	appsConfigList, listDiags := types.ListValue(types.ObjectType{AttrTypes: appAttrTypes}, appConfigValues)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: commonAttrTypes}), diags
	}

	commonObj, commonDiags := types.ObjectValue(commonAttrTypes, map[string]attr.Value{
		"apps_config": appsConfigList,
	})
	diags.Append(commonDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: commonAttrTypes}), diags
	}

	listValue, listDiags := types.ListValue(types.ObjectType{AttrTypes: commonAttrTypes}, []attr.Value{commonObj})
	diags.Append(listDiags...)
	return listValue, diags
}

func FlattenPRAApps(ctx context.Context, apps []applicationsegmentpra.PRAApps) (types.List, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"id":                   types.StringType,
		"app_id":               types.StringType,
		"name":                 types.StringType,
		"description":          types.StringType,
		"application_port":     types.StringType,
		"application_protocol": types.StringType,
		"certificate_id":       types.StringType,
		"certificate_name":     types.StringType,
		"connection_security":  types.StringType,
		"domain":               types.StringType,
		"enabled":              types.BoolType,
		"hidden":               types.BoolType,
		"portal":               types.BoolType,
		"microtenant_id":       types.StringType,
		"microtenant_name":     types.StringType,
	}

	if len(apps) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diag.Diagnostics{}
	}

	values := make([]attr.Value, 0, len(apps))
	var diags diag.Diagnostics
	for _, app := range apps {
		obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"id":                   StringValueOrNull(app.ID),
			"app_id":               StringValueOrNull(app.AppID),
			"name":                 StringValueOrNull(app.Name),
			"description":          StringValueOrNull(app.Description),
			"application_port":     StringValueOrNull(app.ApplicationPort),
			"application_protocol": StringValueOrNull(app.ApplicationProtocol),
			"certificate_id":       StringValueOrNull(app.CertificateID),
			"certificate_name":     StringValueOrNull(app.CertificateName),
			"connection_security":  StringValueOrNull(app.ConnectionSecurity),
			"domain":               StringValueOrNull(app.Domain),
			"enabled":              types.BoolValue(app.Enabled),
			"hidden":               types.BoolValue(app.Hidden),
			"portal":               types.BoolValue(app.Portal),
			"microtenant_id":       StringValueOrNull(app.MicroTenantID),
			"microtenant_name":     StringValueOrNull(app.MicroTenantName),
		})
		diags.Append(objDiags...)
		values = append(values, obj)
	}

	list, listDiags := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, values)
	diags.Append(listDiags...)
	return list, diags
}
