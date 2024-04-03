---
page_title: "zpa_inspection_predefined_controls Data Source - terraform-provider-zpa"
subcategory: "AppProtection"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-custom-controls/API documentation https://help.zscaler.com/zpa/configuring-appprotection-controls-using-api
  Get information about an Inspection Predefined Control in Zscaler Private Access cloud.
---

# zpa_inspection_predefined_controls (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-custom-controls)
* [API documentation](https://help.zscaler.com/zpa/configuring-appprotection-controls-using-api)

Use the **zpa_inspection_predefined_controls** data source to get information about an OWASP predefined control and prefedined control version. This data source is required when creating an inspection profile.

## Example Usage

```terraform
data "zpa_inspection_predefined_controls" "example" {
    name = "Failed to parse request body"
    version = "OWASP_CRS/3.3.0"
}

output "zpa_inspection_predefined_controls" {
    value = data.zpa_inspection_predefined_controls.example
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the predefined control.
* `version` - (Required) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `action` - (Computed)
* `action_value` - (Computed)
* `associated_inspection_profile_names` - (Computed)
  * `id`- (Computed)
  * `name`- (Computed)
* `attachment` - (Computed)
* `control_group` - (Computed)
* `control_number` - (Computed)
* `control_type` - (Computed)
* `creation_time` - (Computed)
* `default_action` - (Computed)
* `default_action_value` - (Computed)
* `description` - (Computed)
* `id` - (Computed)
* `modified_by` - (Computed)
* `modified_time` - (Computed)
* `paranoia_level` - (Computed)
* `protocol_type` - (Computed)
* `severity` - (Computed)
