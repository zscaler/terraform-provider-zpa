---
page_title: "zpa_lss_config_controller Resource - terraform-provider-zpa"
subcategory: "Log Streaming (LSS)"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-log-streaming-service/API documentation https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api
  Creates and manages ZPA LSS Configuration for User Status.
---

# zpa_lss_config_controller (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-log-streaming-service)
* [API documentation](https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api)

The **zpa_lss_config_controller** resource creates and manages Log Streaming Service (LSS) in the Zscaler Private Access cloud for User Status `zpn_auth_log`.

## Example 1 - LSS User Status - Usage

```terraform
# Get Log Type Format - "User Status"
data "zpa_lss_config_log_type_formats" "zpn_auth_log" {
  log_type = "zpn_auth_log"
}

data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

# Retrieve the App Connector Group ID
data "zpa_app_connector_group" "this" {
 name = "Example100"
}

# Retrieve the Identity Provider ID
data "zpa_idp_controller" "this" {
 name = "Idp_Name"
}

# Retrieve the SCIM_GROUP ID(s)
data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "Idp_Name"
}

data "zpa_scim_groups" "sales" {
  name     = "Sales"
  idp_name = "Idp_Name"
}

resource "zpa_lss_config_controller" "lss_user_activity" {
  config {
    name            = "LSS User Status"
    description     = "LSS User Status"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_auth_log.json
    lss_host        = "splunk1.acme.com"
    lss_port        = "5001"
    source_log_type = "zpn_auth_log"
    use_tls         = true
    filter = ["ZPN_STATUS_AUTH_FAILED","ZPN_STATUS_DISCONNECTED", "ZPN_STATUS_AUTHENTICATED"]
  }
  policy_rule_resource {
    name          = "policy_rule_resource_lss_user_status"
    action        = "LOG"
    policy_set_id = data.zpa_policy_type.lss_siem_policy.id
    conditions {
      operator = "OR"
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_exporter", "zpn_client_type_browser_isolation", "zpn_client_type_machine_tunnel", "zpn_client_type_ip_anchoring", "zpn_client_type_edge_connector", "zpn_client_type_zapp", "zpn_client_type_slogger", "zpn_client_type_zapp_partner", "zpn_client_type_branch_connector"]
      }
    }
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
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.this.id ]
  }
}
```

## Schema

### Required

The following arguments are supported:

* `config` - (Required)
  * `name` - (Required)
  * `format` - (Required) The format of the LSS resource. The supported formats are: `JSON`, `CSV`, and `TSV`
  * `lss_host` - (Required) The IP or FQDN of the SIEM (Log Receiver) where logs will be forwarded to.
  * `lss_port` - (Required) The destination port of the SIEM (Log Receiver) where logs will be forwarded to.
  * `source_log_type` - (Required) For `User Status` logs use `zpn_auth_log`. Refer to the [Log Type documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_lss_config_log_type_formats).
  * `connector_groups` - (Required)
        - `id` - (Required) - App Connector Group ID(s) where logs will be forwarded to.

### Optional

In addition to all arguments above, the following attributes are exported:

* `config` - (Required)
  * `description` - (Optional)
  * `enabled` - (Optional)
  * `filter` - (Optional) - The following values are supported: `ZPN_STATUS_AUTH_FAILED`, `ZPN_STATUS_DISCONNECTED`, `ZPN_STATUS_AUTHENTICATED`.
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
      * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `CLIENT_TYPE`
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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)
