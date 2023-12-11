resource "zpa_ba_certificate" "this" {
    name = "server.example.com"
    description = "server.example.com"
    cert_blob = <<-EOT
-----BEGIN CERTIFICATE-----
MIIDyzCCArOgAwIBAgIUekBD+iu64583B3u5ew7Bqj2O5cQwDQYJKoZIhvcNAQEL
BQAwgY0xCzAJBgNVBAYTAkNBMRkwFwYDVQQIDBBCcml0aXNoIENvbHVtYmlhMRIw
EAYDVQQHDAlWYW5jb3V2ZXIxFTATBgNVBAoMDEJELUhhc2hpQ29ycDEVMBMGA1UE
-----END CERTIFICATE-----
    EOT
}
