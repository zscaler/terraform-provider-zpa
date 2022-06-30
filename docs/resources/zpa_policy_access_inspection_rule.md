---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_inspection_rule"
description: |-
  Creates and manages ZPA Policy Access Inspection Rule.
---

# Resource: zpa_policy_inspection_rule

The **zpa_policy_inspection_rule** resource creates a policy inspection access rule in the Zscaler Private Access cloud.

## Example Usage 1

```hcl
# Retrieve Inspection policy type
data "zpa_policy_type" "inspection_policy" {
  policy_type = "INSPECTION_POLICY"
}

#Create Inspection Access Rule
resource "zpa_policy_inspection_rule" "this" {
  name                      = "Example"
  description               = "Example"
  action                    = "INSPECT"
  rule_order                = 1
  operator                  = "AND"
  policy_set_id             = data.zpa_policy_type.inspection_policy.id
  zpn_inspection_profile_id = zpa_inspection_profile.this.id
  conditions {
    operator = "OR"
    operands {
      object_type = "APP"
      lhs         = "id"
      rhs         = zpa_application_segment_inspection.this.id
    }
  }
}
```

## Example Usage 2

```hcl
# Retrieve Inspection policy type
data "zpa_policy_type" "inspection_policy" {
  policy_type = "INSPECTION_POLICY"
}

#Create Inspection Access Rule
resource "zpa_policy_inspection_rule" "this" {
  name                      = "Example"
  description               = "Example"
  action                    = "BYPASS_INSPECT"
  rule_order                = 1
  operator                  = "AND"
  policy_set_id             = data.zpa_policy_type.inspection_policy.id
  conditions {
    operator = "OR"
    operands {
      object_type = "APP"
      lhs         = "id"
      rhs         = zpa_application_segment_inspection.this.id
    }
  }
}
```

### Required

* `name` - (Required) This is the name of the policy inspection rule.
* `policy_set_id` - (Required)

## Attributes Reference

* `action` - (Optional) This is for providing the rule action.
  * The supported actions for a policy inspection rule are: `BYPASS_INSPECT`, or `INSPECT`
* `zpn_inspection_profile_id` (Optional) An inspection profile is required if the `action` is set to `INSPECT`
* `action_id` - (Optional) This field defines the description of the server.
* `bypass_default_rule` - (Optional)
* `custom_msg` - (Optional) This is for providing a customer message for the user.
* `description` - (Optional) This is the description of the access policy rule.
* `operator` (Optional)
* `policy_type` - (Optional)
  * The supported policy type values for a policy inspection rule is: `INSPECTION_POLICY`
* `rule_order` - (Optional)

* `conditions` - (Optional)
  * `negated` - (Optional)
  * `operator` (Optional)
  * `operands`
    * `name` (Optional)
    * `lhs` (Optional)
    * `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.
    * `idp_id` (Optional)
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `CLOUD_CONNECTOR_GROUP`. `TRUSTED_NETWORK`, and `CLIENT_TYPE`.
    * `CLIENT_TYPE` (Optional) - The below options are the only ones supported in a timeout policy rule.
      * `zpn_client_type_exporter`
      * `zpn_client_type_browser_isolation`
      * `zpn_client_type_machine_tunnel`
      * `zpn_client_type_ip_anchoring`
      * `zpn_client_type_edge_connector`
      * `zpn_client_type_zapp`

## Import

Policy Access Inspection Rule can be imported by using `<POLICY INSPECTION RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_inspection_rule.example <policy_inspection_rule_id>
```
