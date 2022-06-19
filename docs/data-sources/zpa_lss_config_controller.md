---
subcategory: "Log Streaming (LSS)"
layout: "zscaler"
page_title: "ZPA: lss_config_controller"
description: |-
  Get information about Log Streaming (LSS) configuration Zscaler Private Access cloud.
---

# Data Source: zpa_lss_config_controller

Use the **zpa_lss_config_controller** data source to get information about a Log Streaming (LSS) configuration resource created in the Zscaler Private Access.

## Example Usage

```hcl
# Retrieve Log Streaming Information by Name
data "zpa_lss_config_controller" "example" {
  name = "testAcc-lss-server"
}
```

```hcl
# Retrieve Log Streaming Information by ID
data "zpa_lss_config_controller" "example" {
  id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This field defines the name of the log streaming resource.
* `id` - (Optional) This field defines the name of the log streaming resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

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
