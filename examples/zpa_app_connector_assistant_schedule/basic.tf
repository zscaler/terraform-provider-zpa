resource "zpa_app_connector_assistant_schedule" "this" {
  frequency = "days"
  frequency_interval = "7"
  enabled = true
  delete_disabled = true
}