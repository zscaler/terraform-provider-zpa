// Create Log Receiver Configuration
resource "zpa_lss_config_controller" "example" {
  config {
    name            = "Example"
    description     = "Example"
    enabled         = true
    format          = data.zpa_lss_config_log_type_formats.zpn_trans_log.json
    lss_host        = "192.168.1.1"
    lss_port        = "11001"
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
        object_type = "APP"
        values      = [ zpa_application_segment.as_lss_config.id ]
      }
    }
  }
  connector_groups {
    id = [ data.zpa_app_connector_group.app_connector_lss.id ]
  }
}

// Create Application Segment
resource "zpa_application_segment" "as_lss_config" {
    name = "App Segment LSS"
    description = "App Segment LSS"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    tcp_port_ranges = ["11001", "11001"]
    domain_names = ["*.acme.com"]
    segment_group_id = data.zpa_segment_group.sg_lss_config.id
    server_groups {
        id = [ data.zpa_server_group.srvg_lss_config.id ]
    }
}

// Retrieve LSS Config Format
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
    log_type="zpn_trans_log"
}

// Retrieve the App Connector Group ID
data "zpa_app_connector_group" "app_connector_lss" {
  name = "App Connector LSS"
}

// Retrieve Segment Group ID
data "zpa_segment_group" "sg_lss_config" {
   name = "Segment Group LSS"
 }

// Retrieve Server Group ID
data "zpa_server_group" "srvg_lss_config" {
  name = "Server Group LSS"
}