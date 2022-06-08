---
subcategory: "Service Edge Group"
layout: "zscaler"
page_title: "ZPA: service_edge_group"
description: |-
  Get information about ZPA Service Edge Group in Zscaler Private Access cloud.
---

# Data Source: zpa_service_edge_group

Use the **zpa_service_edge_group** data source to get information about a service edge group in the Zscaler Private Access cloud. This data source can then be referenced in an App Connector Group. This data source can then be referenced in the following resources:

* Create a server group
* Provisioning Key
* Access policy rule

## Example Usage

```hcl
# ZPA Service Edge Group Data Source by name
data "zpa_service_edge_group" "foo" {
  name = "DataCenter"
}
```

```hcl
# ZPA Service Edge Group Data Source by ID
data "zpa_service_edge_group" "foo" {
  id = "123456789"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service edge group to be exported.
* `id` - (Optional) The ID of the service edge group to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` (string) Description of the Service Edge Group.
* `enabled` - (bool) Whether this App Connector Group is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `latitude` - (string) Latitude of the Service Edge Group. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (string) Longitude of the Service Edge Group.Integer or decimal. With values in the range of `-180` to `180`
* `location` - (string) Location of the Service Edge Group.
* `city_country` - (string) Whether Double Encryption is enabled or disabled for the app.
* `upgrade_day` - (string) App Connectors in this group will attempt to update to a newer version of the software during this specified day
* `upgrade_time_in_secs` - (string) App Connectors in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600`. Integer in seconds (i.e., `-66600`). The integer should be greater than or equal to `0` and less than `86400`, in `15` minute intervals
* `override_version_profile` - (bool) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false` Supported values: `true`, `false`
* `version_profile_id` - (string) ID of the version profile. To learn more, see Version Profile Use Cases.
* `version_profile_name` - (string)
* `version_profile_visibility_scope` - (string)
* `country_code` - (string)
* `dns_query_type` - (string)
* `geo_location_id` - (string)
