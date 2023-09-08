data "zpa_policy_type" "access_policy" {
  policy_type = "ACCESS_POLICY"
}


locals {
  access_rules = [
    { name = "example001", description = "example001", order = 1 },
    { name = "example002", description = "example002", order = 2 },
    { name = "example003", description = "example003", order = 3 },
  ]
}

resource "zpa_policy_access_rule" "this" {
  for_each      = { for rule in local.access_rules : rule.name => rule }
  name          = each.value.name
  description   = each.value.description
  action        = "ALLOW"
  operator      = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id
}

resource "zpa_policy_access_rule_reorder" "this" {
  policy_type = "ACCESS_POLICY"

  dynamic "rules" {
    for_each = [for rule in local.access_rules : { id = zpa_policy_access_rule.example[rule.name].id, order = rule.order }]
    content {
      id    = rules.value.id
      order = rules.value.order
    }
  }
}
