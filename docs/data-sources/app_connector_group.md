---
subcategory: "App Connector Management"
layout: "zpa"
page_title: "ZPA: app_connector_group"
description: |-
  Gets a ZPA App Connector Group details.
  
---
# zpa_app_connector_group

The **zpa_app_connector_group** data source provides details about a specific app connector group created in the Zscaler Private Access cloud. This data source can then be referenced in the following resources:

* Create a server group
* Provisioning Key
* Create an access policy routing to control path selection

## Example Usage

```hcl
# ZPA App Connector Group Data Source
data "zpa_app_connector_group" "foo" {
  name = "DataCenter"
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
