// Create Log Receiver Configuration
resource "zpa_lss_config_controller" "example" {
  config {
    name        = "Example"
    description = "Example"
    enabled     = true
    format      = local.log_type
    lss_host    = "192.168.1.1"
    lss_port    = "11001"
    source_log_type = "zpn_audit_log"
    use_tls         = true
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.example.id ]
  }
}

// Retrieve the App Connector Group ID
data "zpa_app_connector_group" "example" {
  name = "SGIO-Vancouver"
}

// Encode JSON value to string
locals {
  log_type = jsonencode(data.zpa_lss_config_log_type_formats.zpn_audit_log.json)
}

// Retrieve LSS Config Format
data "zpa_lss_config_log_type_formats" "zpn_audit_log" {
    log_type="zpn_audit_log"
}