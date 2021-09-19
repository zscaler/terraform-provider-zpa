data "zpa_ba_certificate" "example" {
 name = "example.acme.com"
}

output "get_zpa_ba_certificate" {
  value = data.zpa_ba_certificate.example
}