/*
// Create Log Receiver Configuration
resource "zpa_lss_config_controller" "example" {
  config {
    name        = "Example"
    description = "Example"
    enabled     = true
    format      = "{\"LogTimestamp\": %j{LogTimestamp:time},\"Customer\": %j{Customer},\"SessionID\": %j{SessionID},\"ConnectionID\": %j{ConnectionID},\"InternalReason\": %j{InternalReason},\"ConnectionStatus\": %j{ConnectionStatus},\"IPProtocol\": %d{IPProtocol},\"DoubleEncryption\": %d{DoubleEncryption},\"Username\": %j{Username},\"ServicePort\": %d{ServicePort},\"ClientPublicIP\": %j{ClientPublicIP},\"ClientPrivateIP\": %j{ClientPrivateIP},\"ClientLatitude\": %f{ClientLatitude},\"ClientLongitude\": %f{ClientLongitude},\"C&hellip;iso8601},\"TimestampZENLastTxConnector\": %j{TimestampZENLastTxConnector:iso8601},\"ZENTotalBytesRxClient\": %d{ZENTotalBytesRxClient},\"ZENBytesRxClient\": %d{ZENBytesRxClient},\"ZENTotalBytesTxClient\": %d{ZENTotalBytesTxClient},\"ZENBytesTxClient\": %d{ZENBytesTxClient},\"ZENTotalBytesRxConnector\": %d{ZENTotalBytesRxConnector},\"ZENBytesRxConnector\": %d{ZENBytesRxConnector},\"ZENTotalBytesTxConnector\": %d{ZENTotalBytesTxConnector},\"ZENBytesTxConnector\": %d{ZENBytesTxConnector},\"Idp\": %j{Idp}}\\n"
    lss_host    = "2.2.2.2"
    lss_port    = "5000"
    filter = ["BRK_MT_SETUP_FAIL_BIND_TO_AST_LOCAL_OWNER",
      "CLT_CONN_FAILED",
      "BRK_MT_TERMINATED_IDLE_TIMEOUT",
      "MT_CLOSED_TLS_CONN_GONE_CLIENT_CLOSED"
    ]
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
        values = ["zpn_client_type_exporter",
          "zpn_client_type_ip_anchoring",
          "zpn_client_type_zapp",
          "zpn_client_type_edge_connector",
          "zpn_client_type_machine_tunnel",
          "zpn_client_type_slogger"
        ]
      }
    }
  }
  connector_groups {
    id = [data.zpa_app_connector_group.example.id]
  }
}

// Retrieve the App Connector Group ID
data "zpa_app_connector_group" "example" {
  name = "Example"
}
*/