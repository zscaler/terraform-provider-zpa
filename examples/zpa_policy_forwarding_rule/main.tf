// Create Policy Forwarding Rule
resource "zpa_policy_forwarding_rule" "crm_application_rule" {
  name                          = "CRM Application"
  description                   = "CRM Application"
  action                        = "BYPASS"
  operator = "AND"
  policy_set_id = data.zpa_policy_type.client_forwarding_policy.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "APP"
      lhs = "id"
      rhs = [data.zpa_application_segment.crm_application.id]
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

// Retrieve Client Forwarding Policy ID
data "zpa_policy_type" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

// Retrieve IDP ID Information
data "zpa_idp_controller" "idp_name" {
 name = "IdP-Name"
}

// Retrieve SCIM Group ID
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IDP-Name"
}