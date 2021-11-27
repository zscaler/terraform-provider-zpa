---
subcategory: "Policy Access Rule"
page_title: "ZPA: policyset_rule"
description: |-
  Creates a ZPA Policy Access Rule.
---

# zpa_policy_access_rule (Resource)

The **zpa_policy_access_rule** resource creates a policy access rule in the Zscaler Private Access cloud.

## Example Usage

```hcl
resource "zpa_policy_access_rule" "gf_engineering" {
  name                          = "GF-Engineering"
  description                   = "GF-Engineering"
  action                        = "ALLOW"
  operator                      = "AND"
  policy_set_id                 = data.zpa_policy_type.access_policy.id

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
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
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

`conditions` - (Optional)

* `negated` (Optional)
* `idp_id` (Optional)
* `operator` (String)
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

Policy access rule can be imported by using `<POLICY ACCESS RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example <policy_access_rule_id>
```
