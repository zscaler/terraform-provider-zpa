terraform {
  required_providers {
    zpa = {
      version = "1.0.0"
      source  = "zscaler.com/zpa/zpa"
    }
  }
}

provider "zpa" {}

resource "zpa_service_edge_group" "service_edge_group_sjc" {
  name                 = "Service Edge Group San Jose"
  description          = "Service Edge Group in San Jose"
  upgrade_day          = "SUNDAY"
  upgrade_time_in_secs = "66600"
  latitude             = "37.3382082"
  longitude            = "-121.8863286"
  location             = "San Jose, CA, USA"
  version_profile_id   = "0"
  trusted_networks {
    id = [data.zpa_trusted_network.example.id]
  }
}

resource "zpa_service_edge_group" "service_edge_group_nyc" {
  name                 = "Service Edge Group New York"
  description          = "Service Edge Group in New York"
  upgrade_day          = "SUNDAY"
  upgrade_time_in_secs = "66600"
  latitude             = "40.7128"
  longitude            = "-73.935242"
  location             = "New York, NY, USA"
  version_profile_id   = "0"
  trusted_networks {
    id = [data.zpa_trusted_network.example.id]
  }
}

data "zpa_trusted_network" "example" {
  name = "Corp-Trusted-Networks"
}