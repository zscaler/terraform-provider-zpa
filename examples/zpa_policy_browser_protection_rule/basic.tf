## Example Usage - Browser Protection with MONITOR ACTION

# Retrieve Policy Types
# Retrieve Identity Provider ID
data "zpa_idp_controller" "this" {
  name = "Idp_Name"
}

# Retrieve SAML Attribute ID
data "zpa_saml_attribute" "email_user_sso" {
  name     = "Email_Users"
  idp_name = "Idp_Name"
}

# Retrieve SAML Attribute ID
data "zpa_saml_attribute" "group_user" {
  name     = "GroupName_Users"
  idp_name = "Idp_Name"
}

# Retrieve SCIM Group ID
data "zpa_scim_groups" "a000" {
  name     = "A000"
  idp_name = "Idp_Name"
}

# Retrieve SCIM Group ID
data "zpa_scim_groups" "b000" {
  name     = "B000"
  idp_name = "Idp_Name"
}

# Create Segment Group
resource "zpa_segment_group" "this" {
  name        = "Example"
  description = "Example"
  enabled     = true
}

# Create Policy Access Rule V2
resource "zpa_policy_browser_protection_rule" "this" {
  name        = "Example"
  description = "Example"
  action      = "MONITOR"

  conditions {
    operator = "OR"
    operands {
      object_type = "APP_GROUP"
      values      = [zpa_segment_group.this.id]
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
      entry_values {
        rhs = "A000"
        lhs = data.zpa_saml_attribute.group_user.id
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
  conditions {
    operator = "OR"
    operands {
      object_type = "USER_PORTAL"
      values      = ["145262059234265326"]
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      values      = ["zpn_client_type_exporter"]
    }
  }
}


## Example Usage - Browser Protection with DO_NOT_MONITOR ACTION

# Retrieve Policy Types
# Retrieve Identity Provider ID
data "zpa_idp_controller" "this" {
  name = "Idp_Name"
}

# Retrieve SAML Attribute ID
data "zpa_saml_attribute" "email_user_sso" {
  name     = "Email_Users"
  idp_name = "Idp_Name"
}

# Retrieve SAML Attribute ID
data "zpa_saml_attribute" "group_user" {
  name     = "GroupName_Users"
  idp_name = "Idp_Name"
}

# Retrieve SCIM Group ID
data "zpa_scim_groups" "a000" {
  name     = "A000"
  idp_name = "Idp_Name"
}

# Retrieve SCIM Group ID
data "zpa_scim_groups" "b000" {
  name     = "B000"
  idp_name = "Idp_Name"
}

# Create Segment Group
resource "zpa_segment_group" "this" {
  name        = "Example"
  description = "Example"
  enabled     = true
}

# Create Policy Access Rule V2
resource "zpa_policy_browser_protection_rule" "this" {
  name        = "Example"
  description = "Example"
  action      = "DO_NOT_MONITOR"

  conditions {
    operator = "OR"
    operands {
      object_type = "APP_GROUP"
      values      = [zpa_segment_group.this.id]
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
      entry_values {
        rhs = "A000"
        lhs = data.zpa_saml_attribute.group_user.id
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
  conditions {
    operator = "OR"
    operands {
      object_type = "USER_PORTAL"
      values      = ["145262059234265326"]
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      values      = ["zpn_client_type_exporter"]
    }
  }
}
