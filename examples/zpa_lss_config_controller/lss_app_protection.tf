# Get Log Type Format - "App Protection"
data "zpa_lss_config_log_type_formats" "zpn_waf_http_exchanges_log" {
  log_type = "zpn_waf_http_exchanges_log"
}

data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

data "zpa_app_connector_group" "this" {
 name = "Example100"
}
resource "zpa_lss_config_controller" "lss_app_protection" {
  config {
    name            = "LSS App Protection"
    description     = "LSS App Protection"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_waf_http_exchanges_log.json
    lss_host        = "splunk1.acme.com"
    lss_port        = "5001"
    source_log_type = "zpn_waf_http_exchanges_log"
    use_tls         = true
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.this.id ]
  }
}
