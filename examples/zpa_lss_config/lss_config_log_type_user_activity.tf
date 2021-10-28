// Create Log Receiver Configuration
resource "zpa_lss_config_controller" "example" {
  config {
    name        = "Example"
    description = "Example"
    enabled     = true
    format      = data.zpa_lss_config_log_type_formats.zpn_trans_log.json
    lss_host    = "192.168.1.1"
    lss_port    = "11001"
    source_log_type = "zpn_trans_log"
    use_tls         = true
  }
  policy_rule_resource {
    name   = "policy_rule_resource-example"
    action = "ALLOW"
    conditions {
      negated  = false
      operator = "OR"
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_exporter"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_ip_anchoring"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_zapp"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_edge_connector"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_machine_tunnel"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_browser_isolation"]
      }
      operands {
        object_type = "CLIENT_TYPE"
        values      = ["zpn_client_type_slogger"]
      }
    }
  }
  connector_groups {
    id = [data.zpa_app_connector_group.example.id]
  }
}

// Retrieve the App Connector Group ID
data "zpa_app_connector_group" "example" {
  name = "SGIO-Vancouver"
}

// Retrieve LSS Config Format
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
    log_type="zpn_trans_log"
}