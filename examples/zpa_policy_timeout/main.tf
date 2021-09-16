terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}

resource "zpa_application_segment" "all_other_services" {
    name = "All Other Services"
    description = "All Other Services"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    tcp_port_ranges = ["1", "52", "54", "65535"]
    domain_names = ["*.securitygeek.io"]
    segment_group_id = zpa_segment_group.sg_all_other_services.id
    server_groups {
        id = [zpa_server_group.all_other_services.id]
    }
}

 resource "zpa_segment_group" "sg_all_other_services" {
   name = "All Other Services"
   description = "All Other Services"
   enabled = true
   policy_migrated = true
 }

 resource "zpa_server_group" "all_other_services" {
  name = "All Other Services"
  description = "All Other Services"
  enabled = true
  dynamic_discovery = true
  app_connector_groups {
    id = [data.zpa_app_connector_group.sgio-vancouver.id]
  }
}

resource "zpa_policy_timeout" "all_other_services" {
  name                          = "All Other Services"
  description                   = "All Other Services"
  action                        = "RE_AUTH"
  reauth_idle_timeout = "600"
  reauth_timeout = "172800"
  rule_order                     = 1
  operator = "AND"
  policy_set_id = data.zpa_policy_timeout.all.id

    conditions {
    negated = false
    operator = "OR"
    operands {
      name =  "All Other Services"
      object_type = "APP_GROUP"
      lhs = "id"
      rhs = zpa_segment_group.sg_all_other_services.id
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

data "zpa_policy_timeout" "all" {
}

data "zpa_app_connector_group" "sgio-vancouver" {
  name = "SGIO-Vancouver"
}

data "zpa_idp_controller" "sgio_user_okta" {
 name = "SGIO-User-Okta"
}

data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "SGIO-User-Okta"
}