terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_idp_controller" "all" {
 name = "SGIO-User-Okta"
}

output "idp_controller" {
    value = data.zpa_idp_controller.all
}