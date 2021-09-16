terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_app_connector_group" "example" {
  name = "SGIO-Vancouver"
}

resource "zpa_application_server" "example20" {
  name                          = "example20.securitygeek.io"
  description                   = "example20.securitygeek.io"
  address                       = "2.2.2.2"
  enabled                       = true
}


resource "zpa_server_group" "example20" {
  name = "example20"
  description = "example20"
  enabled = false
  dynamic_discovery = false
  app_connector_groups {
    id = [data.zpa_app_connector_group.example.id]
  }
  servers {
    id = [zpa_application_server.example20.id]
  }
}


output "all_zpa_server_group" {
  value = zpa_server_group.example20
}