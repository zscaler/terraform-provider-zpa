resource "zpa_private_cloud_group" "this" {
  name                     = "PrivateCloudGroup01"
  description              = "Example private cloud group"
  enabled                  = true
  city_country             = "San Jose, US"
  latitude                 = "37.33874"
  longitude                = "-121.8852525"
  location                 = "San Jose, CA, USA"
  upgrade_day              = "SUNDAY"
  upgrade_time_in_secs     = "66600"
  site_id                  = "72058304855088543"
  version_profile_id       = "0"
  override_version_profile = true
  is_public                = "TRUE"
}
