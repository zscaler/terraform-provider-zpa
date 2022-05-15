---
layout: "zscaler"
page_title: "Zscaler Private Access (ZPA): lss_config_controller"
sidebar_current: "docs-datasource-zpa-lss-config-controller"
description: |-
  Get information about Log Streaming (LSS) configuration Zscaler Private Access cloud.
---

# zpa_lss_config_controller

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
  * `audit_message` - (Computed)
  * `creation_time` - (Computed)
  * `description` - (Computed)
  * `enabled` - (Computed)
  * `filter` - (Computed)
  * `format` - (Computed)
  * `id` - (Computed)
  * `modified_by` - (Computed)
  * `modified_time` - (Computed)
  * `name` - (Computed)
  * `lss_host` - (Computed)
  * `lss_port` - (Computed)
  * `source_log_type` - (Computed)

* `connector_groups` - (Computed)
  * `id` - (Computed)
