terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

resource "zpa_lss_config_controller" "example" {
    config {
    name = "example" 
    description = "example" 
    enabled = true
    //format = "json": "{\"LogTimestamp\": %j{LogTimestamp:time},\"Customer\": %j{Customer},\"SessionID\": %j{SessionID},\"ConnectionID\": %j{ConnectionID},\"InternalReason\": %j{InternalReason},\"ConnectionStatus\": %j{ConnectionStatus},\"IPProtocol\": %d{IPProtocol},\"DoubleEncryption\": %d{DoubleEncryption},\"Username\": %j{Username},\"ServicePort\": %d{ServicePort},\"ClientPublicIP\": %j{ClientPublicIP},\"ClientPrivateIP\": %j{ClientPrivateIP},\"ClientLatitude\": %f{ClientLatitude},\"ClientLongitude\": %f{ClientLongitude},\"ClientCountryCode\": %j{ClientCountryCode},\"ClientZEN\": %j{ClientZEN},\"Policy\": %j{Policy},\"Connector\": %j{Connector},\"ConnectorZEN\": %j{ConnectorZEN},\"ConnectorIP\": %j{ConnectorIP},\"ConnectorPort\": %d{ConnectorPort},\"Host\": %j{Host},\"Application\": %j{Application},\"AppGroup\": %j{AppGroup},\"Server\": %j{Server},\"ServerIP\": %j{ServerIP},\"ServerPort\": %d{ServerPort},\"PolicyProcessingTime\": %d{PolicyProcessingTime},\"ServerSetupTime\": %d{ServerSetupTime},\"TimestampConnectionStart\": %j{TimestampConnectionStart:iso8601},\"TimestampConnectionEnd\": %j{TimestampConnectionEnd:iso8601},\"TimestampCATx\": %j{TimestampCATx:iso8601},\"TimestampCARx\": %j{TimestampCARx:iso8601},\"TimestampAppLearnStart\": %j{TimestampAppLearnStart:iso8601},\"TimestampZENFirstRxClient\": %j{TimestampZENFirstRxClient:iso8601},\"TimestampZENFirstTxClient\": %j{TimestampZENFirstTxClient:iso8601},\"TimestampZENLastRxClient\": %j{TimestampZENLastRxClient:iso8601},\"TimestampZENLastTxClient\": %j{TimestampZENLastTxClient:iso8601},\"TimestampConnectorZENSetupComplete\": %j{TimestampConnectorZENSetupComplete:iso8601},\"TimestampZENFirstRxConnector\": %j{TimestampZENFirstRxConnector:iso8601},\"TimestampZENFirstTxConnector\": %j{TimestampZENFirstTxConnector:iso8601},\"TimestampZENLastRxConnector\": %j{TimestampZENLastRxConnector:iso8601},\"TimestampZENLastTxConnector\": %j{TimestampZENLastTxConnector:iso8601},\"ZENTotalBytesRxClient\": %d{ZENTotalBytesRxClient},\"ZENBytesRxClient\": %d{ZENBytesRxClient},\"ZENTotalBytesTxClient\": %d{ZENTotalBytesTxClient},\"ZENBytesTxClient\": %d{ZENBytesTxClient},\"ZENTotalBytesRxConnector\": %d{ZENTotalBytesRxConnector},\"ZENBytesRxConnector\": %d{ZENBytesRxConnector},\"ZENTotalBytesTxConnector\": %d{ZENTotalBytesTxConnector},\"ZENBytesTxConnector\": %d{ZENBytesTxConnector},\"Idp\": %j{Idp},\"ClientToClient\": %j{c2c}}\\n"
    lss_host = "2.2.2.2"
    lss_port = "5000"
    filter = [ "BRK_MT_SETUP_FAIL_BIND_TO_AST_LOCAL_OWNER",
                "CLT_CONN_FAILED",
                "BRK_MT_TERMINATED_IDLE_TIMEOUT",
                "MT_CLOSED_TLS_CONN_GONE_CLIENT_CLOSED"
            ]
    source_log_type = "zpn_trans_log"
    use_tls = true
    }
    connector_groups {
        id = [ data.zpa_app_connector_group.sgio-vancouver.id ]
    }
}

data "zpa_app_connector_group" "sgio-vancouver" {
  name = "SGIO-Vancouver"
}

/*
data "zpa_lss_config_controller" "example" {
    id = "216196257331287887"
}

output "zpa_lss_config_controller" {
    value = data.zpa_lss_config_controller.example
}
*/