resource "zpa_provisioning_key" "nyc_provisioning_key" {
  name             = "New York Provisioning Key"
  association_type = "SERVICE_EDGE_GRP"
  max_usage        = "10"
  enrollment_cert_id = data.zpa_enrollment_cert.service_edge.id
  zcomponent_id = zpa_service_edge_group.nyc_service_edge_group.id
}

resource "zpa_service_edge_group" "nyc_service_edge_group" {
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