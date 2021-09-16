terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_app_connector_group" "sgio-vancouver" {
  name = "SGIO-Vancouver"
}

resource "zpa_server_group" "srvg_devops" {
  name = "SGIO DevOps"
  description = "SGIO DevOps"
  enabled = true
  dynamic_discovery = false
  app_connector_groups {
    id = [data.zpa_app_connector_group.sgio-vancouver.id]
  }
  servers {
    id = [zpa_application_server.sgio_jenkins.id]
  }
}

 resource "zpa_segment_group" "sg_devops" {
   name = "SGIO DevOps"
   description = "SGIO DevOps"
   enabled = true
   policy_migrated = true
 }

resource "zpa_application_server" "sgio_jenkins" {
  name                          = "jenkins.securitygeek.io"
  description                   = "jenkins.securitygeek.io"
  address                       = "jenkins.securitygeek.io"
  enabled                       = true
}


// DevOps Browser Access
data "zpa_ba_certificate" "jenkins_ba" {
    name = "jenkins.securitygeek.io"
}


resource "zpa_browser_access" "jenkins_browser_access" {
    name = "jenkins_app"
    description = "jenkins_app"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    tcp_port_ranges = ["80", "80", "8080", "8080"]
    domain_names = ["jenkins.securitygeek.io"]
    segment_group_id = zpa_segment_group.sg_devops.id

    clientless_apps {
        name = "jenkins.securitygeek.io"
        application_protocol = "HTTP"
        application_port = "8080"
        certificate_id = data.zpa_ba_certificate.jenkins_ba.id
        trust_untrusted_cert = true
        enabled = true
        domain = "jenkins.securitygeek.io"
    }
    server_groups {
        id = [
            zpa_server_group.srvg_devops.id
        ]
    }
}