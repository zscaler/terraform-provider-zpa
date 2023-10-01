# Get Log Type Format - "App Connector Status"
data "zpa_lss_config_log_type_formats" "zpn_http_trans_log" {
  log_type = "zpn_http_trans_log"
}

data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

data "zpa_app_connector_group" "this" {
 name = "Example100"
}
resource "zpa_lss_config_controller" "lss_web_browser" {
  config {
    name            = "LSS Web Browser"
    description     = "LSS Web Browser"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_http_trans_log.json
    lss_host        = "splunk1.acme.com"
    lss_port        = "5001"
    source_log_type = "zpn_http_trans_log"
    use_tls         = true
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.this.id ]
  }
}