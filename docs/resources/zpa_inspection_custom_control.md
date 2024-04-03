---
page_title: "zpa_inspection_custom_controls Resource - terraform-provider-zpa"
subcategory: "AppProtection"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-custom-controls/API documentation https://help.zscaler.com/zpa/configuring-appprotection-controls-using-api
  Creates and manages Inspection Custom Control in Zscaler Private Access cloud.
---

# zpa_inspection_custom_controls (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-custom-controls)
* [API documentation](https://help.zscaler.com/zpa/configuring-appprotection-controls-using-api)

The **zpa_inspection_custom_controls** resource creates an inspection custom control. This resource can then be referenced in an inspection profile resource.

## Example Usage

```terraform
data "zpa_inspection_profile" "this" {
  name = "Example"
}

resource "zpa_inspection_custom_controls" "this" {
  name           = "Example"
  description    = "Example"
  action         = "PASS"
  default_action = "PASS"
  paranoia_level = "1"
  severity       = "CRITICAL"
  type = "RESPONSE"
  associated_inspection_profile_names {
    id = [data.zpa_inspection_profile.this.id]
  }
  rules {
    names = ["this"]
    type  = "RESPONSE_HEADERS"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
  rules {
    type  = "RESPONSE_BODY"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (Required) The name of the predefined control.
- `version` - (Required) The version of the predefined control, the default is: `OWASP_CRS/3.3.0`
- `action` - (Required) The performed action. Supported values: `PASS`, `BLOCK` and `REDIRECT`
- `action_value` - (Required) Denotes the action
- `name` - (Required) Name of the custom control
- `paranoia_level` - (Required) OWASP Predefined Paranoia Level.
- `protocol_type` - (string) Returned values: `HTTP`, `HTTPS`, `FTP`, `RDP`, `SSH`, `WEBSOCKET`
- `severity` - (Required) Severity of the control number. Supported values: `CRITICAL`, `ERROR`, `WARNING`, `INFO`
- `type` - (Required) Rules to be applied to the request or response type
- `rules` - (Required) Rules of the custom controls applied as conditions `JSON`
  - `conditions` - (Required)
    - `lhs` - (Required) Signifies the key for the object type Supported values: `SIZE`, `VALUE`
    - `op` - (Required) If lhs is set to SIZE, then the user may pass one of the following: `EQ: Equals`, `LE: Less than or equal to`, `GE: Greater than or equal to`. If the lhs is set to `VALUE`, then the user may pass one of the following: `CONTAINS`, `STARTS_WITH`, `ENDS_WITH`, `RX`.
    - `rhs` - (Required) Denotes the value for the given object type. Its value depends on the key. If rules.type is set to REQUEST_METHOD, the conditions.rhs field must have one of the following values: `GET`,`HEAD`, `POST`, `OPTIONS`, `PUT`, `DELETE`, `TRACE`
  - `names` - (Required) Name of the rules. If rules.type is set to `REQUEST_HEADERS`, `REQUEST_COOKIES`, or `RESPONSE_HEADERS`, the rules.name field is required.
  - `type` - (Required) Type value for the rules

### Optional

- `description` - (Optional) Description of the custom control
- `associated_inspection_profile_names` - (Optional) Name of the inspection profile
  - `id`- (Optional)
  - `name`- (Optional)
- `control_rule_json` (Optional) The control rule in JSON format that has the conditions and type of control for the inspection control
- `control_type` - (string) Returned values: `WEBSOCKET_PREDEFINED`, `WEBSOCKET_CUSTOM`, `ZSCALER`, `CUSTOM`, `PREDEFINED`
- `default_action` - (Required) The performed action. Supported values: `PASS`, `BLOCK` and `REDIRECT`
- `default_action_value` - (Optional) This is used to provide the redirect URL if the default action is set to `REDIRECT`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)
