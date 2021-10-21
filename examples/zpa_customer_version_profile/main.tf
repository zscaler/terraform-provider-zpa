terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_customer_version_profile" "example"{
    name = "Default"
}

output "zpa_customer_version_profile" {
    value = data.zpa_customer_version_profile.example
}