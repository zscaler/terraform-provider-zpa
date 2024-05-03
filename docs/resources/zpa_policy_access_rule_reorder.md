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

⚠️ **WARNING:**: Updating the rule order of an access policy configured using `Zscaler Deception` is not supported. When changing the rule order of a regular access policy and there is an access policy configured using Deception, the rule order of the regular access policy must be greater than the rule order for an access policy configured using Deception. Please refer to the [Zscaler API Documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api#:~:text=Updating%20the%20rule,configured%20using%20Deception.) for further details.

## Example Usage 1

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
  rule_orders = [
    { id = zpa_policy_access_rule.example001.id, order = 1 },
    { id = zpa_policy_access_rule.example002.id, order = 2 },
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

## Example Usage 3 - Used when Zscaler Deception Rule Exists

```terraform
# IF NO ZSCALER DECEPTION RULE EXIST, DECREASE THE INDEX TO +1 TO PREVENT DRIFTS
locals {
  policy_config = yamldecode(file("${path.module}/policies.yaml"))
  policies = { for policy in local.policy_config.policies : policy.name => merge(policy, { rule_number = index(local.policy_config.policies, policy) + 2 }) }
}

resource "zpa_policy_access_rule" "rules" {
  for_each      = local.policies
  name          = each.value.name
  action        = each.value.action
  description   = each.value.description
  custom_msg    = try(each.value.custom_msg, null)
  operator      = try(each.value.operator, "AND")
}

resource "zpa_policy_access_rule_reorder" "access_policy_reorder" {
  policy_type = "ACCESS_POLICY"

  dynamic "rules" {
    for_each = local.policies
    content {
      id    = zpa_policy_access_rule.rules[rules.key].id
      order = rules.value.rule_number
    }
  }
}
```

## Example Usage 4 - Similar to Example 3 - No YAML File

```terraform
locals {
  policies = { for index, policy in var.policy_config.policies :
    policy.name => merge(policy, { rule_number = index + 1 })
  }
}

resource "zpa_policy_access_rule" "rules" {
  for_each      = { for rule in local.policies : rule.name => rule }
  name          = each.value.name
  action        = each.value.action
  description   = each.value.description
  custom_msg    = try(each.value.custom_msg, null)
  operator      = try(each.value.operator, "AND")
}


resource "zpa_policy_access_rule_reorder" "access_policy_reorder" {
  policy_type = "ACCESS_POLICY"

  dynamic "rules" {
    for_each = local.policies  # This sets up 'rules' as the variable within the block
    content {
      id    = zpa_policy_access_rule.rules[rules.key].id  # Access 'rules.key' for the map key
      order = rules.value.rule_number  # Use 'rules.value' to get the values from the map
    }
  }
}

variable "policy_config" {
  description = "Configuration for policy rules"
  type = object({
    policies = list(object({
      name        = string
      description = string
      action = string
      // Additional attributes can be included here as needed
    }))
  })

  default = {
    policies = [
      { name = "example001", description = "example001", action = "ALLOW"},
      { name = "example002", description = "example002", action = "DENY" },
      { name = "example003", description = "example003", action = "ALLOW" },
      { name = "example004", description = "example004", action = "DENY" },
      { name = "example005", description = "example005", action = "ALLOW" },
      { name = "example006", description = "example006", action = "DENY" },
      { name = "example007", description = "example007", action = "ALLOW" },
      { name = "example008", description = "example008", action = "DENY" },
      { name = "example009", description = "example009", action = "ALLOW" },
      { name = "example010", description = "example010", action = "DENY" },
    ]
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
