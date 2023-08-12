// Retrieve Timeout Policy Type ID
data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

// Retrieve IDP ID Information
data "zpa_idp_controller" "idp_name" {
 name = "IdP-Name"
}

// Retrieve SCIM Group Information
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP-Name"
}

// Create Policy Timeout Rule
resource "zpa_policy_timeout_rule" "crm_application_rule" {
  name                          = "CRM Application"
  description                   = "CRM Application"
  action                        = "RE_AUTH"
  reauth_idle_timeout           = "600"
  reauth_timeout                = "172800"
  operator                      = "AND"
  policy_set_id = data.zpa_global_policy_timeout.policyset.id

  conditions {
    negated   = false
    operator  = "OR"
    operands {
      object_type = "APP_GROUP"
      lhs = "id"
      rhs = [ data.zpa_segment_group.crm_application.id ]
    }
  }
  conditions {
     negated  = false
     operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [ data.zpa_scim_groups.engineering.id ]
    }
  }
}
