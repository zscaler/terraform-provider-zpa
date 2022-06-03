---
subcategory: "Inspection"
layout: "zscaler"
page_title: "ZPA: zpa_inspection"
description: |-
  Get information about an Inspection Predefined Control in Zscaler Private Access cloud.
---

# Data Source: zpa_inspection_predefined_controls

Use the **zpa_inspection_predefined_controls** data source to get information about an OWASP predefined control and prefedined control version. This data source is required when creating an inspection profile.

## Example Usage

```hcl
data "zpa_inspection_predefined_controls" "example" {
    name = "Failed to parse request body"
    version = "OWASP_CRS/3.3.0"
}

output "zpa_inspection_predefined_controls" {
    value = data.zpa_inspection_predefined_controls.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the predefined control.
* `version` - (Required) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `action` - (Computed)
* `action_value` - (Computed)
* `associated_inspection_profile_names` - (Computed)
  * `id`- (Computed)
  * `name`- (Computed)
* `attachment` - (Computed)
* `control_group` - (Computed)
* `control_number` - (Computed)
* `creation_time` - (Computed)
* `default_action` - (Computed)
* `default_action_value` - (Computed)
* `description` - (Computed)
* `id` - (Computed)
* `modified_by` - (Computed)
* `modified_time` - (Computed)
* `paranoia_level` - (Computed)
* `severity` - (Computed)
