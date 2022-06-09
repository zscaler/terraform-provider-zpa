data "zpa_application_segment_pra" "this" {
  name = "ZPA_PRA_Example"
}

output "zpa_application_segment_pra" {
  value = resource.zpa_application_segment_pra.this
}