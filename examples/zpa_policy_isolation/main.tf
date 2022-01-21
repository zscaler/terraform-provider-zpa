terraform {
    required_providers {
        zpa = {
            version = "2.0.6"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

// data "zpa_global_policy_isolation" "all" {
// }

// output "isolation_policy" {
//     value = data.zpa_global_policy_isolation.all
// }

// resource "zpa_policy_isolation_rule" "isolation_bypass_rule" {
//   name                          = "Isolation Rule"
//   description                   = "Isolation Rule"
//   action                        = "ISOLATE"
//   zpn_cbi_profile_id             = "216196257331286656"
//   rule_order                     = 1
//   operator = "AND"
//   policy_set_id = data.zpa_global_policy_isolation.all.id

//   conditions {
//     negated = false
//     operator = "OR"
//     operands {
//       object_type = "CLIENT_TYPE"
//       lhs = "id"
//       rhs = "zpn_client_type_exporter"
//     }
//   }
// }

// output "zpa_policy_isolation_rule" {
//     value = zpa_policy_isolation_rule.isolation_bypass_rule
// }

// data "zpa_global_policy_isolation" "all" {
// }

// data "zpa_global_access_policy" "access_policy" {
//     policy_type = "ACCESS_POLICY"
// }

data "zpa_global_policy_timeout" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

// data "zpa_global_policy_forwarding" "client_forwarding_policy" {
//     policy_type = "CLIENT_FORWARDING_POLICY"
// }

resource "zpa_policy_timeout_rule" "crm_application_rule" {
  name                          = "CRM Application"
  description                   = "CRM Application"
  action                        = "NO_DOWNLOAD"
  operator = "AND"
  policy_set_id = data.zpa_global_policy_timeout.timeout_policy.id
}

output "zpa_policy_timeout_rule" {
  value = zpa_policy_timeout_rule.crm_application_rule
}