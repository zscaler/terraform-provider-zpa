terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_enrollment_cert" "sales_ba" {
    name = "sales.securitygeek.io"
}

output "zpa_enrollment_cert" {
  value = data.zpa_enrollment_cert.sales_ba
}