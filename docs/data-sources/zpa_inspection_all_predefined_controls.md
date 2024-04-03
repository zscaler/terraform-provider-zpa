---
page_title: "zpa_inspection_all_predefined_controls Data Source - terraform-provider-zpa"
subcategory: "AppProtection"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-custom-controls/API documentation https://help.zscaler.com/zpa/configuring-appprotection-controls-using-api
  Get information about all Inspection Predefined Control by group name in Zscaler Private Access cloud.
---

# zpa_inspection_all_predefined_controls (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-custom-controls)
* [API documentation](https://help.zscaler.com/zpa/configuring-appprotection-controls-using-api)

Use the **zpa_inspection_all_predefined_controls** data source to get information about all OWASP predefined control and prefedined control version by group name. The `Preprocessors` predefined control is the default predefined control, This data source is always required, when creating an inspection profile.

## Example Usage

```terraform
data "zpa_inspection_all_predefined_controls" "this" {
  version    = "OWASP_CRS/3.3.0"
  group_name = "Preprocessors"
}
```

## Schema

### Required

The following arguments are supported:

* `group_name` - (Required) The name of the predefined control.
* `version` - (Required) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `action` - (string)
  * `PASS`
  * `BLOCK`
  * `REDIRECT`
* `action_value` - (string)
* `associated_inspection_profile_names` - (string)
  * `id`- (string)
  * `name`- (string)
* `attachment` - (string)
* `control_group` - (string)
* `control_number` - (string)
* `control_type` - (string) Returned values: `WEBSOCKET_PREDEFINED`, `WEBSOCKET_CUSTOM`, `ZSCALER`, `CUSTOM`, `PREDEFINED`
* `creation_time` - (string)
* `default_action` - (string)
  * `PASS`
  * `BLOCK`
  * `REDIRECT`
* `default_action_value` - (string)
* `description` - (string)
* `id` - (string)
* `modified_by` - (string)
* `modified_time` - (string)
* `paranoia_level` - (string)
* `protocol_type` - (string) Returned values: `HTTP`, `HTTPS`, `FTP`, `RDP`, `SSH`, `WEBSOCKET`
* `severity` - (string)
  * `CRITICAL`
  * `ERROR`
  * `WARNING`
  * `INFO`
