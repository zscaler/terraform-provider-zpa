terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_trusted_network" "example" {
 name = "Corp-Trusted-Networks"
}

output "get_trusted_network" {
  value = data.zpa_trusted_network.example
}