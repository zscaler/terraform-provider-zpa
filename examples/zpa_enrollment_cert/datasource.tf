data "zpa_enrollment_cert" "sales_ba" {
    name = "sales.securitygeek.io"
}

output "zpa_enrollment_cert" {
  value = data.zpa_enrollment_cert.sales_ba
}