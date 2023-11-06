# Get Log Type Format - "User Status"
data "zpa_lss_config_log_type_formats" "zpn_auth_log" {
  log_type = "zpn_auth_log"
}

data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

# Retrieve the App Connector Group ID
data "zpa_app_connector_group" "this" {
 name = "Example100"
}

# Retrieve the Identity Provider ID
data "zpa_idp_controller" "this" {
 name = "Idp_Name"
}

# Retrieve the SCIM_GROUP ID(s)
data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "Idp_Name"
}

data "zpa_scim_groups" "sales" {
  name     = "Sales"
  idp_name = "Idp_Name"
}

resource "zpa_lss_config_controller" "lss_user_activity" {
  config {
    name            = "LSS User Status"
    description     = "LSS User Status"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_auth_log.json
    lss_host        = "splunk1.acme.com"
    lss_port        = "5001"
    source_log_type = "zpn_auth_log"
    use_tls         = true
    filter = ["ZPN_STATUS_AUTH_FAILED","ZPN_STATUS_DISCONNECTED", "ZPN_STATUS_AUTHENTICATED"]
  }
  policy_rule_resource {
    name          = "policy_rule_resource_lss_user_status"
    action        = "LOG"
    policy_set_id = data.zpa_policy_type.lss_siem_policy.id
    conditions {
      negated  = false
      operator = "OR"
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_exporter", "zpn_client_type_browser_isolation", "zpn_client_type_machine_tunnel", "zpn_client_type_ip_anchoring", "zpn_client_type_edge_connector", "zpn_client_type_zapp", "zpn_client_type_slogger", "zpn_client_type_zapp_partner", "zpn_client_type_branch_connector"]
      }
    }
    conditions {
      negated  = false
      operator = "OR"
      operands {
        object_type = "SCIM_GROUP"
        entry_values {
          rhs = data.zpa_scim_groups.engineering.id
          lhs = data.zpa_idp_controller.this.id
        }
        entry_values {
          rhs = data.zpa_scim_groups.sales.id
          lhs = data.zpa_idp_controller.this.id
        }
      }
    }
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.this.id ]
  }
}