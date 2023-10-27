package zpa

import (
	"errors"
	"fmt"

	"github.com/fabiotavarespr/iso3166"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/lssconfigcontroller"
)

// This function checks for ISO3166 Alpha2 Country codes
// It's used within the access policy RHS validation
func isValidAlpha2(code string) bool {
	return iso3166.ExistsIso3166ByAlpha2Code(code)
}

func validateCountryCode(value interface{}, key string) ([]string, []error) {
	var warnings []string
	var errors []error

	code, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", key))
		return warnings, errors
	}

	if !iso3166.ExistsIso3166ByAlpha2Code(code) {
		errors = append(errors, fmt.Errorf("'%s' is not a valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", code))
	}

	return warnings, errors
}

var supportedLSSUserActivity = []string{
	"BRK_MT_SETUP_FAIL_BIND_TO_AST_LOCAL_OWNER",
	"CLT_INVALID_DOMAIN",
	"AST_MT_SETUP_ERR_HASH_TBL_FULL",
	"AST_MT_SETUP_ERR_CONN_PEER",
	"BRK_MT_SETUP_FAIL_REJECTED_BY_POLICY_APPROVAL",
	"BRK_MT_SETUP_FAIL_ICMP_RATE_LIMIT_NUM_APP_EXCEEDED",
	"EXPTR_MT_TLS_SETUP_FAIL_VERSION_MISMATCH",
	"BRK_MT_SETUP_FAIL_RATE_LIMIT_LOOP_DETECTED",
	"CLT_INVALID_TAG",
	"AST_MT_SETUP_ERR_NO_SYSTEM_FD",
	"AST_MT_SETUP_ERR_NO_PROCESS_FD",
	"BROKER_NOT_ENABLED",
	"AST_MT_SETUP_ERR_AST_CFG_DISABLED",
	"BRK_MT_SETUP_FAIL_TOO_MANY_FAILED_ATTEMPTS",
	"BRK_MT_AUTH_NO_SAML_ASSERTION_IN_MSG",
	"BRK_MT_SETUP_FAIL_CTRL_BRK_CANNOT_FIND_CONNECTOR",
	"INVALID_DOMAIN",
	"BRK_MT_TERMINATED_BRK_SWITCHED",
	"AST_MT_SETUP_ERR_OPEN_SERVER_CLOSE",
	"AST_MT_SETUP_ERR_BIND_TO_AST_LOCAL_OWNER",
	"NO_CONNECTOR_AVAILABLE",
	"BRK_MT_AUTH_SAML_CANNOT_ADD_ATTR_TO_HEAP",
	"EXPTR_MT_TLS_SETUP_FAIL_NOT_TRUSTED_CA",
	"AST_MT_SETUP_TIMEOUT_NO_ACK_TO_BIND",
	"CLT_PORT_UNREACHABLE",
	"C2C_CLIENT_CONN_EXPIRED",
	"BRK_MT_SETUP_FAIL_BIND_TO_CLIENT_LOCAL_OWNER",
	"BRK_MT_AUTH_SAML_CANNOT_ADD_ATTR_TO_HASH",
	"BRK_MT_SETUP_FAIL_REPEATED_DISPATCH",
	"AST_MT_SETUP_ERR_OPEN_SERVER_ERROR",
	"DSP_MT_SETUP_FAIL_DISCOVERY_TIMEOUT",
	"CUSTOMER_NOT_ENABLED",
	"BRK_CONN_UPGRADE_REQUEST_FAILED",
	"C2C_MTUNNEL_FAILED_FORWARD",
	"EXPTR_MT_TLS_SETUP_FAIL_CERT_CHAIN_ISSUE",
	"AST_MT_SETUP_ERR_RATE_LIMIT_REACHED",
	"BRK_MT_SETUP_FAIL_RATE_LIMIT_NUM_APP_EXCEEDED",
	"CLT_WRONG_PORT",
	"AST_MT_SETUP_TIMEOUT_CANNOT_CONN_TO_SERVER",
	"BRK_MT_AUTH_SAML_FINGER_PRINT_FAIL",
	"AST_MT_SETUP_ERR_NO_EPHEMERAL_PORT",
	"BRK_CONN_UPGRADE_REQUEST_FORBIDDEN",
	"AST_MT_SETUP_ERR_OPEN_SERVER_CONN",
	"CLT_PROBE_FAILED",
	"AST_MT_SETUP_ERR_APP_NOT_FOUND",
	"AST_MT_SETUP_ERR_OPEN_BROKER_CONN",
	"BRK_MT_SETUP_FAIL_ICMP_RATE_LIMIT_EXCEEDED",
	"AST_MT_SETUP_ERR_OPEN_SERVER_TIMEOUT",
	"C2C_MTUNNEL_BAD_STATE",
	"CLT_DUPLICATE_TAG",
	"AST_MT_SETUP_TIMEOUT",
	"CLT_DOUBLEENCRYPT_NOT_SUPPORTED",
	"BRK_MT_SETUP_FAIL_CANNOT_SEND_MT_COMPLETE",
	"BRK_MT_SETUP_FAIL_BIND_RECV_IN_BAD_STATE",
	"APP_NOT_AVAILABLE",
	"BRK_MT_AUTH_SAML_NO_USER_ID",
	"AST_MT_SETUP_TIMEOUT_CANNOT_CONN_TO_BROKER",
	"DSP_MT_SETUP_FAIL_MISSING_HEALTH",
	"AST_MT_SETUP_ERR_DUP_MT_ID",
	"AST_MT_SETUP_ERR_BIND_GLOBAL_OWNER",
	"BRK_MT_TERMINATED_APPROVAL_TIMEOUT",
	"AST_MT_SETUP_ERR_BIND_ACK",
	"CLT_CONN_FAILED",
	"BRK_MT_SETUP_FAIL_ACCESS_DENIED",
	"AST_MT_SETUP_ERR_INIT_FOHH_MCONN",
	"AST_MT_SETUP_ERR_MEM_LIMIT_REACHED",
	"BRK_MT_SETUP_FAIL_DUPLICATE_TAG_ID",
	"BRK_MT_AUTH_SAML_FAILURE",
	"AST_MT_SETUP_ERR_PRA_UNAVAILABLE",
	"C2C_MTUNNEL_NOT_FOUND",
	"MT_CLOSED_INTERNAL_ERROR",
	"DSP_MT_SETUP_FAIL_CANNOT_SEND_TO_BROKER",
	"CLT_READ_FAILED",
	"BRK_MT_SETUP_FAIL_CANNOT_SEND_TO_DISPATCHER",
	"AST_MT_SETUP_ERR_BROKER_BIND_FAIL",
	"BRK_MT_SETUP_FAIL_RATE_LIMIT_EXCEEDED",
	"CLT_INVALID_CLIENT",
	"BRK_MT_SETUP_FAIL_APP_NOT_FOUND",
	"C2C_NOT_AVAILABLE",
	"AST_MT_SETUP_ERR_MAX_SESSIONS_REACHED",
	"BRK_MT_AUTH_TWO_SAML_ASSERTION_IN_MSG",
	"AST_MT_SETUP_ERR_CPU_LIMIT_REACHED",
	"AST_MT_SETUP_ERR_NO_DNS_TO_SERVER",
	"CLT_PROTOCOL_NOT_SUPPORTED",
	"BRK_MT_AUTH_ALREADY_FAILED",
	"BRK_MT_SETUP_FAIL_CONNECTOR_GROUPS_MISSING",
	"BRK_MT_SETUP_FAIL_SCIM_INACTIVE",
	"EXPTR_MT_TLS_SETUP_FAIL_PEER",
	"BRK_MT_AUTH_SAML_DECODE_FAIL",
	"AST_MT_SETUP_ERR_BRK_HASH_TBL_FULL",
	"APP_NOT_REACHABLE",
	"BRK_MT_SETUP_TIMEOUT",
	"BRK_MT_TERMINATED_IDLE_TIMEOUT",
	"MT_CLOSED_DTLS_CONN_GONE_CLIENT_CLOSED",
	"MT_CLOSED_DTLS_CONN_GONE",
	"MT_CLOSED_DTLS_CONN_GONE_AST_CLOSED",
	"MT_CLOSED_TLS_CONN_GONE_SCIM_USER_DISABLE",
	"MT_CLOSED_TLS_CONN_GONE_CLIENT_CLOSED",
	"MT_CLOSED_TLS_CONN_GONE",
	"OPEN_OR_ACTIVE_CONNECTION",
	"MT_CLOSED_TLS_CONN_GONE_AST_CLOSED",
	"ZPN_ERR_SCIM_INACTIVE",
	"BRK_MT_CLOSED_FROM_ASSISTANT",
	"MT_CLOSED_TERMINATED",
	"AST_MT_TERMINATED",
	"BRK_MT_CLOSED_FROM_CLIENT",
	"BRK_MT_TERMINATED",
	"BRK_MT_SETUP_FAIL_NO_POLICY_FOUND",
	"BRK_MT_SETUP_FAIL_REJECTED_BY_POLICY",
	"BRK_MT_SETUP_FAIL_SAML_EXPIRED",
}

