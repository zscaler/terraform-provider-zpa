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
        object_type = "APP_GROUP"
        values      = [zpa_segment_group.other_lss_name.id]
      }
    }
  }
  connector_groups {
    id = [data.zpa_app_connector_group.example.id]
  }
}

// Create Segment Group
resource "zpa_segment_group" "sg_lss_config" {
  name = "Segment Group LSS"
  description = "Segment Group LSS"
  enabled = true
  policy_migrated = true
 }

// Retrieve LSS Config Format
data "zpa_lss_config_log_type_formats" "zpn_trans_log" {
    log_type="zpn_trans_log"
}

// Retrieve the App Connector Group ID
data "zpa_app_connector_group" "app_connector_lss" {
  name = "App Connector LSS"
}
