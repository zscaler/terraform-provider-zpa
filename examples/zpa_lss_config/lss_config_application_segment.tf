/*
// Create Log Receiver Configuration
resource "zpa_lss_config_controller" "example" {
  config {
    name            = "Example"
    description     = "Example"
    enabled         = true
    audit_message   = "{\"logType\":\"User Activity\",\"tcpPort\":\"11001\",\"appConnectorGroups\":[{\"name\":\"SGIO-Vancouver\",\"id\":\"216196257331281931\"}],\"domainOrIpAddress\":\"192.168.1.1\",\"logStreamContent\":\"{\\\"LogTimestamp\\\": %j{LogTimestamp:time},\\\"Customer\\\": %j{Customer},\\\"SessionID\\\": %j{SessionID},\\\"ConnectionID\\\": %j{ConnectionID},\\\"InternalReason\\\": %j{InternalReason},\\\"ConnectionStatus\\\": %j{ConnectionStatus},\\\"IPProtocol\\\": %d{IPProtocol},\\\"DoubleEncryption\\\": %d{DoubleEncryption},\\\"Username\\\": %j{Username},\\\"ServicePort\\\": %d{ServicePort},\\\"ClientPublicIP\\\": %j{ClientPublicIP},\\\"ClientPrivateIP\\\": %j{ClientPrivateIP},\\\"ClientLatitude\\\": %f{ClientLatitude},\\\"ClientLongitude\\\": %f{ClientLongitude},\\\"ClientCountryCode\\\": %j{ClientCountryCode},\\\"ClientZEN\\\": %j{ClientZEN},\\\"Policy\\\": %j{Policy},\\\"Connector\\\": %j{Connector},\\\"ConnectorZEN\\\": %j{ConnectorZEN},\\\"ConnectorIP\\\": %j{ConnectorIP},\\\"ConnectorPort\\\": %d{ConnectorPort},\\\"Host\\\": %j{Host},\\\"Application\\\": %j{Application},\\\"AppGroup\\\": %j{AppGroup},\\\"Server\\\": %j{Server},\\\"ServerIP\\\": %j{ServerIP},\\\"ServerPort\\\": %d{ServerPort},\\\"PolicyProcessingTime\\\": %d{PolicyProcessingTime},\\\"ServerSetupTime\\\": %d{ServerSetupTime},\\\"TimestampConnectionStart\\\": %j{TimestampConnectionStart:iso8601},\\\"TimestampConnectionEnd\\\": %j{TimestampConnectionEnd:iso8601},\\\"TimestampCATx\\\": %j{TimestampCATx:iso8601},\\\"TimestampCARx\\\": %j{TimestampCARx:iso8601},\\\"TimestampAppLearnStart\\\": %j{TimestampAppLearnStart:iso8601},\\\"TimestampZENFirstRxClient\\\": %j{TimestampZENFirstRxClient:iso8601},\\\"TimestampZENFirstTxClient\\\": %j{TimestampZENFirstTxClient:iso8601},\\\"TimestampZENLastRxClient\\\": %j{TimestampZENLastRxClient:iso8601},\\\"TimestampZENLastTxClient\\\": %j{TimestampZENLastTxClient:iso8601},\\\"TimestampConnectorZENSetupComplete\\\": %j{TimestampConnectorZENSetupComplete:iso8601},\\\"TimestampZENFirstRxConnector\\\": %j{TimestampZENFirstRxConnector:iso8601},\\\"TimestampZENFirstTxConnector\\\": %j{TimestampZENFirstTxConnector:iso8601},\\\"TimestampZENLastRxConnector\\\": %j{TimestampZENLastRxConnector:iso8601},\\\"TimestampZENLastTxConnector\\\": %j{TimestampZENLastTxConnector:iso8601},\\\"ZENTotalBytesRxClient\\\": %d{ZENTotalBytesRxClient},\\\"ZENBytesRxClient\\\": %d{ZENBytesRxClient},\\\"ZENTotalBytesTxClient\\\": %d{ZENTotalBytesTxClient},\\\"ZENBytesTxClient\\\": %d{ZENBytesTxClient},\\\"ZENTotalBytesRxConnector\\\": %d{ZENTotalBytesRxConnector},\\\"ZENBytesRxConnector\\\": %d{ZENBytesRxConnector},\\\"ZENTotalBytesTxConnector\\\": %d{ZENTotalBytesTxConnector},\\\"ZENBytesTxConnector\\\": %d{ZENBytesTxConnector},\\\"Idp\\\": %j{Idp}}\\\\n\",\"name\":\"LSS App Connector Status\",\"description\":\"LSS App Connector Status\",\"sessionStatuses\":[],\"enabled\":true,\"useTls\":true,\"policy\":{\"policyType\":\"Log Receiver Policy\",\"name\":\"SIEM selection rule for LSS App Connector Status\",\"action\":\"LOG\",\"ruleOrder\":\"1\"}}"
    format          = "{\"LogTimestamp\": %j{LogTimestamp:time},\"Customer\": %j{Customer},\"SessionID\": %j{SessionID},\"ConnectionID\": %j{ConnectionID},\"InternalReason\": %j{InternalReason},\"ConnectionStatus\": %j{ConnectionStatus},\"IPProtocol\": %d{IPProtocol},\"DoubleEncryption\": %d{DoubleEncryption},\"Username\": %j{Username},\"ServicePort\": %d{ServicePort},\"ClientPublicIP\": %j{ClientPublicIP},\"ClientPrivateIP\": %j{ClientPrivateIP},\"ClientLatitude\": %f{ClientLatitude},\"ClientLongitude\": %f{ClientLongitude},\"ClientCountryCode\": %j{ClientCountryCode},\"ClientZEN\": %j{ClientZEN},\"Policy\": %j{Policy},\"Connector\": %j{Connector},\"ConnectorZEN\": %j{ConnectorZEN},\"ConnectorIP\": %j{ConnectorIP},\"ConnectorPort\": %d{ConnectorPort},\"Host\": %j{Host},\"Application\": %j{Application},\"AppGroup\": %j{AppGroup},\"Server\": %j{Server},\"ServerIP\": %j{ServerIP},\"ServerPort\": %d{ServerPort},\"PolicyProcessingTime\": %d{PolicyProcessingTime},\"ServerSetupTime\": %d{ServerSetupTime},\"TimestampConnectionStart\": %j{TimestampConnectionStart:iso8601},\"TimestampConnectionEnd\": %j{TimestampConnectionEnd:iso8601},\"TimestampCATx\": %j{TimestampCATx:iso8601},\"TimestampCARx\": %j{TimestampCARx:iso8601},\"TimestampAppLearnStart\": %j{TimestampAppLearnStart:iso8601},\"TimestampZENFirstRxClient\": %j{TimestampZENFirstRxClient:iso8601},\"TimestampZENFirstTxClient\": %j{TimestampZENFirstTxClient:iso8601},\"TimestampZENLastRxClient\": %j{TimestampZENLastRxClient:iso8601},\"TimestampZENLastTxClient\": %j{TimestampZENLastTxClient:iso8601},\"TimestampConnectorZENSetupComplete\": %j{TimestampConnectorZENSetupComplete:iso8601},\"TimestampZENFirstRxConnector\": %j{TimestampZENFirstRxConnector:iso8601},\"TimestampZENFirstTxConnector\": %j{TimestampZENFirstTxConnector:iso8601},\"TimestampZENLastRxConnector\": %j{TimestampZENLastRxConnector:iso8601},\"TimestampZENLastTxConnector\": %j{TimestampZENLastTxConnector:iso8601},\"ZENTotalBytesRxClient\": %d{ZENTotalBytesRxClient},\"ZENBytesRxClient\": %d{ZENBytesRxClient},\"ZENTotalBytesTxClient\": %d{ZENTotalBytesTxClient},\"ZENBytesTxClient\": %d{ZENBytesTxClient},\"ZENTotalBytesRxConnector\": %d{ZENTotalBytesRxConnector},\"ZENBytesRxConnector\": %d{ZENBytesRxConnector},\"ZENTotalBytesTxConnector\": %d{ZENTotalBytesTxConnector},\"ZENBytesTxConnector\": %d{ZENBytesTxConnector},\"Idp\": %j{Idp}}\\n"
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
        values      = [data.zpa_application_segment.example.id]
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

// Retrieve the Application Segment ID
data "zpa_application_segment" "example" {
  name = "Example"
}
*/