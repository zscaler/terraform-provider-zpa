// Create Server Group - Dynamic Discovery "True"
resource "zpa_server_group" "example" {
  name = "Example"
  description = "Example"
  enabled = true
  dynamic_discovery = true
  app_connector_groups {
    id = [data.zpa_app_connector_group.example.id]
  }
}

data "zpa_app_connector_group" "example" {
  name = "Example"
}

// Create Server Group - Dynamic Discovery "False"
resource "zpa_server_group" "example" {
  name = "Example"
  description = "Example"
  enabled = false
  dynamic_discovery = false
  app_connector_groups {
    id = [data.zpa_app_connector_group.example.id]
  }
  servers {
    id = [zpa_application_server.example.id]
  }
}

data "zpa_app_connector_group" "example" {
  name = "Example"
}

resource "zpa_application_server" "example" {
  name                          = "Server1"
  description                   = "Server1"
  address                       = "192.168.1.1"
  enabled                       = true
}