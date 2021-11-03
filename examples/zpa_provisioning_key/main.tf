// Create Provisioning Key for Service Edge Group
resource "zpa_provisioning_key" "usa_provisioning_key" {
  name             = "AWS Provisioning Key"
  association_type = "SERVICE_EDGE_GRP"
  max_usage        = "10"
  enrollment_cert_id = data.zpa_enrollment_cert.service_edge.id
  zcomponent_id = zpa_service_edge_group.service_edge_group_nyc.id
}

resource "zpa_service_edge_group" "service_edge_group_nyc" {
  name                 = "Service Edge Group New York"
  description          = "Service Edge Group in New York"
  upgrade_day          = "SUNDAY"
  upgrade_time_in_secs = "66600"
  latitude             = "40.7128"
  longitude            = "-73.935242"
  location             = "New York, NY, USA"
  version_profile_id   = "0"
}

data "zpa_enrollment_cert" "service_edge" {
    name = "Service Edge"
}

// Create Provisioning Key for App Connector Group
resource "zpa_provisioning_key" "canada_provisioning_key" {
  name             = "Canada Provisioning Key"
  association_type = "CONNECTOR_GRP"
  max_usage        = "10"
  enrollment_cert_id = data.zpa_enrollment_cert.connector.id
  zcomponent_id = zpa_app_connector_group.canada_connector_group.id
}

resource "zpa_app_connector_group" "canada_connector_group" {
  name                          = "Canada Connector Group"
  description                   = "Canada Connector Group"
  enabled                       = true
  city_country                  = "Langley, CA"
  country_code                  = "CA"
  latitude                      = "49.1041779"
  longitude                     = "-122.6603519"
  location                      = "Langley City, BC, Canada"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}

data "zpa_enrollment_cert" "connector" {
    name = "Connector"
}