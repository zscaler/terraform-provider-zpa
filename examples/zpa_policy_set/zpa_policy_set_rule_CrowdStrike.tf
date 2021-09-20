terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

// CrowdStrike_ZTA_Score_Policy
resource "zpa_policy_access_rule" "crwd_zta_score_40" {
  name                          = "CrowdStrike_ZTA_Score_40"
  description                   = "CrowdStrike_ZTA_Score_40"
  action                        = "DENY"
  rule_order                     = 2
  operator = "AND"
  policy_set_id = data.zpa_global_access_policy.all.id
  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "POSTURE"
      lhs = data.zpa_posture_profile.crwd_zta_score_40.posture_udid
      rhs = false
    }
  }
  conditions {
     negated = false
     operator = "OR"
    operands {
      object_type = "SAML"
      lhs = data.zpa_saml_attribute.email_user_sso.id
      rhs_list = ["wguilherme@securitygeek.io", "wguilherme2@securitygeek.io"]
      idp_id = data.zpa_idp_controller.sgio_user_okta.id
    }
  }
}

data "zpa_global_access_policy" "all" {}

data "zpa_idp_controller" "sgio_user_okta" {
 name = "SGIO-User-Okta"
}

data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_SGIO-User-Okta"
}

data "zpa_posture_profile" "crwd_zta_score_40" {
 name = "CrowdStrike_ZPA_ZTA_40"
}