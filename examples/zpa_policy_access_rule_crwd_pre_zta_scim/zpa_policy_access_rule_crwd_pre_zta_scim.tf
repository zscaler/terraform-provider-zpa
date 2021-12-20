// CrowdStrike_ZTA_Score_Policy
resource "zpa_policy_access_rule" "crwd_zpa_pre_zta" {
  name                          = "CrowdStrike_ZPA_Pre-ZTA"
  description                   = "CrowdStrike_ZPA_Pre-ZTA"
  action                        = "DENY"
  rule_order                    = 1
  operator = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id
  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "POSTURE"
      lhs = data.zpa_posture_profile.crwd_zpa_pre_zta.posture_udid
      rhs = false
    }
  }
  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [data.zpa_scim_groups.engineering.id]
    }
  }
}

data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP_Name"
}

data "zpa_posture_profile" "crwd_zpa_pre_zta" {
 name = "CrowdStrike_ZPA_Pre-ZTA"
}