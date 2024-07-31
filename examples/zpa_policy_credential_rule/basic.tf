resource "zpa_segment_group" "this" {
  name        = "Example_Segment_Group"
  description = "Example_Segment_Group"
  enabled     = true
}

resource "zpa_application_segment_pra" "this" {
  name             = "ZPA_PRA_Example"
  description      = "ZPA_PRA_Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_range = [
    {
    from = "3389"
    to = "3389"
    },
    {
    from = "22"
    to = "22"
    }
  ]
  domain_names     = ["ssh_pra.example.com", "rdp_pra.example.com"]
  segment_group_id = zpa_segment_group.this.id
  common_apps_dto {
    apps_config {
      name                 = "rdp_pra.example.com"
      domain               = "rdp_pra.example.com"
      application_protocol = "RDP"
      connection_security  = "ANY"
      application_port     = "3389"
      enabled              = true
      app_types            = ["SECURE_REMOTE_ACCESS"]
    }
    apps_config {
      name                 = "ssh_pra.example.com"
      domain               = "ssh_pra.example.com"
      application_protocol = "SSH"
      application_port     = "22"
      enabled              = true
      app_types            = ["SECURE_REMOTE_ACCESS"]
    }
  }
}

data "zpa_application_segment_by_type" "ssh_pra" {
  application_type = "SECURE_REMOTE_ACCESS"
  name             = "ssh_pra"
  depends_on       = [zpa_application_segment_pra.this]
}

data "zpa_application_segment_by_type" "rdp_pra" {
  application_type = "SECURE_REMOTE_ACCESS"
  name             = "rdp_pra"
  depends_on       = [zpa_application_segment_pra.this]
}

data "zpa_ba_certificate" "this" {
  name = "pra01.bd-hashicorp.com"
}

data "zpa_idp_controller" "this" {
  name = "BD_Okta_Users"
}

data "zpa_scim_groups" "a000" {
  name     = "A000"
  idp_name = "BD_Okta_Users"
}

data "zpa_saml_attribute" "email_user_sso" {
  name     = "Email_BD_Okta_Users"
  idp_name = "BD_Okta_Users"
}

resource "zpa_pra_portal_controller" "this" {
  name                      = "pra01.bd-hashicorp.com"
  description               = "pra01.bd-hashicorp.com"
  enabled                   = true
  domain                    = "pra01.bd-hashicorp.com"
  certificate_id            = data.zpa_ba_certificate.this.id
  user_notification         = "Created with Terraform"
  user_notification_enabled = true
}

resource "zpa_pra_console_controller" "rdp_pra" {
  name        = "RDP_PRA_Console"
  description = "Created with Terraform"
  enabled     = true
  pra_application {
    id = data.zpa_application_segment_by_type.rdp_pra.id
  }
  pra_portals {
    id = [zpa_pra_portal_controller.this.id]
  }
}

resource "zpa_pra_console_controller" "ssh_pra" {
  name        = "SSH_PRA_Console"
  description = "Created with Terraform"
  enabled     = true
  pra_application {
    id = data.zpa_application_segment_by_type.ssh_pra.id
  }
  pra_portals {
    id = [zpa_pra_portal_controller.this.id]
  }
}

resource "zpa_pra_credential_controller" "this" {
    name = "John Carrow"
    description = "Created with Terraform"
    credential_type = "USERNAME_PASSWORD"
    user_domain = "acme.com"
    username = "jcarrow"
    password = "************"
}

resource "zpa_policy_credential_rule" "this" {
  name          = "Example_Credential_Rule"
  description   = "Example_Credential_Rule"
  action        = "INJECT_CREDENTIALS"

  credential {
    id = zpa_pra_credential_controller.this.id
  }

  conditions {
    operator = "OR"
    operands {
      object_type = "CONSOLE"
      values      = [ zpa_pra_console_controller.rdp_pra.id, zpa_pra_console_controller.ssh_pra.id ]
    }
  }

  conditions {
    operator = "OR"
    operands {
      object_type = "SAML"
      entry_values {
        rhs = "jcarrow@acme.com"
        lhs = data.zpa_saml_attribute.email_user_sso.id
      }
    }
    operands {
      object_type = "SCIM_GROUP"
      entry_values {
        rhs = data.zpa_scim_groups.a000.id
        lhs = data.zpa_idp_controller.this.id
      }
    }
  }
}