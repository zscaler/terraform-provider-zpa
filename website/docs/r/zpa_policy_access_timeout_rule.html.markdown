---
layout: "zscaler"
page_title: "Zscaler Private Access (ZPA): policy_timeout_rule"
sidebar_current: "docs-resource-zpa-policy-timeout-rule"
description: |-
  Creates and manages ZPA Policy Timeout Access Rule.
---

# zpa_policy_timeout_rule (Resource)

The **zpa_policy_timeout_rule** resource creates a policy timeout rule in the Zscaler Private Access cloud.

## Example Usage

```hcl
resource "zpa_policy_timeout_rule" "test_policy_timeout"  {
  name                          = "test1-policy-timeout-rule"
  description                   = "test1-policy-timeout-rule"
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
      object_type = "CLIENT_TYPE"
      lhs = "id"
      rhs = "zpn_client_type_zapp"
    }
  }
  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs = "id"
      rhs = "zpn_client_type_browser_isolation"
    }
  }
  conditions {
    negated  = false
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [data.zpa_scim_groups.engineering.id]
    }
  }
  depends_on = [
  data.zpa_policy_type.access_policy,
  data.zpa_idp_controller.idp_name,
  data.zpa_scim_groups.engineering,
  ]
}

# Get Global Timeout Policy ID
data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

# Get IdP ID
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

# Get SCIM Group attribute ID
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP_Name"
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required)
* `reauth_default_rule` (Required)
* `reauth_idle_timeout` (Required)

## Attributes Reference

* `action` (Optional) This is for providing the rule action.
* `bypass_default_rule` (Optional)
* `custom_msg` (Optional) This is for providing a customer message for the user.
* `description` (Optional) This is the description of the access policy rule.
* `operator` (Optional)
* `policy_type` (Optional)
* `reauth_default_rule` (Required)
* `reauth_idle_timeout` (Required)
* `reauth_timeout` (Optional)
* `rule_order` (Optional)

* `conditions` - (Optional)
  * `negated` - (Optional)
  * `operator` (Optional)
  * `operands`
    * `name` (Optional)
    * `lhs` (Optional)
    * `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.
    * `idp_id` (Optional)
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `SAML`, `SCIM`, `SCIM_GROUP`, `IDP`, `CLIENT_TYPE`,  `POSTURE`
    * `CLIENT_TYPE` (Optional) - The below options are the only ones supported in a timeout policy rule.
      * `zpn_client_type_zapp`
      * `zpn_client_type_browser_isolation`
      * `zpn_client_type_exporter`

## Import

Policy access timeout can be imported by using `<POLICY TIMEOUT RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_timeout_rule.example <policy_timeout_rule_id>
```
