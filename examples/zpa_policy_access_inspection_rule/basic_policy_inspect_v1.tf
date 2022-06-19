// Retrieve Policy Set ID
data "zpa_policy_type" "inspection_policy" {
  policy_type = "INSPECTION_POLICY"
}

resource "zpa_policy_inspection_rule" "this" {
  name                      = "Example"
  description               = "Example"
  action                    = "INSPECT"
  rule_order                = 1
  operator                  = "AND"
  policy_set_id             = data.zpa_policy_type.inspection_policy.id
  zpn_inspection_profile_id = zpa_inspection_profile.this.id
  conditions {
    operator = "OR"
    operands {
      object_type = "APP"
      lhs         = "id"
      rhs         = zpa_application_segment_inspection.this.id
    }
  }
}

data "zpa_inspection_predefined_controls" "this" {
  name    = "Failed to parse request body"
  version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_all_predefined_controls" "default_predefined_controls" {
  version    = "OWASP_CRS/3.3.0"
  group_name = "preprocessors"
}

resource "zpa_inspection_profile" "this" {
  name           = "Example"
  description    = "Example"
  paranoia_level = "1"
  dynamic "predefined_controls" {
    for_each = data.zpa_inspection_all_predefined_controls.default_predefined_controls.list
    content {
      id           = predefined_controls.value.id
      action       = predefined_controls.value.action == "" ? predefined_controls.value.default_action : predefined_controls.value.action
      action_value = predefined_controls.value.action_value
    }
  }
  predefined_controls {
    id     = data.zpa_inspection_predefined_controls.this.id
    action = "BLOCK"
  }
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
      app_types            = ["INSPECT"]
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
