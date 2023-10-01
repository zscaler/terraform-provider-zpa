// Retrieve Log Receiver Information
data "zpa_lss_config_controller" "example" {
  id = zpa_lss_config_controller.example
}

output "zpa_lss_config_controller" {
  value = data.zpa_lss_config_controller.example
}
