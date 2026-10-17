// Create Application Segment with Tags
resource "zpa_application_segment" "crm_application" {
  name             = "CRM Application"
  description      = "CRM Application"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["80", "80"]
  domain_names     = ["crm.example.com"]
  segment_group_id = zpa_segment_group.crm_app_group.id
  server_groups {
    id = [zpa_server_group.crm_servers.id]
  }

  tags {
    tag_key_id   = zpa_tag_key.environment.id
    tag_value_id = zpa_tag_key.environment.tag_values[0].id
  }
}

// Create Tag Namespace
resource "zpa_tag_namespace" "this" {
  name        = "CRM Namespace"
  description = "CRM Namespace"
  enabled     = true
}

// Create Tag Key with Tag Values
resource "zpa_tag_key" "environment" {
  name         = "Environment"
  description  = "Environment"
  enabled      = true
  namespace_id = zpa_tag_namespace.this.id

  tag_values {
    name = "Production"
  }

  tag_values {
    name = "Staging"
  }
}

// Create Server Group
resource "zpa_server_group" "crm_servers" {
  name              = "CRM Servers"
  description       = "CRM Servers"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [data.zpa_app_connector_group.this.id]
  }
  servers {
    id = [zpa_application_server.crm_app_server.id]
  }
}

// Create Application Server
resource "zpa_application_server" "crm_app_server" {
  name        = "CRM App Server"
  description = "CRM App Server"
  address     = "crm.example.com"
  enabled     = true
}

// Create Segment Group
resource "zpa_segment_group" "crm_app_group" {
  name        = "CRM App group"
  description = "CRM App group"
  enabled     = true
}

// Retrieve App Connector Group
data "zpa_app_connector_group" "this" {
  name = "Example"
}
