---
subcategory: "App Connector Group"
layout: "zscaler"
page_title: "ZPA: app_connector_group"
description: |-
  Creates and manages ZPA App Connector Groups.
---

# Resource: zpa_app_connector_group

The **zpa_app_connector_group** resource creates a and manages app connector groups in the Zscaler Private Access (ZPA) cloud. This resource can then be associated with the following resoueces: server groups, log receivers and access policies.

## Example Usage

```hcl
# Create a App Connector Group
resource "zpa_app_connector_group" "example" {
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the App Connector Group.
* `description` (Optional) Description of the App Connector Group.
* `enabled` - (Optional) Whether this App Connector Group is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `latitude` - (Required) Latitude of the App Connector Group. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (Required) Longitude of the App Connector Group. Integer or decimal. With values in the range of `-180` to `180`
* `location` - (Required) Location of the App Connector Group.
* `city_country` - (Optional) Whether Double Encryption is enabled or disabled for the app.
* `upgrade_day` - (Optional) App Connectors in this group will attempt to update to a newer version of the software during this specified day
* `upgrade_time_in_secs` - (Optional) App Connectors in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600`. Integer in seconds (i.e., `-66600`). The integer should be greater than or equal to `0` and less than `86400`, in `15` minute intervals
* `override_version_profile` - (Optional) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false` Supported values: `true`, `false`
* `version_profile_id` - (Optional) ID of the version profile. To learn more, see Version Profile Use Cases.
* `version_profile_name` - (Optional)
* `version_profile_visibility_scope` - (Optional)
* `country_code` - (Optional)
* `dns_query_type` - (Optional)
* `geo_location_id` - (Optional)
* `geo_location_id` - (Optional)

## Attributes Reference

* `id` - The ID of the Group Role Assignment.

## Import

App Connector Group can be imported by using `<APP CONNECTOR GROUP ID>` or `<APP CONNECTOR GROUP NAME>`as the import ID.

```shell
terraform import zpa_app_connector_group.example <app_connector_group_id>
```

or

```shell
terraform import zpa_app_connector_group.example <app_connector_group_name>
```
