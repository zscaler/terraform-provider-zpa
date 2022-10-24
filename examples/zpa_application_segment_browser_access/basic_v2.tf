// Create Browser Access Application Segment
resource "zpa_application_segment_browser_access" "crm_browser_access" {
  name             = "CRM Application"
  description      = "CRM Application"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  tcp_port_range {
    from = "80"
    to   = "80"
  }
  tcp_port_range {
    from = "8080"
    to   = "8080"
  }
  domain_names     = ["jenkins.securitygeek.io"]
  segment_group_id = data.zpa_segment_group.crm_app_group.id

  clientless_apps {
    name                 = "jenkins.securitygeek.io"
    application_protocol = "HTTP"
    application_port     = "8080"
    certificate_id       = data.zpa_ba_certificate.jenkins_ca.id
    trust_untrusted_cert = true
    enabled              = true
    domain               = "jenkins.securitygeek.io"
  }
  server_groups {
    id = [
      data.zpa_server_group.crm_servers.id
    ]
  }
}

// Browser Access Certificate
data "zpa_ba_certificate" "jenkins_ca" {
  name = "jenkins.securitygeek.io"
}

// Server Group
data "zpa_server_group" "crm_servers" {
  name = "CRM Servers"
}

// Segment Group
data "zpa_segment_group" "crm_app_group" {
  name = "CRM App group"
}