---
subcategory: "LSS Config Controller"
layout: "zpa"
page_title: "ZPA: lss_config_controller"
description: |-
  Gets details of a Log Streaming (LSS) configuration resource.
---

# zpa_lss_config_controller

The **zpa_lss_config_controller** data source provides details about a specific Log Streaming (LSS) configuration resource created in the Zscaler Private Access.

## Example Usage

```hcl
# Retrieve Log Streaming Information
data "zpa_lss_config_controller" "example" {
  id = zpa_lss_config_controller.example
}

output "zpa_lss_config_controller" {
  value = data.zpa_lss_config_controller.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This field defines the name of the log streaming resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String)
* `enabled` - (Boolean)
* `format` - (String)
* `lss_host` - (String)
* `lss_port` - (String)
* `lss_port` - (String)
* `filter` - (String)
* `source_log_type` - (String)
