terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}


data "zpa_lss_config_client_types" "example" {}

output "zpa_lss_config_client_types"{
    value = data.zpa_lss_config_client_types.example
}


data "zpa_lss_config_status_codes" "example" {}

output "zpa_lss_config_status_codes"{
    value = data.zpa_lss_config_status_codes.example
}

data "zpa_lss_config_log_type_formats" "example" {
    log_type="zpn_auth_log"
}

output "zpa_lss_config_log_type_formats"{
    value = data.zpa_lss_config_log_type_formats.example
}
