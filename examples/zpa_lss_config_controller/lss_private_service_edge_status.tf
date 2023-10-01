# Get Log Type Format - "Private Service Edge Status"
data "zpa_lss_config_log_type_formats" "zpn_sys_auth_log" {
  log_type = "zpn_sys_auth_log"
}

data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

data "zpa_app_connector_group" "this" {
 name = "Example100"
}
resource "zpa_lss_config_controller" "lss_private_service_edge_group" {
  config {
    name            = "LSS Private Service Edge Status"
    description     = "LSS Private Service Edge Status"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_sys_auth_log.json
    lss_host        = "splunk1.acme.com"
    lss_port        = "5001"
    source_log_type = "zpn_sys_auth_log"
    use_tls         = true
    filter = ["ZPN_STATUS_AUTH_FAILED", "ZPN_STATUS_DISCONNECTED", "ZPN_STATUS_AUTHENTICATED"]
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.this.id ]
  }
}
