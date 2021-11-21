data "zpa_trusted_network" "example" {
 name = "Corp-Trusted-Networks"
}

output "get_trusted_network" {
  value = data.zpa_trusted_network.example
}