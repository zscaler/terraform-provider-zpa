terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_global_policy_isolation" "all" {
}

output "isolation_policy" {
    value = data.zpa_global_policy_isolation.all
}

resource "zpa_policy_isolation_rule" "isolation_bypass_rule" {
  name                          = "Isolation Bypass Rule"
  description                   = "Isolation Bypass Rule"
  action                        = "ISOLATE"
  rule_order                     = 1
  operator = "AND"
  policy_set_id = data.zpa_global_policy_isolation.all.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs = "id"
      rhs = "zpn_client_type_exporter"
    }
  }
}

output "zpa_policy_isolation_rule" {
    value = zpa_policy_isolation_rule.isolation_bypass_rule
}
