---
page_title: "zpa_inspection_profile Resource - terraform-provider-zpa"
subcategory: "AppProtection"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-browser-protection-profiles/API documentation https://help.zscaler.com/zpa/configuring-appprotection-profiles-using-api
  Creates and manages Inspection Profile in Zscaler Private Access cloud.
---

# zpa_inspection_profile (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-browser-protection-profiles)
* [API documentation](https://help.zscaler.com/zpa/configuring-appprotection-profiles-using-api)

The  **zpa_inspection_profile** resource creates an inspection profile in the Zscaler Private Access cloud. This resource can then be referenced in an inspection custom control resource.

**NOTE** There are several ways to set up the Inspection Profile due to its complex data structure

## Example Usage Using Dynamic Blocks

```terraform
data "zpa_inspection_all_predefined_controls" "this" {
  version    = "OWASP_CRS/3.3.0"
  group_name = "Preprocessors"
}

data "zpa_inspection_predefined_controls" "this" {
  name     = "Failed to parse request body"
  version  = "OWASP_CRS/3.3.0"
}

resource "zpa_inspection_profile" "this" {
  name                          = "Example"
  description                   = "Example"
  paranoia_level                = "1"
  predefined_controls_version   = "OWASP_CRS/3.3.0"
  incarnation_number            = "6"
  controls_info {
    control_type = "PREDEFINED"
  }
  dynamic "predefined_controls" {
    for_each = data.zpa_inspection_all_predefined_controls.default_predefined_controls.list
    content {
      id           = predefined_controls.value.id
      action       = predefined_controls.value.action == "" ? predefined_controls.value.default_action : predefined_controls.value.action
      action_value = predefined_controls.value.action_value
    }
  }
  predefined_controls {
    id     = data.zpa_inspection_predefined_controls.this.id
    action = "BLOCK"
  }
  global_control_actions = [
    "PREDEFINED:PASS",
    "CUSTOM:NONE",
    "OVERRIDE_ACTION:COMMON"
  ]
  common_global_override_actions_config = {
    "PREDEF_CNTRL_GLOBAL_ACTION" : "PASS",
    "IS_OVERRIDE_ACTION_COMMON" : "TRUE"
  }
}
```

## Example Usage Using Locals and Dynamic Blocks

```terraform
locals {
  group_names = {
    default_predefined_controls = "preprocessors" // Mandatory and must always be included
    group1 = "Protocol Issues"
    group2 = "Environment and port scanners"
    group3 = "Remote Code Execution"
    group4 = "Remote file inclusion"
    group5 = "Local File Inclusion"
    group6 = "Request smuggling or Response split or Header injection"
    group7 = "PHP Injection"
    group8 = "XSS"
    group9 = "SQL Injection"
    group10 = "Session Fixation"
    group11 = "Deserialization Issues"
    group12 = "Anomalies"
    group13 = "Request smuggling or Response split or Header injection"
  }
}

# Dynamically create data sources using for_each
data "zpa_inspection_all_predefined_controls" "all" {
  for_each = local.group_names
  group_name = each.value
}

# Combine the data source results into a single list
locals {
  combined_predefined_controls = flatten([
    for ds in data.zpa_inspection_all_predefined_controls.all : ds.list
  ])
}

resource "zpa_inspection_profile" "example" {
  name           = "Example"
  description    = "Example"
  paranoia_level = "2"
  incarnation_number = "6"
  controls_info {
    control_type = "PREDEFINED"
  }

  dynamic "predefined_controls" {
    for_each = local.combined_predefined_controls
    content {
      id     = predefined_controls.value.id
      action = predefined_controls.value.action == "" ? predefined_controls.value.default_action : predefined_controls.value.action
    }
  }
  global_control_actions = [
    "PREDEFINED:PASS",
    "CUSTOM:NONE",
    "OVERRIDE_ACTION:COMMON"
  ]
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (Required) The name of the inspection profile.
- `description` - (Optional) Description of the inspection profile.
- `paranoia_level` - (Required) OWASP Predefined Paranoia Level. Range: [1-4], inclusive
- `predefined_controls` - (Required) The predefined controls. The default predefined control `Preprocessors` is mandatory and must be always included in the request by default. Individual `predefined_controls` can be set by using the data source `data_source_zpa_predefined_controls` or by group using the data source `zpa_inspection_all_predefined_controls`.
  - `id` - (Required) ID of the predefined control
  - `action` - (Required) The action of the predefined control. Supported values: `PASS`, `BLOCK` and `REDIRECT`

-> **NOTE:** When assigning predefined controls by control group, use the data source `zpa_inspection_all_predefined_controls` with the following parameters:
- `group_name` = "preprocessors"

### Optional

- `attachment` (Optional) Control attachment
- `control_group` (Optional) Control group
- `associate_all_controls` (Optional) When set to `true`, `ALL` predefined controls are automatically associated with the profile. If set to `false`, `ALL` predefined controls are dissociated from the profile.

- `custom_controls` - (Optional) Types for custom controls
  - `type` (Optional) Types for custom controls
  - `control_rule_json` (Optional) Custom controls string in JSON format
  - `rules` - (Optional) Rules of the custom controls applied as conditions `JSON`
    - `conditions` - (Optional)
      - `lhs` - (Optional) Signifies the key for the object type Supported values: `SIZE`, `VALUE`
      - `op` - (Optional) If lhs is set to SIZE, then the user may pass one of the following: `EQ: Equals`, `LE: Less than or equal to`, `GE: Greater than or equal to`. If the lhs is set to `VALUE`, then the user may pass one of the following: `CONTAINS`, `STARTS_WITH`, `ENDS_WITH`, `RX`.
      - `rhs` - (Optional) Denotes the value for the given object type. Its value depends on the key. If rules.type is set to REQUEST_METHOD, the conditions.rhs field must have one of the following values: `GET`,`HEAD`, `POST`, `OPTIONS`, `PUT`, `DELETE`, `TRACE`
    - `names` - (Optional) Name of the rules. If rules.type is set to `REQUEST_HEADERS`, `REQUEST_COOKIES`, or `RESPONSE_HEADERS`, the rules.name field is required.
    - `type` - (Optional) Type value for the rules
    - `version` - (Optional) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`

- `associated_inspection_profile_names` - (Optional) Name of the inspection profile
  - `id`- (Optional)
  - `name`- (Optional)

- `common_global_override_actions_config` - (Optional)
- `controls_info` - (Optional) Types for custom controls
  - `control_type` - (string) Control types. Supported Values: `WEBSOCKET_PREDEFINED`, `WEBSOCKET_CUSTOM`, `CUSTOM`, `PREDEFINED`, `ZSCALER`
  - `count` - (Optional) Control information counts `Long`

- `web_socket_controls` - (string)
  - `id` - (string) ID of the predefined control
  - `action` - (string) The action of the predefined control. Supported values: `PASS`, `BLOCK` and `REDIRECT`
  - `action_value` - (string) Value for the predefined controls action. This field is only required if the action is set to REDIRECT. This field is only required if the action is set to `REDIRECT`.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)
