terraform {
  required_providers {
    zpa = {
      version = "2.0.9"
      source  = "zscaler.com/zpa/zpa"
    }
  }
  required_version = ">= 0.13"
}

provider "zpa" {}

data "zpa_policy_access_rule" "all_other_services" {
  name                          = "All Other Services"
}

output "zpa_policy_access_rule" {
    value = data.zpa_policy_access_rule.all_other_services
}