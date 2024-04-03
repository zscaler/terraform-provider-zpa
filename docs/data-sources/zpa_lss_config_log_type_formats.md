---
page_title: "zpa_lss_config_log_type_formats Data Source - terraform-provider-zpa"
subcategory: "Log Streaming (LSS)"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-log-streaming-service/
  API documentation https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api
  Get information about all all LSS log type format details.
---

# zpa_lss_config_log_type_formats (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-log-streaming-service)
* [API documentation](https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api)

Use the **zpa_lss_config_log_type_formats** data source to get information about all LSS log type formats in the Zscaler Private Access cloud. This data source is required when creating an LSS Config Controller resource.

## Example Usage

```terraform
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
  log_type = "zpn_trans_log"
}

data "zpa_lss_config_log_type_formats" "zpn_auth_log" {
  log_type = "zpn_auth_log"
}

data "zpa_lss_config_log_type_formats" "zpn_ast_auth_log" {
  log_type = "zpn_ast_auth_log"
}

data "zpa_lss_config_log_type_formats" "zpn_http_trans_log" {
  log_type = "zpn_http_trans_log"
}

data "zpa_lss_config_log_type_formats" "zpn_audit_log" {
  log_type = "zpn_audit_log"
}

data "zpa_lss_config_log_type_formats" "zpn_sys_auth_log" {
  log_type = "zpn_sys_auth_log"
}

data "zpa_lss_config_log_type_formats" "zpn_ast_comprehensive_stats" {
  log_type = "zpn_ast_comprehensive_stats"
}

data "zpa_lss_config_log_type_formats" "zpn_waf_http_exchanges_log" {
 log_type = "zpn_waf_http_exchanges_log"
}

data "zpa_lss_config_log_type_formats" "zpn_pbroker_comprehensive_stats" {
  log_type = "zpn_pbroker_comprehensive_stats"
}
```

### Read-Only

The following arguments are supported:

* `log_type` - (Required) The type of log to be exported.
  * `zpn_trans_log`
  * `zpn_auth_log`
  * `zpn_ast_auth_log`
  * `zpn_http_trans_log`
  * `zpn_audit_log`
  * `zpn_sys_auth_log`
  * `zpn_ast_comprehensive_stats`
  * `zpn_waf_http_exchanges_log`
  * `zpn_pbroker_comprehensive_stats`
