// Trusted Network Policy Rule
resource "zpa_policyset_rule" "corp_trusted_network" {
  name                          = "Corp Trusted Network"
  description                   = "Corp Trusted Network"
  action                        = "ALLOW"
  rule_order                    = 1
  operator = "AND"
  policy_set_id = data.zpa_policy_set_global.all.id
  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "TRUSTED_NETWORK"
      lhs = data.zpa_trusted_network.corp_trusted_network.network_id
      rhs = true
    }
  }
}

data "zpa_policy_set_global" "all" {}

data "zpa_trusted_network" "corp_trusted_network" {
 name = "Corp-Trusted-Networks"
}

