---
subcategory: "Inspection"
layout: "zscaler"
page_title: "ZPA: zpa_inspection"
description: |-
  Get information about an Inspection Custom Control in Zscaler Private Access cloud.
---

# Data Source: zpa_inspection_custom_controls

Use the **zpa_inspection_custom_controls** data source to get information about an inspection custom control. This data source can be associated with an inspection profile.

## Example Usage

```hcl
data "zpa_inspection_custom_controls" "example" {
    name = "ZPA_Inspection_Custom_Control"
}
```

```hcl
data "zpa_inspection_custom_controls" "example" {
    id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the predefined control.
* `version` - (Required) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `action` - (string) The performed action. Supported values: `PASS`, `BLOCK` and `REDIRECT`
* `action_value` - (string) Denotes the action
* `associated_inspection_profile_names` - (string) Name of the inspection profile
  * `id`- (string)
  * `name`- (string)
* `control_rule_json` (string) The control rule in JSON format that has the conditions and type of control for the inspection control
* `control_number` - (string)
* `creation_time` - (string)
* `default_action` - (string) The performed action. Supported values: `PASS`, `BLOCK` and `REDIRECT`
* `default_action_value` - (string) This is used to provide the redirect URL if the default action is set to `REDIRECT`
* `description` - (string) Description of the custom control
* `id` - (string)
* `name` - (string) Name of the custom control
* `modified_by` - (string)
* `modified_time` - (string)
* `paranoia_level` - (string) OWASP Predefined Paranoia Level.
* `severity` - (string) Severity of the control number. Supported values: `CRITICAL`, `ERROR`, `WARNING`, `INFO`
* `type` - (string) Rules to be applied to the request or response type
* `rules` - (string) Rules of the custom controls applied as conditions `JSON`
  * `conditions` - (string)
    * `lhs` - (string) Signifies the key for the object type Supported values: `SIZE`, `VALUE`
    * `op` - (string) If lhs is set to SIZE, then the user may pass one of the following: `EQ: Equals`, `LE: Less than or equal to`, `GE: Greater than or equal to`. If the lhs is set to `VALUE`, then the user may pass one of the following: `CONTAINS`, `STARTS_WITH`, `ENDS_WITH`, `RX`.
    * `rhs` - (string) Denotes the value for the given object type. Its value depends on the key. If rules.type is set to REQUEST_METHOD, the conditions.rhs field must have one of the following values: `GET`,`HEAD`, `POST`, `OPTIONS`, `PUT`, `DELETE`, `TRACE`
  * `names` - (string) Name of the rules. If rules.type is set to `REQUEST_HEADERS`, `REQUEST_COOKIES`, or `RESPONSE_HEADERS`, the rules.name field is required.
  * `type` - (string) Type value for the rules
