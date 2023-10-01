// Create Log Receiver Configuration
resource "zpa_lss_config_controller" "example" {
  config {
    name        = "Example"
    description = "Example"
    enabled     = true
    format      = data.zpa_lss_config_log_type_formats.zpn_ast_auth_log.json
    lss_host    = "192.168.1.1"
    lss_port    = "11001"
    source_log_type = "zpn_ast_auth_log"
    use_tls         = true
    filter = [
                "ZPN_STATUS_AUTH_FAILED",
                "ZPN_STATUS_DISCONNECTED",
                "ZPN_STATUS_AUTHENTICATED"
            ]
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.example.id ]
  }
}

// Retrieve the App Connector Group ID
data "zpa_app_connector_group" "example" {
  name = "Example"
}

// Retrieve LSS Config Format
data "zpa_lss_config_log_type_formats" "zpn_ast_auth_log" {
    log_type="zpn_ast_auth_log"
}
