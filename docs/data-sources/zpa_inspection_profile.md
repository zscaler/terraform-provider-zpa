---
page_title: "zpa_inspection_profile Data Source - terraform-provider-zpa"
subcategory: "AppProtection"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-browser-protection-profiles
  API documentation https://help.zscaler.com/zpa/configuring-appprotection-profiles-using-api
  Get information about an Inspection Profile in Zscaler Private Access cloud.
---

# zpa_inspection_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-browser-protection-profiles)
* [API documentation](https://help.zscaler.com/zpa/configuring-appprotection-profiles-using-api)

Use the **zpa_inspection_profile** data source to get information about an inspection profile in the Zscaler Private Access cloud. This resource can then be referenced in an inspection custom control resource.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
data "zpa_inspection_profile" "this" {
  name = "Example"
}
```

## Schema

### Required

* `name` - (String) This field defines the name of the inspection profile.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) This field defines the id of the inspection profile.
* `description` - (string) Description of the inspection profile.
* `paranoia_level` - (string) OWASP Predefined Paranoia Level. Range: [1-4], inclusive
* `predefined_controls` - (string) The predefined controls
  * `id` - (string) ID of the predefined control
  * `action` - (string) The action of the predefined control. Supported values: `PASS`, `BLOCK` and `REDIRECT`
  * `action_value` - (string) Value for the predefined controls action. This field is only required if the action is set to REDIRECT. This field is only required if the action is set to `REDIRECT`.
  * `attachment` (string) Control attachment
  * `control_group` (string) Control group

* `custom_controls` - (string) Types for custom controls
  * `type` (string) Types for custom controls
  * `control_rule_json` (string) Custom controls string in JSON format
  * `rules` - (string) Rules of the custom controls applied as conditions `JSON`
    * `conditions` - (string)
      * `lhs` - (string) Signifies the key for the object type Supported values: `SIZE`, `VALUE`
      * `op` - (string) If lhs is set to SIZE, then the user may pass one of the following: `EQ: Equals`, `LE: Less than or equal to`, `GE: Greater than or equal to`. If the lhs is set to `VALUE`, then the user may pass one of the following: `CONTAINS`, `STARTS_WITH`, `ENDS_WITH`, `RX`.
      * `rhs` - (string) Denotes the value for the given object type. Its value depends on the key. If rules.type is set to REQUEST_METHOD, the conditions.rhs field must have one of the following values: `GET`,`HEAD`, `POST`, `OPTIONS`, `PUT`, `DELETE`, `TRACE`
    * `names` - (string) Name of the rules. If rules.type is set to `REQUEST_HEADERS`, `REQUEST_COOKIES`, or `RESPONSE_HEADERS`, the rules.name field is required.
    * `type` - (string) Type value for the rules
    * `version` - (string) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`

* `associated_inspection_profile_names` - (string) Name of the inspection profile
  * `id`- (string)
  * `name`- (string)

* `common_global_override_actions_config` - (string)
* `controls_info` - (string) Types for custom controls
  * `control_type` - (string) Control types. Supported Values: `WEBSOCKET_PREDEFINED`, `WEBSOCKET_CUSTOM`, `CUSTOM`, `PREDEFINED`, `ZSCALER`
  * `count` - (string) Control information counts `Long`
* `web_socket_controls` - (string)
  * `id` - (string) ID of the predefined control
  * `action` - (string) The action of the predefined control. Supported values: `PASS`, `BLOCK` and `REDIRECT`
  * `action_value` - (string) Value for the predefined controls action. This field is only required if the action is set to REDIRECT. This field is only required if the action is set to `REDIRECT`.
