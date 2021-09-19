data "zpa_application_server" "example" {
    name = "server1"
}

output "get_zpa_application_server" {
  value = data.zpa_application_server.example-group
}