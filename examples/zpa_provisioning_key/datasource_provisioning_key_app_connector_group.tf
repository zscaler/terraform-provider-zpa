// Retrieve Provisioning Key for App Connector Group
data "zpa_provisioning_key" "example" {
    name            = "App_Connector_Provisioning_Key"
    association_type = "CONNECTOR_GRP"
}

output "zpa_provisioning_key_example" {
    value = data.zpa_provisioning_key.example
}

// NOTE: ASSOCIATION_TYPE is madantory due to API requirement.