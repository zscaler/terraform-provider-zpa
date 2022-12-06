---
subcategory: "Log Streaming (LSS)"
layout: "zscaler"
page_title: "ZPA: lss_config_controller"
description: |-
  Creates and manages ZPA LSS Configuration.
---

# Resource: zpa_lss_config_controller

The **zpa_lss_config_controller** resource creates and manages Log Streaming Service (LSS) in the Zscaler Private Access cloud.

## Example 1 Usage

```hcl
# Get Log Type Format
data "zpa_lss_config_log_type_formats" "zpn_ast_auth_log" {
  log_type = "zpn_ast_auth_log"
}

# Create Log Receiver Configuration
resource "zpa_lss_config_controller" "example" {
  config {
    name            = "Example"
    description     = "Example"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_ast_auth_log.json
    lss_host        = "splunk.acme.com"
    lss_port        = "11000"
    source_log_type = "zpn_ast_auth_log"
    use_tls         = true
    filter = [
      "ZPN_STATUS_AUTH_FAILED",
      "ZPN_STATUS_DISCONNECTED",
      "ZPN_STATUS_AUTHENTICATED"
    ]
  }
  connector_groups {
    id = [ zpa_app_connector_group.example.id ]
  }
}
```

## Example 2 Usage

```hcl
# Get Log Type Format
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
  log_type = "zpn_trans_log"
}

data "zpa_policy_type" "lss_siem_policy" {
    policy_type = "SIEM_POLICY"
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
  }
  policy_rule_resource {
    name   = "policy_rule_resource-lss_user_activity"
    action = "ALLOW"
    policy_set_id = data.zpa_policy_type.lss_siem_policy.id
    conditions {
      negated  = false
      operator = "OR"
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_exporter"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_ip_anchoring"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_zapp"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_edge_connector"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_machine_tunnel"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_browser_isolation"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_slogger"]
      }
    }
  }
  connector_groups {
    id = [ zpa_app_connector_group.example.id ]
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
  * `source_log_type` - (Required) Refer to the log type documentation
  * `connector_groups` - (Required)
        - `id` - (Required) - App Connector Group ID(s) where logs will be forwarded to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `config` - (Required)
  * `description` - (Optional)
  * `enabled` - (Optional)
  * `filter` - (Optional)
  * `use_tls` - (Optional)
  * `source_log_type` - (Required)
  * `connector_groups` - (Required)
        - `id` - (Required) - App Connector Group ID(s) where logs will be forwarded to.

* `policy_rule_resource` - (Optional)
  * `name` - (Optional)
  * `action` - (Optional)
  * `audit_message` - (Optional)
  * `custom_msg` - (Optional)
  * `connector_groups` - (Optional)
        - `id` - (Optional) - App Connector Group ID(s) where logs will be forwarded to.
  * `app_server_groups` - (Optional)
        - `id` - (Optional) - Server Group ID(s).
    * `conditions` - (Optional)
    * `negated` - (Optional)
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands`
      * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `SAML`, `SCIM`, `SCIM_GROUP`
      * `values` (Optional) The below values are supported when choosing `object_type` of type `CLIENT_TYPE`.
            - `zpn_client_type_exporter`
            - `zpn_client_type_browser_isolation`
            - `zpn_client_type_machine_tunnel`
            - `zpn_client_type_ip_anchoring`
            - `zpn_client_type_edge_connector`
            - `zpn_client_type_zapp`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)
