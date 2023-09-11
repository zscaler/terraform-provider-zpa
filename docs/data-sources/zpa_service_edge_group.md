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
* `city_country` - (string) Whether Double Encryption is enabled or disabled for the app.
* `country_code` - (string)
* `creation_time` - (string)
* `geo_location_id` - (string)
* `is_public` - (string)
* `latitude` - (string) Latitude of the Service Edge Group. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (string) Longitude of the Service Edge Group.Integer or decimal. With values in the range of `-180` to `180`
* `location` - (string) Location of the Service Edge Group.
* `modified_by` - (string)
* `modified_time` - (string)
* `upgrade_day` - (string) App Connectors in this group will attempt to update to a newer version of the software during this specified day
* `upgrade_time_in_secs` - (string) App Connectors in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600`. Integer in seconds (i.e., `-66600`). The integer should be greater than or equal to `0` and less than `86400`, in `15` minute intervals
* `override_version_profile` - (bool) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false` Supported values: `true`, `false`
* `version_profile_id` - (String) ID of the version profile.
  Exported values are:
  * ``0`` = ``Default``
  * ``1`` = ``Previous Default``
  * ``2`` = ``New Release``
* `version_profile_name` - (String)
  Exported values are:
  * ``Default`` = ``0``
  * ``Previous Default`` = ``1``
  * ``New Release`` = ``2``
* `version_profile_visibility_scope` - (string)
  Exported values are:
  * ``ALL``
  * ``NONE``
  * ``CUSTOM``
  * `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
  * `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

* `service_edges` - (string)
  * `name` (string)
  * `application_start_time` (string)
  * `service_edge_group_id` (string)
  * `service_edge_group_name` (string)
  * `control_channel_status` (string)
  * `creation_time` (string)
  * `modified_by` (string)
  * `modified_time` (string)
  * `ctrl_broker_name` (string)
  * `current_version` (string)
  * `expected_upgrade_time` (string)
  * `expected_version` (string)
  * `figerprint` (string)
  * `ip_acl` (string)
  * `issued_cert_id` (string)
  * `last_broker_connect_time` (string)
  * `last_broker_connect_time_duration` (string)
  * `last_broker_disconnect_time` (string)
  * `last_broker_disconnect_time_duration` (string)
  * `last_upgrade_time` (string)
  * `latitude` (string)
  * `listen_ips` (string)
  * `location` (string)
  * `longitude` (string)
  * `provisioning_key_id` (string)
  * `provisioning_key_name` (string)
  * `platform` (string)
  * `previous_version` (string)
  * `private_ip` (string)
  * `public_ip` (string)
  * `publish_ips` (string)
  * `sarge_version` (string)
  * `enrollment_cert` (string)
  * `upgrade_attempt` (string)
  * `upgrade_status` (string)
  * `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
  * `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

* `trusted_networks` - (string)
  * `creation_time` (string)
  * `domain` (string)
  * `id` (string)
  * `master_customer_id` (string)
  * `modified_by` (string)
  * `modified_time` (string)
  * `name` (string)
  * `network_id` (string)
  * `zscaler_cloud` (string)

:warning: Notice that certificate and public_keys are omitted from the output.
