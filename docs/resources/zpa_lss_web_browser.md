---
page_title: "zpa_lss_config_controller Resource - terraform-provider-zpa"
subcategory: "Log Streaming (LSS)"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-log-streaming-service/API documentation https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api
  Creates and manages ZPA LSS Configuration for Browser Access.
---

# zpa_lss_config_controller (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-log-streaming-service)
* [API documentation](https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api)

The **zpa_lss_config_controller** resource creates and manages Log Streaming Service (LSS) in the Zscaler Private Access cloud for Browser Access `zpn_http_trans_log`.

## Example 1 - LSS Browser Access - Usage

```terraform
# Get Log Type Format - "Browser Access"
data "zpa_lss_config_log_type_formats" "zpn_http_trans_log" {
  log_type = "zpn_http_trans_log"
}

data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

data "zpa_app_connector_group" "this" {
 name = "Example100"
}
resource "zpa_lss_config_controller" "lss_browser_access" {
  config {
    name            = "LSS Browser Access"
    description     = "LSS Browser Access"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_http_trans_log.json
    lss_host        = "splunk1.acme.com"
    lss_port        = "5001"
    source_log_type = "zpn_http_trans_log"
    use_tls         = true
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
  * `source_log_type` - (Required) For `Private Service Edge Status` logs use `zpn_sys_auth_log`. Refer to the [Log Type documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_lss_config_log_type_formats).
  * `connector_groups` - (Required)
        - `id` - (Required) - App Connector Group ID(s) where logs will be forwarded to.

### Optional

In addition to all arguments above, the following attributes are exported:

* `config` - (Required)
  * `description` - (Optional)
  * `enabled` - (Optional)
  * `use_tls` - (Optional)
  * `source_log_type` - (Required) For `Web Browser` logs use `zpn_http_trans_log`. Refer to the [Log Type documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_lss_config_log_type_formats).
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

## LSS Source Log Type Table

|       Source Log Type                     |            Description                 |
|-------------------------------------------|----------------------------------------|
|        `zpn_trans_log`                    |        `User Activity`                 |
|        `zpn_auth_log`                     |         `User Status`                  |
|        `zpn_ast_auth_log`                 |        `App Connector Status`          |
|        `zpn_http_trans_log`               |         `Web Browser`                  |
|        `zpn_audit_log`                    |         `Audit Logs`                   |
|        `zpn_sys_auth_log`                 |         `Private Service Edge Status`  |
|        `zpn_ast_comprehensive_stats`      |         `App Connector Metrics`        |
|        `zpn_pbroker_comprehensive_stats`  |         `Private Service Edge Metrics` |
|        `zpn_waf_http_exchanges_log`       |         `ZPA App Protection`           |

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)
