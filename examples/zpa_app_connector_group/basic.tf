resource "zpa_app_connector_group" "example" {
  name                          = "Example"
  description                   = "Example"
  enabled                       = true
  city_country                  = "San Jose, US"
  country_code                  = "US"
  latitude                      = "37.3382082"
  longitude                     = "-121.8863286"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_name          = "New Release"
  dns_query_type                = "IPV4"
  pra_enabled                   = true
  waf_disabled                  = false
}