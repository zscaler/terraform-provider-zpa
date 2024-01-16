// Retrieve All Assistant Schedules
data "zpa_app_connector_assistant_schedule" "this" {}

// Retrieve A Specific Assistant Schedule by ID
data "zpa_app_connector_assistant_schedule" "this" {
    id = "1"
}

// Retrieve A Specific Assistant Schedule by the Customer ID
data "zpa_app_connector_assistant_schedule" "this" {
    customer_id = "1234567891012"
}

