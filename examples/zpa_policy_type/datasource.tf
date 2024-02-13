// Get information for "GLOBAL_POLICY" ID
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

// Get information for "CAPABILITIES_POLICY" ID
data "zpa_policy_type" "capabilities_policy" {
    policy_type = "CAPABILITIES_POLICY"
}

// Get information for "CLIENT_FORWARDING_POLICY" ID
data "zpa_policy_type" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

// Get information for "CREDENTIAL_POLICY" ID
data "zpa_policy_type" "credential_policy" {
    policy_type = "CREDENTIAL_POLICY"
}

// Get information for "INSPECTION_POLICY" ID
data "zpa_policy_type" "inspection_policy" {
    policy_type = "INSPECTION_POLICY"
}

// Get information for "INSPECTION_POLICY" ID
data "zpa_policy_type" "inspection_policy" {
    policy_type = "ISOLATION_POLICY"
}


// Get information for "REAUTH_POLICY" ID
data "zpa_policy_type" "reauth_policy" {
    policy_type = "REAUTH_POLICY"
}

// Get information for "TIMEOUT_POLICY" ID
data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

data "zpa_policy_type" "inspection_policy" {
    policy_type = "REDIRECTION_POLICY"
}

