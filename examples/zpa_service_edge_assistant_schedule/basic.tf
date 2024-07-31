resource "zpa_service_edge_assistant_schedule" "this" {
  frequency = "days"
  frequency_interval = "7"
  enabled = true
  delete_disabled = true
}