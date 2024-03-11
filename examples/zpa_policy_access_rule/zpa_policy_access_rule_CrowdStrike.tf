// CrowdStrike_ZTA_Score_Policy
resource "zpa_policy_access_rule" "crwd_zta_score_40" {
  name                          = "CrowdStrike_ZTA_Score_40"
  description                   = "CrowdStrike_ZTA_Score_40"
  action                        = "DENY"
  operator = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id
  conditions {
    operator = "OR"
    operands {
      object_type = "POSTURE"
      lhs = data.zpa_posture_profile.crwd_zta_score_40.posture_udid
      rhs = false
    }
  }
  conditions {
     operator = "OR"
    operands {
      object_type = "SAML"
      lhs = data.zpa_saml_attribute.email_user_sso.id
      rhs_list = ["user1@acme.com", "user2@acme.com"]
      idp_id = data.zpa_idp_controller.user_idp_name.id
    }
  }
}

// Retrieve Policy Type ID
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

// Retrieve IDP ID information
data "zpa_idp_controller" "sgio_user_okta" {
 name = "User_IDP_Name"
}

// Retrieve SAML Attribute information
data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_User_IDP_Name"
}

// Retrieve Posture Profile UUID information
data "zpa_posture_profile" "crwd_zta_score_40" {
 name = "CrowdStrike_ZPA_ZTA_40"
}