terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}


resource "zpa_application_server" "example10" {
  name                          = "example10.securitygeek.io"
  description                   = "example10.securitygeek.io"
  address                       = "1.1.1.1"
  enabled                       = true
}

