terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

resource "zpa_service_edge_group" "example" {
  name                          = "Example"
  description                   = "Example"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  latitude                      = "49.1041779"
  longitude                     = "-122.6603519"
  location                      = "Langley City, BC, Canada"
  version_profile_id = "0"
  trusted_networks {
    id = [ data.zpa_trusted_network.example.id ]
  }
}


data "zpa_service_edge_group" "example" {
  name = "Example"
}

output "zpa_service_edge_group" {
  value = data.zpa_service_edge_group.example
}

data "zpa_trusted_network" "example" {
 name = "Corp-Trusted-Networks"
}