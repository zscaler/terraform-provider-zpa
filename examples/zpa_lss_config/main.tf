terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_lss_config_controller" "example" {
    id = "216196257331287880"
}

output "zpa_lss_config_controller" {
    value = data.zpa_lss_config_controller.example
}