terraform {
  required_providers {
    zpa = {
      source  = "zscaler.com/zpa/zpa"
      version = "2.2.2"
    }
  }
}
provider "zpa" {}

data "zpa_app_connector_controller" "example" {
  name = "AWS-VPC100-App-Connector"
}

output "zpa_app_connector_controller" {
  value = data.zpa_app_connector_controller.example
}
