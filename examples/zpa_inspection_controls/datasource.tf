data "zpa_inspection_predefined_controls" "example1" {
  name = "Failed to parse request body"
  version = "OWASP_CRS/3.3.0"
}

output "zpa_inspection_predefined_controls" {
  value = data.zpa_inspection_predefined_controls.example1
}

data "zpa_inspection_predefined_controls" "example2" {
  name = "Multipart request body failed strict validation"
  version = "OWASP_CRS/3.3.0"
}

output "zpa_inspection_predefined_controls" {
  value = data.zpa_inspection_predefined_controls.example2
}