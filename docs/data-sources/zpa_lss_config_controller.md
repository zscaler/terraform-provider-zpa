---
page_title: "zpa_lss_config_controller Data Source - terraform-provider-zpa"
subcategory: "Log Streaming (LSS)"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-log-streaming-service/
  API documentation https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api
  Get information about Log Streaming (LSS) configuration Zscaler Private Access cloud.
---

# zpa_lss_config_controller (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-log-streaming-service)
* [API documentation](https://help.zscaler.com/zpa/configuring-log-streaming-service-configurations-using-api)

Use the **zpa_lss_config_controller** data source to get information about a Log Streaming (LSS) configuration resource created in the Zscaler Private Access.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
# Retrieve Log Streaming Information by Name
data "zpa_lss_config_controller" "example" {
  name = "testAcc-lss-server"
}
```

```terraform
# Retrieve Log Streaming Information by ID
data "zpa_lss_config_controller" "example" {
  id = "1234567890"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) This field defines the name of the log streaming resource.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) This field defines the name of the log streaming resource.
* `config` - (Computed)
  * `audit_message` - (string)
  * `creation_time` - (string)
  * `description` - (string)
  * `enabled` - (bool)
  * `filter` - (string)
  * `format` - (string)
  * `id` - (string)
  * `modified_by` - (string)
  * `modified_time` - (string)
  * `name` - (string)
  * `lss_host` - (string)
  * `lss_port` - (string)
  * `source_log_type` - (string)

* `connector_groups` - (Computed)
  * `id` - (string)
