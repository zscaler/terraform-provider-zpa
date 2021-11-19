// Retrieve Provisioning Key for App Connector Group
data "zpa_provisioning_key" "example1" {
    name            = "App_Connector_Provisioning_Key"
    association_type = "CONNECTOR_GRP"
}

output "zpa_provisioning_key_example1" {
    value = data.zpa_provisioning_key.example1
}


// Retrieve Provisioning Key for Service Edge Group
data "zpa_provisioning_key" "example2" {
    name            = "Service_Edge_Provisioning_Key"
    association_type = "SERVICE_EDGE_GRP"
}

output "zpa_provisioning_key_example2" {
    value = data.zpa_provisioning_key.example2
}