var supportedAppConnectorStatus = []string{
	"ZPN_STATUS_AUTH_FAILED", "ZPN_STATUS_DISCONNECTED", "ZPN_STATUS_AUTHENTICATED",
}

var supportedPrivateServiceEdgeStatus = []string{
	"ZPN_STATUS_AUTH_FAILED", "ZPN_STATUS_DISCONNECTED", "ZPN_STATUS_AUTHENTICATED",
}

var supportedLSSUserStatus = []string{
	"ZPN_STATUS_AUTH_FAILED", "ZPN_STATUS_DISCONNECTED", "ZPN_STATUS_AUTHENTICATED",
}

var noFilterSupportLogTypes = map[string]string{
	"zpn_http_trans_log":              "Web Browser",
	"zpn_audit_log":                   "Audit Logs",
	"zpn_waf_http_exchanges_log":      "AppProtection",
	"zpn_ast_comprehensive_stats":     "App Connector Metrics",
	"zpn_pbroker_comprehensive_stats": "Private Service Edge Metrics",
}

var supportedClientTypes = map[string]struct{}{
	"zpn_client_type_exporter":          {},
	"zpn_client_type_browser_isolation": {},
	"zpn_client_type_machine_tunnel":    {},
	"zpn_client_type_ip_anchoring":      {},
	"zpn_client_type_edge_connector":    {},
	"zpn_client_type_zapp":              {},
	"zpn_client_type_slogger":           {},
	"zpn_client_type_zapp_partner":      {},
	"zpn_client_type_branch_connector":  {},
}

