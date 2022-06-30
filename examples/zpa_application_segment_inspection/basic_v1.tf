// ZPA Inspection Application Segment using HTTPS

// Retrieve Certificate
data "zpa_ba_certificate" "jenkins" {
  name = "jenkins.example.com"
}

// Create Inspection Application Segment
resource "zpa_application_segment_inspection" "this" {
  name             = "ZPA_Inspection_Example"
  description      = "ZPA_Inspection_Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["443", "443"]
  domain_names     = ["jenkins.example.com"]
  segment_group_id = zpa_segment_group.this.id
  server_groups {
    id = [zpa_server_group.this.id]
  }
  common_apps_dto {
    apps_config {
      name                 = "jenkins.example.com"
      domain               = "jenkins.example.com"
      application_protocol = "HTTPS"
      application_port     = "443"
      certificate_id       = data.zpa_ba_certificate.jenkins.id
      enabled              = true
      app_types            = [ "INSPECT" ]
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
    id = [data.zpa_app_connector_group.this.id]
  }
}