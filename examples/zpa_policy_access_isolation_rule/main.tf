data "zpa_policy_type" "isolation_policy" {
    policy_type = "ISOLATION_POLICY"
}

data "zpa_isolation_profile" "isolation_profile" {
    name = "zpa_isolation_profile"
}

resource "zpa_policy_isolation_rule" "this" {
  name                          = "Example_Isolation_Policy"
  description                   = "Example_Isolation_Policy"
  action                        = "ISOLATE"
  operator = "AND"
  policy_set_id = data.zpa_policy_type.isolation_policy.id
  zpn_isolation_profile_id = data.zpa_isolation_profile.isolation_profile.id

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs = "id"
      rhs = "zpn_client_type_exporter"
    }
  }
}