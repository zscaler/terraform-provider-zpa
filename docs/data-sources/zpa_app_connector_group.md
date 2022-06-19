---
subcategory: "App Connector Group"
layout: "zscaler"
page_title: "ZPA: app_connector_group"
description: |-
  Get information about ZPA App Connector Group in Zscaler Private Access cloud.
---

# Data Source: zpa_app_connector_group

Use the **zpa_app_connector_group** data source to get information about a app connector group in the Zscaler Private Access cloud. This data source can then be referenced in an App Connector Group. This data source can then be referenced in the following resources:

* Create a server group
* Provisioning Key
* Access policy rule

## Example Usage

```hcl
# ZPA App Connector Group Data Source
data "zpa_app_connector_group" "foo" {
  name = "DataCenter"
}
```

```hcl
# ZPA App Connector Group Data Source
data "zpa_app_connector_group" "foo" {
  id = "123456789"
}
```

## Argument Reference

In addition to all arguments above, the following attributes are exported:

* `name` - (Required) Name of the App Connector Group.
* `id` - (Required) ID of the App Connector Group.

## Attribute Reference

The following attributes are exported:

* `description` (String) Description of the App Connector Group.
* `enabled` - (String) Whether this App Connector Group is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `latitude` - (String) Latitude of the App Connector Group. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (String) Longitude of the App Connector Group. Integer or decimal. With values in the range of `-180` to `180`
* `location` - (String) Location of the App Connector Group.
* `city_country` - (String) Whether Double Encryption is enabled or disabled for the app.
* `upgrade_day` - (String) App Connectors in this group will attempt to update to a newer version of the software during this specified day
* `upgrade_time_in_secs` - (String) App Connectors in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600`. Integer in seconds (i.e., `-66600`). The integer should be greater than or equal to `0` and less than `86400`, in `15` minute intervals
* `override_version_profile` - (bool) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false` Supported values: `true`, `false`
* `version_profile_id` - (String) ID of the version profile. To learn more, see Version Profile Use Cases.
* `version_profile_name` - (String)
* `version_profile_visibility_scope` - (String)
* `country_code` - (String)
* `dns_query_type` - (String)
* `geo_location_id` - (String)
