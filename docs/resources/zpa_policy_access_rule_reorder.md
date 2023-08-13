---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_access_rule_reorder"
description: |-
  Creates and Updates rule orders in all ZPA Policy Access types.
---

# Resource: zpa_policy_access_rule_reorder

The **zpa_policy_access_rule_reorder** is a dedicated resource to manage and update rule_orders in any of the supported ZPA Policy Access types Zscaler Private Access cloud.

⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of this resource for all ZPA policy types.

## Example Usage 1

```hcl
data "zpa_policy_type" "access_policy" {
  policy_type = "ACCESS_POLICY"
}

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

```hcl
data "zpa_policy_type" "access_policy" {
  policy_type = "ACCESS_POLICY"
}

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

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``
* `policy_type` (Required) - Supported values:
  * ``ACCESS_POLICY or GLOBAL_POLICY``
  * ``TIMEOUT_POLICY or REAUTH_POLICY``
  * ``BYPASS_POLICY or CLIENT_FORWARDING_POLICY``
  * ``INSPECTION_POLICY``

## Attributes Reference

* `rules` - (Required)
  * `id` - (Required) - The ID of the rule to which the order number will be applied.
  * `order` (Required) - The order number that should be applied to the respective rule ID.
