
// Retrieve Policy Types
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

// Retrieve Trusted Network NetworkID Information
data "zpa_trusted_network" "corp_trusted_network" {
 name = "Corp-Trusted-Networks"
}


// Trusted Network Policy Rule
resource "zpa_policy_access_rule" "corp_trusted_network" {
  name                          = "Corp Trusted Network"
  description                   = "Corp Trusted Network"
  action                        = "ALLOW"
  rule_order                    = 1
  operator = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id
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
