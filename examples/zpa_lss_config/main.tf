terraform {
  required_providers {
    zpa = {
      version = "1.0.0"
      source  = "zscaler.com/zpa/zpa"
    }
  }
}

provider "zpa" {}

resource "zpa_lss_config_controller" "example2" {
  config {
    name        = "example2"
    description = "example2"
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
    name   = "policy_rule_resource-example2"
    action = "ALLOW"
    conditions {
      negated  = false
      operator = "OR"
      operands {
        object_type = "CLIENT_TYPE"
        values = ["zpn_client_type_exporter"]
      }
    }
  }
  connector_groups {
    id = [data.zpa_app_connector_group.sgio-vancouver.id]
  }
}

data "zpa_app_connector_group" "sgio-vancouver" {
  name = "SGIO-Vancouver"
}


data "zpa_lss_config_controller" "example" {
  id = zpa_lss_config_controller.example2.id
}

output "zpa_lss_config_controller" {
  value = data.zpa_lss_config_controller.example
}

