// Create Application Segment Privileged Remote Access (PRA)

resource "zpa_application_segment_pra" "this" {
  name             = "ZPA_PRA_Example"
  description      = "ZPA_PRA_Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = [ "22", "22", "3389", "3389" ]
  domain_names     = [ "ssh_pra.example.com", "rdp_pra.example.com" ]
  segment_group_id = zpa_segment_group.this.id
  server_groups {
    id = [ zpa_server_group.this.id ]
  }
  common_apps_dto {
    apps_config {
      name                 = "ssh_pra"
      domain               = "ssh_pra.example.com"
      application_protocol = "SSH"
      application_port     = "22"
      enabled = true
      app_types = [ "SECURE_REMOTE_ACCESS" ]
    }
      apps_config {
      name                 = "rdp_pra"
      domain               = "rdp_pra.example.com"
      application_protocol = "RDP"
      connection_security  = "ANY"
      application_port     = "3389"
      enabled = true
      app_types = [ "SECURE_REMOTE_ACCESS" ]
    }
  }
}

// Create Segment Group
resource "zpa_segment_group" "this" {
  name            = "Example_Segment_Group"
  description     = "Example_Segment_Group"
  enabled         = true
}

// Create Server Group
resource "zpa_server_group" "this" {
  name              = "Example_Server_Group"
  description       = "Example_Server_Group"
  enabled           = true
  dynamic_discovery = true
  app_connector_groups {
    id = [ data.zpa_app_connector_group.this.id ]
  }
}