data "zpa_policy_type" "this" {
  policy_type = "ACCESS_POLICY"
}

resource "zpa_policy_access_rule" "this" {
  name          = "Example"
  description   = "Example"
  action        = "ALLOW"
  operator      = "AND"
  policy_set_id = data.zpa_policy_type.this.id
}

resource "zpa_policy_access_rule_reorder" "this" {
  policy_type = "ACCESS_POLICY"
  rules = {
    id = zpa_policy_access_rule.this.id
  }
}