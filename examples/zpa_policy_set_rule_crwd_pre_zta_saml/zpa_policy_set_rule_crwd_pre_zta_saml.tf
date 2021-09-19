// CrowdStrike_ZTA_Score_Policy
resource "zpa_policyset_rule" "crwd_zpa_pre_zta" {
  name                          = "CrowdStrike_ZPA_Pre-ZTA"
  description                   = "CrowdStrike_ZPA_Pre-ZTA"
  action                        = "DENY"
  rule_order                    = 1
  operator = "AND"
  policy_set_id = data.zpa_policy_set_global.all.id
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
      object_type = "SAML"
      lhs = data.zpa_saml_attribute.email_user_sso.id
      rhs = "user1@acme.com"
      idp_id = data.zpa_idp_controller.idp_name.id
    }
  }
}

data "zpa_policy_set_global" "all" {}

data "zpa_posture_profile" "crwd_zpa_pre_zta" {
 name = "CrowdStrike_ZPA_Pre-ZTA"
}

data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_SGIO-User-Okta"
}