// Get information for "GLOBAL_POLICY" ID
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

output "zpa_policy_type_access_policy" {
    value = data.zpa_policy_type.access_policy
}

// Get information for "TIMEOUT_POLICY" ID
data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

output "zpa_policy_type_timeout_policy" {
    value = data.zpa_policy_type.timeout_policy
}

// Get information for "REAUTH_POLICY" ID
data "zpa_policy_type" "reauth_policy" {
    policy_type = "REAUTH_POLICY"
}

output "zpa_policy_type_reauth_policy" {
    value = data.zpa_policy_type.reauth_policy
}

// Get information for "SIEM_POLICY" ID
data "zpa_policy_type" "siem_policy" {
    policy_type = "SIEM_POLICY"
}

output "zpa_policy_type_siem_policy" {
    value = data.zpa_policy_type.siem_policy
}

// Get information for "CLIENT_FORWARDING_POLICY" ID
data "zpa_policy_type" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

output "zpa_policy_type_client_forwarding_policy" {
    value = data.zpa_policy_type.client_forwarding_policy
}



