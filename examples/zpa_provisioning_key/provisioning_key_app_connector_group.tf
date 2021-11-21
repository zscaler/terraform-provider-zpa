// Create Provisioning Key for App Connector Group
resource "zpa_provisioning_key" "nyc_provisioning_key" {
  name             = "New York Provisioning Key"
  association_type = "CONNECTOR_GRP"
  max_usage        = "10"
  enrollment_cert_id = data.zpa_enrollment_cert.connector.id
  zcomponent_id = zpa_app_connector_group.nyc_connector_group.id
}

resource "zpa_app_connector_group" "nyc_connector_group" {
  name                          = "App Connector Group New York"
  description                   = "App Connector Group New York"
  enabled                       = true
  city_country                  = "New York, NY"
  country_code                  = "USA"
  latitude                      = "49.1041779"
  longitude                     = "-122.6603519"
  location                      = "New York, NY, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}

data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}
*/