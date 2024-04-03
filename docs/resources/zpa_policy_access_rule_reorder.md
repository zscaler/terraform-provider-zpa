---
page_title: "policy_access_rule_reorder Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-access-policy
  API documentation https://help.zscaler.com/zpa/configuring-access-policies-using-api
  Creates and Updates rule orders in all ZPA Policy Access types.
---

# policy_access_rule_reorder (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-access-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api)

The **zpa_policy_access_rule_reorder** is a dedicated resource to manage and update `rule_orders` in any of the supported ZPA Policy Access types Zscaler Private Access cloud.

⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of this resource for all ZPA policy types.

## Example Usage 1

```terraform
resource "zpa_policy_access_rule" "example001" {
  name          = "example001"
  description   = "example001"
  action        = "ALLOW"
  operator      = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id
}

resource "zpa_policy_access_rule" "example002" {
  name          = "example002"
  description   = "example002"
  action        = "ALLOW"
  operator      = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id
}

locals {
  rule_orders = [
    { id = zpa_policy_access_rule.example001.id, order = 1 },
    { id = zpa_policy_access_rule.example002.id, order = 2 },
  ]
}

resource "zpa_policy_access_rule_reorder" "access_policy_reorder" {
  policy_set_id = data.zpa_policy_type.access_policy.id
  policy_type   = "ACCESS_POLICY"

  dynamic "rules" {
    for_each = local.rule_orders
    content {
      id    = rules.value.id
      order = rules.value.order
    }
  }
}
```

## Example Usage 2

```terraform
resource "zpa_policy_access_rule" "example001" {
  name          = "example001"
  description   = "example001"
  action        = "ALLOW"
  operator      = "AND"
}

resource "zpa_policy_access_rule" "example002" {
  name          = "example002"
  description   = "example002"
  action        = "ALLOW"
  operator      = "AND"
}

locals {
  # Define a map with rule names as keys and their desired order as values.
  rule_order_map = {
    "example001"        = 1,
    "example002"        = 2,
  }

  rule_orders = [
    { id = zpa_policy_access_rule.example001.id, order = lookup(local.rule_order_map, "example001") },
    { id = zpa_policy_access_rule.example002.id, order = lookup(local.rule_order_map, "example002") },
  ]
}

resource "zpa_policy_access_rule_reorder" "access_policy_reorder" {
  policy_type   = "ACCESS_POLICY"

  dynamic "rules" {
    for_each = local.rule_orders
    content {
      id    = rules.value.id
      order = rules.value.order
    }
  }
}
```

## Schema

### Required

- `name` (String) This is the name of the policy rule.
- `policy_type` (String) - Supported values:
  - ``ACCESS_POLICY or GLOBAL_POLICY``
  - ``TIMEOUT_POLICY or REAUTH_POLICY``
  - ``BYPASS_POLICY or CLIENT_FORWARDING_POLICY``
  - ``INSPECTION_POLICY``
  - ``ISOLATION_POLICY``
  - ``CREDENTIAL_POLICY``
  - ``CAPABILITIES_POLICY``
  - ``CLIENTLESS_SESSION_PROTECTION_POLICY``

- `rules` - (Block Set)
  - `id` - (String) - The ID of the rule to which the order number will be applied.
  - `order` (String) - The order number that should be applied to the respective rule ID.
