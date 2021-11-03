// Create Browser Access Application Segment
resource "zpa_browser_access" "crm_browser_access" {
  name             = "CRM Application"
  description      = "CRM Application"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  tcp_port_ranges  = ["80", "80", "8080", "8080"]
  domain_names     = ["crm.example.com"]
  segment_group_id = data.zpa_segment_group.crm_app_group.id

  clientless_apps {
    name                 = "crm.example.com"
    application_protocol = "HTTP"
    application_port     = "8080"
    certificate_id       = data.zpa_ba_certificate.crm_ba.id
    trust_untrusted_cert = true
    enabled              = true
    domain               = "crm.example.com"
  }
  server_groups {
    id = [
      data.zpa_server_group.crm_servers.id
    ]
  }
}

// Browser Access Certificate
data "zpa_ba_certificate" "crm_ba" {
  name = "crm.example.com"
}

// Server Group
data "zpa_server_group" "crm_servers" {
  name = "CRM Servers"
}

// Segment Group
data "zpa_segment_group" "crm_app_group" {
  name = "CRM App group"
}