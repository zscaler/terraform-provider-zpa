resource "zpa_private_cloud" "this" {
  name               = "PrivateCloud01"
  description        = "Example private cloud"
  enabled            = true
  re_enroll_period   = "180"
  fire_drill_enabled = false
  sitec_preferred    = false
  remote_lss         = true

  # Private Cloud Group Controller Association
  site_controller_group_ids {
    id = [zpa_private_cloud_group.this.id]
  }
  # App Connector Group Controller Association
  assistant_groups_ids {
    id = [zpa_app_connector_group.this.id]
  }
  # Service Edge Group Controller Association
  private_broker_group_ids {
    id = [zpa_service_edge_group.this.id]
  }
}

# Private Cloud Group
resource "zpa_private_cloud_group" "this" {
  name                     = "PrivateCloudGroup01"
  description              = "Example private cloud group"
  enabled                  = true
  country_code             = "US"
  city_country             = "San Jose, US"
  latitude                 = "37.33874"
  longitude                = "-121.8852525"
  location                 = "San Jose, CA, USA"
  upgrade_day              = "SUNDAY"
  upgrade_time_in_secs     = "66600"
  version_profile_id       = "0"
  override_version_profile = true
  is_public                = "TRUE"
}

# App Connector Group
resource "zpa_app_connector_group" "this" {
  name                          = "AppConnectorGroup01"
  description                   = "AppConnectorGroup01"
  enabled                       = true
  city_country                  = "San Jose, US"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  version_profile_id            = "0"
  override_version_profile      = true
  dns_query_type                = "IPV4_IPV6"
  use_in_dr_mode                = false
}

# Service Edge Group
resource "zpa_service_edge_group" "this" {
  name                        = "ServiceEdgeGroup01"
  description                 = "ServiceEdgeGroup01"
  enabled                     = true
  is_public                   = false
  upgrade_day                 = "SUNDAY"
  city_country                = "San Jose, US"
  country_code                = "US"
  latitude                    = "37.338"
  longitude                   = "-121.8863"
  location                    = "San Jose, CA, US"
  upgrade_time_in_secs        = "66600"
  version_profile_id            = "0"
  override_version_profile      = true
}