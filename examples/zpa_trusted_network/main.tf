terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}


// Testing Data Source Trusted Network
data "zpa_trusted_network" "example" {
 name = "SGIO-Trusted-Networks"
}

output "all_trusted_network" {
  value = data.zpa_trusted_network.example
}