// Retrieve Provisioning Key for Service Edge Group
data "zpa_provisioning_key" "example" {
    name             = "Service_Edge_Provisioning_Key"
    association_type = "SERVICE_EDGE_GRP"
}

output "zpa_provisioning_key_example" {
    value = data.zpa_provisioning_key.example
}

// NOTE: ASSOCIATION_TYPE is madantory due to API requirement.