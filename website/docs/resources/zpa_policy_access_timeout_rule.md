---
subcategory: "Policy Timeout Rule"
layout: "zpa"
page_title: "ZPA: policy_timeout"
description: |-
  Creates a ZPA Policy Timeout Rule.
---
# zpa_policy_timeout_rule (Resource)

The **zpa_policy_timeout_rule** resource creates a policy timeout rule in the Zscaler Private Access cloud.

## Example Usage

```hcl
resource "zpa_policy_timeout_rule" "example_timeout_access_rule" {
  name                          = "Example Timeout Access Rule"
  description                   = "Example Timeout Access Rule"
  action                        = "RE_AUTH"
  reauth_idle_timeout           = "600"
  reauth_timeout                = "172800"
  operator                      = "AND"
  policy_set_id                 = data.zpa_policy_type.timeout_policy.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs = "id"
      rhs = "zpn_client_type_exporter"
    }
  }
  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [data.zpa_scim_groups.engineering.id]
    }
  }
}
```

```hcl
data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}
```

```hcl
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}
```

```hcl
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP_Name"
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required)

## Attributes Reference

* `action` (Optional) This is for providing the rule action.
* `action_id` (Optional) This field defines the description of the server.
* `bypass_default_rule` (Boolean)
* `custom_msg` (Optional) This is for providing a customer message for the user.
* `description` (Optional) This is the description of the access policy rule.
* `operator` (Optional)
* `policy_type` (Optional)
* `priority` (Optional)
* `reauth_default_rule` (Optional)
* `reauth_idle_timeout` (Optional)
* `reauth_timeout` (Optional)
* `rule_order` (Optional)

`conditions` - (Optional)

* `negated` (Optional)
* `idp_id` (Optional)
* `operator` (Optional)
* `name` (Optional)
* `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `CLOUD_CONNECTOR_GROUP`. TRUSTED_NETWORK is only supported for CLIENT_TYPE

`operands`

* `lhs` (Optional)
* `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.

`app_connector_groups`

* `id` - (Optional) The ID of this resource.

`app_server_groups`

* `id` - (Optional) The ID of this resource.

## Import

Policy access timeout can be imported; use `<POLICY Access RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_timeout_rule.example 216196257331290863
```
