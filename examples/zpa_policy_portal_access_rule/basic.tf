data "zpa_saml_attribute" "email_user_sso" {
  name     = "Email_Okta_Users"
  idp_name = "Okta_Users"
}

data "zpa_idp_controller" "this" {
  name = "Okta_Users"
}

data "zpa_scim_groups" "a000" {
  idp_name = "Okta_Users"
  name     = "A000"
}

data "zpa_scim_groups" "b000" {
  idp_name = "Okta_Users"
  name     = "B000"
}

resource "zpa_policy_portal_access_rule" "example" {
  name        = "Portal Access Rule"
  description = "Allow portal access with specific capabilities"
  action      = "CHECK_PRIVILEGED_PORTAL_CAPABILITIES"

  privileged_portal_capabilities {
    delete_file             = true
    access_uninspected_file = true
    request_approvals       = true
    review_approvals        = true
  }

  conditions {
    operator = "OR"
    operands {
      object_type = "PRIVILEGE_PORTAL"
      values      = ["216196257331387235"]
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "COUNTRY_CODE"
      entry_values {
        lhs = "BR"
        rhs = "true"
      }
      entry_values {
        lhs = "CA"
        rhs = "true"
      }
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "SAML"
      entry_values {
        rhs = "user1@acme.com"
        lhs = data.zpa_saml_attribute.email_user_sso.id
      }
    }
    operands {
      object_type = "SCIM_GROUP"
      entry_values {
        rhs = data.zpa_scim_groups.a000.id
        lhs = data.zpa_idp_controller.this.id
      }
      entry_values {
        rhs = data.zpa_scim_groups.b000.id
        lhs = data.zpa_idp_controller.this.id
      }
    }
  }
}
