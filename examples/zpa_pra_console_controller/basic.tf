# Creates Privileged Remote Access Application Segment"
resource "zpa_application_segment_pra" "this" {
  name             = "Example"
  description      = "Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["3389", "3389"]
  domain_names     = [ "rdp_pra.example.com"]
  segment_group_id = zpa_segment_group.this.id
  common_apps_dto {
    apps_config {
      name                 = "rdp_pra"
      domain               = "rdp_pra.example.com"
      application_protocol = "RDP"
      connection_security  = "ANY"
      application_port     = "3389"
      enabled              = true
      app_types            = ["SECURE_REMOTE_ACCESS"]
    }
  }
}

data "zpa_application_segment_by_type" "this" {
    application_type = "SECURE_REMOTE_ACCESS"
    name = "rdp_pra"
    depends_on = [zpa_application_segment_pra.this]
}

# Creates Segment Group for Application Segment"
resource "zpa_segment_group" "this" {
  name        = "Example"
  description = "Example"
  enabled     = true
}

# Retrieves the Browser Access Certificate
data "zpa_ba_certificate" "this" {
  name = "pra01.example.com"
}

# Creates PRA Portal"
resource "zpa_pra_portal_controller" "this1" {
  name                      = "pra01.example.com"
  description               = "pra01.example.com"
  enabled                   = true
  domain                    = "pra01.example.com"
  certificate_id            = data.zpa_ba_certificate.this.id
  user_notification         = "Created with Terraform"
  user_notification_enabled = true
}


resource "zpa_pra_console_controller" "ssh_pra" {
  name        = "ssh_console"
  description = "Created with Terraform"
  enabled     = true
  pra_application {
    id = data.zpa_application_segment_by_type.this.id
  }
  pra_portals {
    id = [zpa_pra_portal_controller.this.id]
  }
}