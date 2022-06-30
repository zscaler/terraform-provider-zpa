---
subcategory: "Inspection"
layout: "zscaler"
page_title: "ZPA: zpa_inspection"
description: |-
  Get information about all Inspection Predefined Control by group name in Zscaler Private Access cloud.
---

# Data Source: zpa_inspection_all_predefined_controls

Use the **zpa_inspection_all_predefined_controls** data source to get information about all OWASP predefined control and prefedined control version by group name. The `Preprocessors` predefined control is the default predefined control, This data source is always required, when creating an inspection profile.

## Example Usage

```hcl
data "zpa_inspection_all_predefined_controls" "this" {
  version    = "OWASP_CRS/3.3.0"
  group_name = "Preprocessors"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required) The name of the predefined control.
* `version` - (Required) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`

## Attribute Reference

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
* `severity` - (string)
  * `CRITICAL`
  * `ERROR`
  * `WARNING`
  * `INFO`
