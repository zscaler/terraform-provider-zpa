data "zpa_inspection_predefined_controls" "example" {
  name    = "Failed to parse request body"
  version = "OWASP_CRS/3.3.0"
}

output "zpa_inspection_predefined_controls" {
  value = data.zpa_inspection_predefined_controls.example
}

data "zpa_inspection_all_predefined_controls" "example" {
  version = "OWASP_CRS/3.3.0"
}

output "zpa_inspection_all_predefined_controls" {
  value = data.zpa_inspection_all_predefined_controls.example
}
