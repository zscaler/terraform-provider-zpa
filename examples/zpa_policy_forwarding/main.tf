terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

data "zpa_policy_forwarding" "all" {
}

data "zpa_idp_controller" "sgio_user_okta" {
 name = "SGIO-User-Okta"
}

// Okta IDP SCIM Groups
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "SGIO-User-Okta"
}

// Application Connector Groups
data "zpa_app_connector_group" "sgio-vancouver" {
  name = "SGIO-Vancouver"
}

resource "zpa_policy_forwarding" "example" {
  name                          = "example"
  description                   = "example"
  action                        = "INTERCEPT_ACCESSIBLE"
  operator = "AND"
  policy_set_id = data.zpa_policy_forwarding.all.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      name =  "SGIO DevOps Servers"
      object_type = "APP"
      lhs = "id"
      rhs = zpa_application_segment.as_sgio_devops.id
    }
  }
  conditions {
     negated = false
     operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.sgio_user_okta.id
      rhs = data.zpa_scim_groups.engineering.id
      idp_id = data.zpa_idp_controller.sgio_user_okta.id
    }
  }
}

resource "zpa_application_segment" "as_sgio_devops" {
    name = "SGIO DevOps Servers"
    description = "SGIO DevOps Servers"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    tcp_port_ranges = ["8080", "8080"]
    domain_names = ["jenkins.securitygeek.io"]
    segment_group_id = zpa_segment_group.sg_sgio_devops.id
    server_groups {
        id = [zpa_server_group.sgio_devops_servers.id]
    }
}

  resource "zpa_segment_group" "sg_sgio_devops" {
   name = "SGIO DevOps Servers"
   description = "SGIO DevOps Servers"
   enabled = true
   policy_migrated = true
 }

 resource "zpa_server_group" "sgio_devops_servers" {
  name = "SGIO DevOps Servers"
  description = "SGIO DevOps Servers"
  enabled = true
  dynamic_discovery = false
    servers {
    id = [
      zpa_application_server.jenkins.id,
    ]
  }
  app_connector_groups {
    id = [data.zpa_app_connector_group.sgio-vancouver.id]
  }
}

resource "zpa_application_server" "jenkins" {
  name                          = "jenkins.securitygeek.io"
  description                   = "jenkins.securitygeek.io"
  address                       = "jenkins.securitygeek.io"
  enabled                       = true
}