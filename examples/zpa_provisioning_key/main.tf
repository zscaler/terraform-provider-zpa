terraform {
  required_providers {
    zpa = {
      version = "1.0.0"
      source  = "zscaler.com/zpa/zpa"
    }
  }
}

provider "zpa" {}

resource "zpa_provisioning_key" "example" {
  name             = "zpa_provisioning_key_example"
  association_type = "SERVICE_EDGE_GRP"
  max_usage        = "1"
  enrollment_cert_id = "10242"
  zcomponent_id = "216196257331288679"
}

data "zpa_provisioning_key" "example" {
  name             = zpa_provisioning_key.example.name
  association_type = zpa_provisioning_key.example.association_type
}

output "zpa_provisioning_key" {
  value = data.zpa_provisioning_key.example
}