var validObjectTypesForAuthLog = map[string]struct{}{
	"IDP":         {},
	"SCIM":        {},
	"SCIM_GROUP":  {},
	"SAML":        {},
	"CLIENT_TYPE": {},
}

var validObjectTypesForTransLog = map[string]struct{}{
	"IDP":         {},
	"SCIM":        {},
	"SCIM_GROUP":  {},
	"SAML":        {},
	"APP":         {},
	"APP_GROUP":   {},
	"CLIENT_TYPE": {},
}

func validateLSSConfigControllerFilters(sourceLogType, objectType, filter string, values []string, operands []lssconfigcontroller.PolicyRuleResourceOperands) error {
	// New logic to check if filters are not supported for specific log types
	if filter != "" {
		if displayName, found := noFilterSupportLogTypes[sourceLogType]; found {
			return fmt.Errorf("filter is not supported for source log type %s - %s", sourceLogType, displayName)
		}
	}
	// Filter validation
	if filter != "" {
		isValidFilter := false
		switch sourceLogType {
		case "zpn_auth_log":
			for _, validFilter := range supportedLSSUserStatus {
				if filter == validFilter {
					isValidFilter = true
					break
				}
			}
		case "zpn_trans_log":
			for _, validFilter := range supportedLSSUserActivity {
				if filter == validFilter {
					isValidFilter = true
					break
				}
			}
		case "zpn_ast_auth_log":
			for _, validFilter := range supportedAppConnectorStatus {
				if filter == validFilter {
					isValidFilter = true
					break
				}
			}
		case "zpn_sys_auth_log":
			for _, validFilter := range supportedPrivateServiceEdgeStatus {
				if filter == validFilter {
					isValidFilter = true
					break
				}
			}
		default:
			return fmt.Errorf("invalid source_log_type: %s", sourceLogType)
		}
		if !isValidFilter {
			return fmt.Errorf("invalid filter: %s for source_log_type: %s", filter, sourceLogType)
		}
	}
	// Object type validation
	if objectType != "" {
		var validObjectTypes map[string]struct{}
		switch sourceLogType {
		case "zpn_auth_log":
			validObjectTypes = validObjectTypesForAuthLog
		case "zpn_trans_log":
			validObjectTypes = validObjectTypesForTransLog
		default:
			return fmt.Errorf("invalid source_log_type: %s", sourceLogType)
		}

		if _, valid := validObjectTypes[objectType]; !valid {
			return fmt.Errorf("invalid object_type: %s for source_log_type: %s", objectType, sourceLogType)
		}

		// Value validation for CLIENT_TYPE
		if objectType == "CLIENT_TYPE" {
			for _, value := range values {
				if _, exists := supportedClientTypes[value]; !exists {
					return fmt.Errorf("invalid value: %s for object_type CLIENT_TYPE", value)
				}
			}
		}
	}

	// Entry value validation for specific object types
	for _, operand := range operands {
		objectType := operand.ObjectType // Overriding the outer objectType variable for clarity within this loop
		if objectType == "SCIM_GROUP" || objectType == "SCIM" || objectType == "SAML" || objectType == "IDP" {
			if operand.OperandsResourceLHSRHSValue == nil || len(*operand.OperandsResourceLHSRHSValue) == 0 {
				return fmt.Errorf("entry_values must be provided for object_type %s", objectType)
			}
			for _, entryValue := range *operand.OperandsResourceLHSRHSValue {
				if entryValue.LHS == "" || entryValue.RHS == "" {
					return errors.New("both lhs and rhs must be provided in entry_values")
				}
			}
		} else if operand.OperandsResourceLHSRHSValue != nil && len(*operand.OperandsResourceLHSRHSValue) != 0 {
			return fmt.Errorf("entry_values should not be provided for object_type %s", objectType)
		}
	}

	return nil
}
