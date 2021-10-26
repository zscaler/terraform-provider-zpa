---
subcategory: "App Connector Management"
layout: "zpa"
page_title: "ZPA: app_connector_group"
description: |-
  Creates a ZPA App Connector Group.
  
---
# zpa_app_connector_group

The **zpa_app_connector_group** resource creates an app connector group in the Zscaler Private Access cloud. This resource can then be associated with a provisioning key, policy access and connector resources.

## Example Usage

```hcl
# ZPA App Connector Group Data Source
resource "zpa_app_connector_group" "example" {
  name                          = "Example"
  description                   = "Example"
  enabled                       = true
  city_country                  = "New York, CA"
  country_code                  = "US"
  latitude                      = "37.3382082"
  longitude                     = "-121.8863286"
  location                      = "San Jose, CA, USA"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the App Connector Group.
* `latitude` - (Required) Latitude of the App Connector Group. With values in the range of -90 to 90
* `longitude` - (Required) Longitude of the App Connector Group. With values in the range of -180 to 180
* `location` - (Required) Location of the App Connector Group.
* `version_profile_id` - (Optional)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` - (Optional) Whether this App Connector Group is enabled or not. Supported values: `true`, `false`
* `description` - (Optional) Description of the App Connector Group.
* `upgrade_day` - (Optional) App Connectors in this group will attempt to update to a newer version of the software during this specified day. Default value: `SUNDAY`. List of valid days (i.e., Sunday, Monday)
* `dns_query_type` - (Optional) Whether to enable IPv4 or IPv6, or both, for DNS resolution of all applications in the App Connector Group. Default: `IPV4_IPV6`. Supported values: `IPV4_IPV6`, `IPV4`, `IPV6`
* `upgrade_time_in_secs` - (Optional) App Connectors in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600`. Integer in seconds (i.e., -66600). The integer should be greater than or equal to `0` and less than `86400`, in `15 minute` intervals
* `override_version_profile` - (Optional) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false`. Supported values: `true`, `false`
* `city_country` - (Optional)
* `country_code` - (Optional)
* `geolocation_id` - (Optional)
* `siem_appconnector_group` - (Optional)
