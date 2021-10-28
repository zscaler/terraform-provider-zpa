
terraform {
  required_providers {
    zpa = {
      version = "1.0.0"
      source  = "zscaler.com/zpa/zpa"
    }
  }
}

provider "zpa" {}

data "zpa_global_access_policy" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

output "zpa_global_access_policy_access_policy" {
    value = data.zpa_global_access_policy.access_policy
}

data "zpa_global_access_policy" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

output "zpa_global_access_policy_timeout_policy" {
    value = data.zpa_global_access_policy.timeout_policy
}

data "zpa_global_access_policy" "reauth_policy" {
    policy_type = "REAUTH_POLICY"
}

output "zpa_global_access_policy_reauth_policy" {
    value = data.zpa_global_access_policy.reauth_policy
}

data "zpa_global_access_policy" "siem_policy" {
    policy_type = "SIEM_POLICY"
}

output "zpa_global_access_policy_siem_policy" {
    value = data.zpa_global_access_policy.siem_policy
}

data "zpa_global_access_policy" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

output "zpa_global_access_client_forwarding_policy" {
    value = data.zpa_global_access_policy.client_forwarding_policy
}



