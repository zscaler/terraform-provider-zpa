// Retrieve Policy Set ID
data "zpa_policy_type" "inspection_policy" {
    policy_type = "INSPECTION_POLICY"
}

resource "zpa_policy_inspection_rule" "this" {
  name          = "Example"
  description   = "Example"
  action        = "BYPASS_INSPECT"
  operator      = "AND"
  policy_set_id = data.zpa_policy_type.inspection_policy.id

  conditions {
    operator = "OR"
    operands {
      object_type = "APP"
      lhs         = "id"
      rhs         = zpa_application_segment_inspection.this.id
    }
  }
  depends_on = [ zpa_application_segment_inspection.this ]
}

// Create Inspection Application Segment
resource "zpa_application_segment_inspection" "this" {
  name             = "ZPA_Inspection_Example"
  description      = "ZPA_Inspection_Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["80", "80"]
  domain_names     = ["jenkins.example.com"]
  segment_group_id = zpa_segment_group.this.id
  server_groups {
    id = [zpa_server_group.this.id]
  }
  common_apps_dto {
    apps_config {
      name                 = "jenkins.example.com"
      domain               = "jenkins.example.com"
      application_protocol = "HTTP"
      application_port     = "80"
      certificate_id       = data.zpa_ba_certificate.jenkins.id
      enabled              = true
      app_types            = ["INSPECT"]
    }
  }
  depends_on = [ zpa_segment_group.this, zpa_server_group.this ]
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
