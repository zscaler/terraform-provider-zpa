data "zpa_service_edge_controller" "example" {
  name = "On-Prem-PSE"
}

output "zpa_service_edge_controller" {
  value = data.zpa_service_edge_controller.example
}