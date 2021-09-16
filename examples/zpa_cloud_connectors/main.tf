terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}



data "zpa_cloud_connector_group" "all" {
}

output "all_cloud_connector_group" {
  value = data.zpa_cloud_connector_group.all
}
