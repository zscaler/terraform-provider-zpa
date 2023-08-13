data "zpa_policy_type" "access_policy" {
  policy_type = "ACCESS_POLICY"
}

data "zpa_scim_attribute_header" "givenName" {
  name     = "name.givenName"
  idp_name = "IdP_Name"
}

data "zpa_scim_attribute_header" "familyName" {
  name     = "name.familyName"
  idp_name = "IdP_Name"
}

data "zpa_posture_profile" "crwd_zpa_pre_zta" {
  name = "CrowdStrike_ZPA_Pre-ZTA"
}

// CrowdStrike_ZTA_Score_Policy
resource "zpa_policy_access_rule" "crwd_zpa_pre_zta" {
  name          = "CrowdStrike_ZPA_Pre-ZTA"
  description   = "CrowdStrike_ZPA_Pre-ZTA"
  action        = "DENY"
  operator      = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id
  conditions {
    negated  = false
    operator = "OR"
    operands {
      object_type = "POSTURE"
      lhs         = data.zpa_posture_profile.crwd_zpa_pre_zta.posture_udid
      rhs         = false
    }
  }
  conditions {
    negated  = false
    operator = "OR"
    operands {
      object_type = "SCIM"
      idp_id      = data.zpa_scim_attribute_header.givenName.idp_id
      lhs         = data.zpa_scim_attribute_header.givenName.id
      rhs         = "John"
    }
    operands {
      object_type = "SCIM"
      idp_id      = data.zpa_scim_attribute_header.familyName.idp_id
      lhs         = data.zpa_scim_attribute_header.familyName.id
      rhs         = "Smith"
    }
  }
}
