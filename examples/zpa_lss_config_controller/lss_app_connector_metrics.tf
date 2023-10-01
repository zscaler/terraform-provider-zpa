# Get Log Type Format - "App Connector Metrics"
data "zpa_lss_config_log_type_formats" "zpn_ast_comprehensive_stats" {
  log_type = "zpn_ast_comprehensive_stats"
}

data "zpa_policy_type" "lss_siem_policy" {
  policy_type = "SIEM_POLICY"
}

data "zpa_app_connector_group" "this" {
 name = "Example100"
}
resource "zpa_lss_config_controller" "lss_app_connector_metrics" {
  config {
    name            = "LSS App Connector Metrics"
    description     = "LSS App Connector Metrics"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_ast_comprehensive_stats.json
    lss_host        = "splunk1.acme.com"
    lss_port        = "5001"
    source_log_type = "zpn_ast_comprehensive_stats"
    use_tls         = true
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.this.id ]
  }
}
