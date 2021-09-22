// Create Policy Forwarding Rule
resource "zpa_policy_forwarding_rule" "crm_application_rule" {
  name                          = "CRM Application"
  description                   = "CRM Application"
  action                        = "BYPASS"
  operator = "AND"
  policy_set_id = data.zpa_global_policy_forwarding.policyset.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "APP"
      lhs = "id"
      rhs = [zpa_application_segment.crm_application.id]
    }
  }
  conditions {
     negated = false
     operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [data.zpa_scim_groups.engineering.id]
    }
  }
}

// Create Application Segment
resource "zpa_application_segment" "crm_application" {
    name = "CRM Application"
    description = "CRM Application"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    is_cname_enabled = true
    tcp_port_ranges = ["80", "80"]
    domain_names = ["crm.example.com"]
    segment_group_id = zpa_segment_group.crm_app_group.id
    server_groups {
        id = [ zpa_server_group.crm_servers.id ]
    }
}

// Create Server Group
resource "zpa_server_group" "crm_servers" {
  name = "CRM Servers"
  description = "CRM Servers"
  enabled = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ data.zpa_app_connector_group.dc_connector_group.id ]
  }
  servers {
    id = [ zpa_application_server.crm_app_server.id ]
  }
}

// Create Application Server
resource "zpa_application_server" "crm_app_server" {
  name                          = "CRM App Server"
  description                   = "CRM App Server"
  address                       = "crm.example.com"
  enabled                       = true
}

// Create Segment Group
resource "zpa_segment_group" "crm_app_group" {
  name = "CRM App group"
  description = "CRM App group"
  enabled = true
  policy_migrated = true
}

// Retrieve App Connector Group
data "zpa_app_connector_group" "dc_connector_group" {
  name = "DC Connector Group"
}

data "zpa_global_policy_forwarding" "policyset" {
}

data "zpa_idp_controller" "idp_name" {
 name = "IdP-Name"
}

// Okta IDP SCIM Groups
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP-Name"
}