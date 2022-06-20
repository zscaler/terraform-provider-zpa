---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_access_rule"
description: |-
  Creates and manages ZPA Policy Access Rule with SCIM Attribute Header conditions.
---

# Resource: zpa_policy_access_rule

The **zpa_policy_access_rule** resource creates and manages a policy access rule with SCIM attribute header conditions in the Zscaler Private Access cloud.

## Example Usage

```hcl
data "zpa_policy_type" "access_policy" {
  policy_type = "ACCESS_POLICY"
}

data "zpa_scim_attribute_header" "givenName" {
  name     = "name.givenName"
  idp_name = "IdP_Name"
}

data "zpa_scim_attribute_header" "familyName" {
  name     = "name.familyName"
  idp_name = "IdP_Name"
}

resource "zpa_policy_access_rule" "browser_access_rule" {
  name                          = "test1-ba-policy-access"
  description                   = "test1-ba-policy-access"
  action                        = "ALLOW"
  operator                      = "AND"
  policy_set_id                 = data.zpa_policy_type.access_policy.id

  conditions {
    negated  = false
    operator = "OR"
    operands {
      object_type = "SCIM"
      idp_id      = data.zpa_scim_attribute_header.givenName.idp_id
      lhs         = data.zpa_scim_attribute_header.givenName.id
      rhs         = "John"
    }
    operands {
      object_type = "SCIM"
      idp_id      = data.zpa_scim_attribute_header.familyName.idp_id
      lhs         = data.zpa_scim_attribute_header.familyName.id
      rhs         = "Smith"
    }
  }
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required)

## Attributes Reference

* `action` (String) This is for providing the rule action.
* `action_id` (String) This field defines the description of the server.
* `bypass_default_rule` (Boolean)
* `custom_msg` (String) This is for providing a customer message for the user.
* `description` (String) This is the description of the access policy rule.
* `operator` (String)
* `policy_type` (String)
* `priority` (String)
* `reauth_default_rule` (Boolean)
* `reauth_idle_timeout` (String)
* `reauth_timeout` (String)
* `rule_order` (String)

* `conditions` - (Optional)
  * `negated` - (Optional)
  * `operator` (Optional)
  * `operands`
    * `name` (Optional)
    * `lhs` (Optional)
    * `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.
    * `idp_id` (Optional)
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `CLOUD_CONNECTOR_GROUP`. `TRUSTED_NETWORK`, and `CLIENT_TYPE`.
    * `CLIENT_TYPE` (Optional) - The below options are the only ones supported in an access policy rule.
      * `zpn_client_type_exporter`
      * `zpn_client_type_browser_isolation`
      * `zpn_client_type_machine_tunnel`
      * `zpn_client_type_ip_anchoring`
      * `zpn_client_type_edge_connector`
      * `zpn_client_type_zapp`

* `app_connector_groups`
  * `id` - (Optional) The ID of an app connector group resource

* `app_server_groups`
  * `id` - (Optional) The ID of a server group resource

## Import

Policy Access Rule for Browser Access can be imported by using`<POLICY ACCESS RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example <policy_access_rule_id>
```
