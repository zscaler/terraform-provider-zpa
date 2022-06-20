resource "zpa_application_server" "example" {
  name                          = "Example"
  description                   = "Example"
  address                       = "192.168.1.1"
  enabled                       = true
}
