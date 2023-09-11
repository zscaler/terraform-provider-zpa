resource "zpa_cloud_browser_isolation_certificate" "this" {
    name = "CBI Certificate"
    pem = file("cert.pem")
}

# Warning: Certificate must be in PEM format