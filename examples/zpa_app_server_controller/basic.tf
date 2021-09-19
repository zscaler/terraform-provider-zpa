resource "zpa_application_server" "example" {
  name                          = "Example"
  description                   = "Example"
  address                       = "192.168.1.1"
  enabled                       = true
  app_server_group_ids          = [data.zpa_server_group.example.id]
}

data "zpa_server_group" "example"{
    name = "Server-Group-Example"
}