---
subcategory: "Log Streaming (LSS)"
layout: "zscaler"
page_title: "ZPA: lss_config_controller"
description: |-
  Creates and manages ZPA LSS Configuration for User Activity.
---

# Resource: zpa_lss_config_controller

The **zpa_lss_config_controller** resource creates and manages Log Streaming Service (LSS) in the Zscaler Private Access cloud for User Activity `zpn_trans_log`.

## Example 1 - LSS User Activity - Usage

```hcl
# Retrieve the Log Type Format
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
  log_type = "zpn_trans_log"
}

# Retrieve the Policy Set ID from Policy Type `SIEM_POLICY`
data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

# Retrieve the App Connector Group ID
data "zpa_app_connector_group" "this" {
 name = "Example100"
}

# Retrieve the Application Segment(s) ID
data "zpa_application_segment" "app01" {
  name = "app01"
}

data "zpa_application_segment" "app02" {
  name = "app02"
}

# Retrieve the Segment Group(s) ID
data "zpa_segment_group" "this" {
  name = "Example100"
}

# Retrieve the Identity Provider ID
data "zpa_idp_controller" "this" {
 name = "IdP_Name"
}

# Retrieve the SCIM_GROUP ID(s)
data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "IdP_Name"
}

data "zpa_scim_groups" "sales" {
  name     = "Sales"
  idp_name = "IdP_Name"
}

resource "zpa_lss_config_controller" "lss_user_activity" {
  config {
    name            = "LSS User Activity"
    description     = "LSS User Activity"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_trans_log.json
    lss_host        = "splunk.acme.com"
    lss_port        = "11001"
    source_log_type = "zpn_trans_log"
    use_tls         = true
    filter = ["BRK_MT_SETUP_FAIL_BIND_TO_AST_LOCAL_OWNER", "CLT_INVALID_DOMAIN","AST_MT_SETUP_ERR_HASH_TBL_FULL", "AST_MT_SETUP_ERR_CONN_PEER", "BRK_MT_SETUP_FAIL_REJECTED_BY_POLICY_APPROVAL", "BRK_MT_SETUP_FAIL_ICMP_RATE_LIMIT_NUM_APP_EXCEEDED"
  ]
  }
  policy_rule_resource {
    name          = "policy_rule_resource-lss_user_activity"
    action        = "LOG"
    policy_set_id = data.zpa_policy_type.lss_siem_policy.id
    conditions {
      operator = "OR"
      operands {
        object_type = "SCIM_GROUP"
        entry_values {
          rhs = data.zpa_scim_groups.engineering.id
          lhs = data.zpa_idp_controller.this.id
        }
        entry_values {
          rhs = data.zpa_scim_groups.sales.id
          lhs = data.zpa_idp_controller.this.id
        }
      }
    }
    conditions {
      operator = "OR"
      operands {
        object_type = "APP"
        values      = [ data.zpa_application_segment.app01.id, data.zpa_application_segment.app02.id ]
      }
      operands {
        object_type = "APP_GROUP"
        values      = [ data.zpa_segment_group.this.id ]
      }
    }
    conditions {
      operator = "OR"
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_exporter", "zpn_client_type_ip_anchoring", "zpn_client_type_zapp", "zpn_client_type_edge_connector", "zpn_client_type_machine_tunnel", "zpn_client_type_browser_isolation", "zpn_client_type_slogger", "zpn_client_type_zapp_partner", "zpn_client_type_branch_connector"]
      }
    }
  }
  connector_groups {
    id = [data.zpa_app_connector_group.this.id]
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `config` - (Required)
  * `name` - (Required)
  * `format` - (Required) The format of the LSS resource. The supported formats are: `JSON`, `CSV`, and `TSV`
  * `lss_host` - (Required) The IP or FQDN of the SIEM (Log Receiver) where logs will be forwarded to.
  * `lss_port` - (Required) The destination port of the SIEM (Log Receiver) where logs will be forwarded to.
  * `source_log_type` - (Required) For `User Activity` logs use `zpn_trans_log`. Refer to the [Log Type documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_lss_config_log_type_formats).
  * `connector_groups` - (Required)
        - `id` - (Required) - App Connector Group ID(s) where logs will be forwarded to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `config` - (Required)
  * `description` - (Optional)
  * `enabled` - (Optional)
  * `filter` - (Optional) - Refers to the session status codes to exclude from the Log Streaming Service. If you don't make any selections, then the policy includes all status codes. [LSS Filter Status Codes](#lss-filter-status-codes)
    * Refer to the filter table for reference.
  * `use_tls` - (Optional)
  * `source_log_type` - (Required) Refer to the [Log Type documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_lss_config_log_type_formats).
    * `zpn_trans_log - "User Activity"`
    * `zpn_auth_log - "User Status"`
    * `zpn_ast_auth_log - "App Connector Status"`
    * `zpn_http_trans_log - "Web Browser"`
    * `zpn_audit_log - "Audit Logs"`
    * `zpn_sys_auth_log - "Private Service Edge Status"`
    * `zpn_ast_comprehensive_stats - "App Connector Metrics"`
    * `zpn_pbroker_comprehensive_stats - "Private Service Edge Metrics"`
    * `zpn_waf_http_exchanges_log`
  * `connector_groups` - (Required)
        - `id` - (Required) - App Connector Group ID(s) where logs will be forwarded to.

* `policy_rule_resource` - (Optional)
  * `name` - (Optional)
  * `action` - (Optional) - Supported Value(s) are: `LOG`
  * `audit_message` - (Optional)
  * `custom_msg` - (Optional)
    * `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
      * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `CLIENT_TYPE`
      * `values` (Optional) The below values are supported when choosing `object_type` of type `CLIENT_TYPE`.
            - `zpn_client_type_exporter - "Web Browser"`
            - `zpn_client_type_browser_isolation - "Cloud Browser"`
            - `zpn_client_type_machine_tunnel - "Machine Tunnel"`
            - `zpn_client_type_ip_anchoring - "ZIA Service Edge"`
            - `zpn_client_type_edge_connector - "Cloud Connector"`
            - `zpn_client_type_branch_connector - "Branch Connector"`
            - `zpn_client_type_zapp - "Client Connector"`
            - `zpn_client_type_slogger - "ZPA LSS"`
            - `zpn_client_type_zapp_partner - "Client Connector Partner"`

    * `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
      * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `SCIM`, `SCIM_GROUP`, `SAML`, `IDP`
      * `entry_values` (Optional)
        * `lhs` - (Optional) -  The Identity Provider ID
        * `rhs` - (Optional) - The SCIM, SCIM Group, or SAML Attribute ID

## LSS Filter Status Codes

|                     LSS Filter Options                                      |
|:---------------------------------------------------------------------------:|
| `BRK_MT_SETUP_FAIL_BIND_TO_AST_LOCAL_OWNER`,                                |
| `CLT_INVALID_DOMAIN, AST_MT_SETUP_ERR_HASH_TBL_FULL`                        |
| `BRK_MT_SETUP_FAIL_REJECTED_BY_POLICY_APPROVAL`                             |
| `BRK_MT_SETUP_FAIL_ICMP_RATE_LIMIT_NUM_APP_EXCEEDED`,                       |
| `EXPTR_MT_TLS_SETUP_FAIL_VERSION_MISMATCH`,                                 |
| `BRK_MT_SETUP_FAIL_RATE_LIMIT_LOOP_DETECTED`,                               |
| `CLT_INVALID_TAG`, `AST_MT_SETUP_ERR_NO_SYSTEM_FD`                          |
| `AST_MT_SETUP_ERR_NO_PROCESS_FD`, `BROKER_NOT_ENABLED`                      |
| `BROKER_NOT_ENABLED`, `AST_MT_SETUP_ERR_AST_CFG_DISABLED`,                  |
| `BRK_MT_SETUP_FAIL_TOO_MANY_FAILED_ATTEMPTS`,                               |
| `BRK_MT_AUTH_NO_SAML_ASSERTION_IN_MSG`,                                     |
| `INVALID_DOMAIN`, `CLT_WRONG_PORT`                                          |
| `BRK_MT_TERMINATED_BRK_SWITCHED`, `AST_MT_SETUP_ERR_CONN_PEER`              |
| `AST_MT_SETUP_ERR_OPEN_SERVER_CLOSE`,                                       |
| `AST_MT_SETUP_ERR_BIND_TO_AST_LOCAL_OWNER`, `NO_CONNECTOR_AVAILABLE`        |
| `BRK_MT_AUTH_SAML_CANNOT_ADD_ATTR_TO_HEAP`,                                 |
| `EXPTR_MT_TLS_SETUP_FAIL_NOT_TRUSTED_CA`,                                   |
| `AST_MT_SETUP_TIMEOUT_NO_ACK_TO_BIND`                                       |
| `CLT_PORT_UNREACHABLE`, `C2C_CLIENT_CONN_EXPIRED`                           |
| `BRK_MT_SETUP_FAIL_BIND_TO_CLIENT_LOCAL_OWNER`,                             |
| `BRK_MT_AUTH_SAML_CANNOT_ADD_ATTR_TO_HASH`,                                 |
| `BRK_MT_SETUP_FAIL_REPEATED_DISPATCH`,                                      |
| `AST_MT_SETUP_ERR_OPEN_SERVER_ERROR`, `DSP_MT_SETUP_FAIL_DISCOVERY_TIMEOUT` |
| `CUSTOMER_NOT_ENABLED`, `BRK_CONN_UPGRADE_REQUEST_FAILED`                   |
| `C2C_MTUNNEL_FAILED_FORWARD`, `EXPTR_MT_TLS_SETUP_FAIL_CERT_CHAIN_ISSUE`,   |
| `AST_MT_SETUP_ERR_RATE_LIMIT_REACHED`,                                      |
| `BRK_MT_SETUP_FAIL_RATE_LIMIT_NUM_APP_EXCEEDED`,                            |
| `AST_MT_SETUP_TIMEOUT_CANNOT_CONN_TO_SERVER`                                |
| `BRK_MT_AUTH_SAML_FINGER_PRINT_FAIL`,                                       |
| `AST_MT_SETUP_ERR_NO_EPHEMERAL_PORT`,                                       |
| `BRK_CONN_UPGRADE_REQUEST_FORBIDDEN`,                                       |
| `AST_MT_SETUP_ERR_OPEN_SERVER_CONN`, `CLT_PROBE_FAILED`                     |
| `AST_MT_SETUP_ERR_APP_NOT_FOUND`, `AST_MT_SETUP_ERR_OPEN_BROKER_CONN`       |
| `BRK_MT_SETUP_FAIL_ICMP_RATE_LIMIT_EXCEEDED`                                |
| `AST_MT_SETUP_ERR_OPEN_SERVER_TIMEOUT`,                                     |
| `C2C_MTUNNEL_BAD_STATE`, `CLT_DUPLICATE_TAG`, `AST_MT_SETUP_TIMEOUT`        |
| `CLT_DOUBLEENCRYPT_NOT_SUPPORTED`,                                          |
| `BRK_MT_SETUP_FAIL_CANNOT_SEND_MT_COMPLETE`,                                |
| `BRK_MT_SETUP_FAIL_BIND_RECV_IN_BAD_STATE`,                                 |
| `APP_NOT_AVAILABLE`, `BRK_MT_AUTH_SAML_NO_USER_ID`,                         |
| `AST_MT_SETUP_TIMEOUT_CANNOT_CONN_TO_BROKER`,                               |
| `DSP_MT_SETUP_FAIL_MISSING_HEALTH`, `AST_MT_SETUP_ERR_DUP_MT_ID`            |
| `AST_MT_SETUP_ERR_BIND_GLOBAL_OWNER`, `AST_MT_SETUP_ERR_BIND_ACK`           |
| `BRK_MT_TERMINATED_APPROVAL_TIMEOUT`, `CLT_CONN_FAILED`                     |
| `BRK_MT_SETUP_FAIL_ACCESS_DENIED`,                                          |
| `AST_MT_SETUP_ERR_INIT_FOHH_MCONN`,                                         |
| `AST_MT_SETUP_ERR_MEM_LIMIT_REACHED`,                                       |
| `BRK_MT_SETUP_FAIL_DUPLICATE_TAG_ID`, `BRK_MT_AUTH_SAML_FAILURE`            |
| `AST_MT_SETUP_ERR_PRA_UNAVAILABLE`,                                         |
| `C2C_MTUNNEL_NOT_FOUND`, `MT_CLOSED_INTERNAL_ERROR`                         |
| `DSP_MT_SETUP_FAIL_CANNOT_SEND_TO_BROKER`, `CLT_READ_FAILED`                |
| `BRK_MT_SETUP_FAIL_CANNOT_SEND_TO_DISPATCHER`,                              |
| `AST_MT_SETUP_ERR_BROKER_BIND_FAIL`, `CLT_INVALID_CLIENT`                   |
| `BRK_MT_SETUP_FAIL_RATE_LIMIT_EXCEEDED`,                                    |
| `BRK_MT_SETUP_FAIL_APP_NOT_FOUND`, `C2C_NOT_AVAILABLE`                      |
| `AST_MT_SETUP_ERR_MAX_SESSIONS_REACHED`,                                    |
| `BRK_MT_AUTH_TWO_SAML_ASSERTION_IN_MSG`,                                    |
| `AST_MT_SETUP_ERR_CPU_LIMIT_REACHED`,                                       |
| `AST_MT_SETUP_ERR_NO_DNS_TO_SERVER`,                                        |
| `CLT_PROTOCOL_NOT_SUPPORTED`, `BRK_MT_AUTH_ALREADY_FAILED`                  |
| `BRK_MT_SETUP_FAIL_CONNECTOR_GROUPS_MISSING`,                               |
| `BRK_MT_SETUP_FAIL_SCIM_INACTIVE`,                                          |
| `EXPTR_MT_TLS_SETUP_FAIL_PEER`, `BRK_MT_AUTH_SAML_DECODE_FAIL`              |
| `AST_MT_SETUP_ERR_BRK_HASH_TBL_FULL`, `APP_NOT_REACHABLE`                   |
| `BRK_MT_TERMINATED_IDLE_TIMEOUT`,                                           |
| `MT_CLOSED_DTLS_CONN_GONE_CLIENT_CLOSED`, `MT_CLOSED_DTLS_CONN_GONE`        |
| `MT_CLOSED_DTLS_CONN_GONE_AST_CLOSED`,                                      |
| `MT_CLOSED_TLS_CONN_GONE_SCIM_USER_DISABLE`,                                |
| `MT_CLOSED_TLS_CONN_GONE_CLIENT_CLOSED`,                                    |
| `MT_CLOSED_TLS_CONN_GONE`, `OPEN_OR_ACTIVE_CONNECTION`                      |
| `MT_CLOSED_TLS_CONN_GONE_AST_CLOSED`,                                       |
| `ZPN_ERR_SCIM_INACTIVE`, `BRK_MT_CLOSED_FROM_ASSISTANT`,                    |
| `MT_CLOSED_TERMINATED`, `AST_MT_TERMINATED`                                 |
| `BRK_MT_CLOSED_FROM_CLIENT`, `BRK_MT_TERMINATED`                            |
| `BRK_MT_SETUP_FAIL_NO_POLICY_FOUND`,                                        |
| `BRK_MT_SETUP_FAIL_REJECTED_BY_POLICY`,                                     |
| `BRK_MT_SETUP_FAIL_SAML_EXPIRED`,                                           |

## LSS Source Log Type Table

|       Source Log Type                         |            Description                    |
|-----------------------------------------------|-------------------------------------------|
|        `zpn_trans_log`                        |        `User Activity`                    |
|        `zpn_auth_log`                         |         `User Status`                     |
|        `zpn_ast_auth_log`                     |        `App Connector Status`             |
|        `zpn_http_trans_log`                   |         `Web Browser`                     |
|        `zpn_audit_log`                        |         `Audit Logs`                      |
|        `zpn_sys_auth_log`                     |         `Private Service Edge Status`     |
|        `zpn_ast_comprehensive_stats`          |         `App Connector Metrics`           |
|        `zpn_pbroker_comprehensive_stats`      |         `Private Service Edge Metrics`    |
|        `zpn_waf_http_exchanges_log`           |         `ZPA App Protection`              |

## LSS Object Type `CLIENT_TYPE` Values Table

|       Source Log Type                          |          Description                       |
|----------------------------------------------- |--------------------------------------------|
|        `zpn_client_type_exporter`              |          `Web Browser`                     |
|        `zpn_client_type_browser_isolation`     |         `Cloud Browser`                    |
|        `zpn_client_type_machine_tunnel`        |         `Machine Tunnel`                   |
|        `zpn_client_type_ip_anchoring`          |         `ZIA Service Edge`                 |
|        `zpn_client_type_edge_connector`        |         `Cloud Connector`                  |
|        `zpn_client_type_branch_connector`      |         `Branch Connector`                 |
|        `zpn_client_type_zapp`                  |         `Client Connector`                 |
|        `zpn_client_type_slogger`               |              `ZPA LSS`                     |
|        `zpn_client_type_zapp_partner`          |         `Client Connector Partner`         |

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)
