# ZPA Application Segment resource
resource "zpa_application_segment" "this" {
    name              = "Example"
    description       = "Example"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
    tcp_port_ranges   = ["8080", "8080"]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.this.id
    server_groups {
        id = [ zpa_server_group.this.id]
    }
    depends_on = [ zpa_server_group.this, zpa_segment_group.this]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "this" {
  name            = "Example"
  description     = "Example"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "this" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.this.id ]
  }
  depends_on = [ zpa_app_connector_group.this ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "this" {
  name                          = "Example"
  description                   = "Example"
  enabled                       = true
  city_country                  = "San Jose, CA"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}

# Create PRA Approval Controller
resource "zpa_pra_approval_controller" "this" {
    email_ids = ["jdoe@acme.com"]
    start_time = "Tue, 07 Mar 2024 11:05:30 PST"
    end_time = "Tue, 07 Jun 2024 11:05:30 PST"
    status = "FUTURE"
    applications {
      id = [zpa_application_segment.this.id]
    }
    working_hours {
      days = ["FRI", "MON", "SAT", "SUN", "THU", "TUE", "WED"]
      start_time = "00:10"
      start_time_cron = "0 0 8 ? * MON,TUE,WED,THU,FRI,SAT"
      end_time = "09:15"
      end_time_cron = "0 15 17 ? * MON,TUE,WED,THU,FRI,SAT"
      timezone = "America/Vancouver"
    }
}