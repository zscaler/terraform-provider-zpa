---
subcategory: "Service Edge Controller"
layout: "zscaler"
page_title: "ZPA: service_edge_controller"
description: |-
  Get information about Service Edge Controller in Zscaler Private Access cloud.
---

# Data Source: zpa_service_edge_controller

Use the **zpa_service_edge_controller** data source to get information about a service edge controller in the Zscaler Private Access cloud. This data source can then be referenced in a Service Edge Group and Provisioning Key.

## Example Usage

```hcl
# ZPA Service Edge Controller Data Source
data "zpa_service_edge_controller" "example" {
  name = "On-Prem-PSE"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service edge controller to be exported.
* `id` - (Optional) The ID of the service edge controllerto be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` - (bool) Whether this Service Edge Controller is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `description` (string) - Description of the App Connector.
* `app_connector_group_name` (Computed) - Expected values: UNKNOWN/ZPN_STATUS_AUTHENTICATED(1)/ZPN_STATUS_DISCONNECTED
* `latitude` - (string) Latitude of the Service Edge Controller. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (string) Longitude of the Service Edge Controller. Integer or decimal. With values in the range of `-180` to `180`
* `location` - (string) Location of the Service Edge Controller.
* `application_start_time` (string)
* `app_connector_group_id` (string)
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
* `provisioning_key_id` (string)
* `provisioning_key_name` (string)
* `platform` (string)
* `previous_version` (string)
* `private_ip` (string)
* `public_ip` (string)
* `sarge_version` (string)
* `enrollment_cert` (string)
* `upgrade_attempt` (string)
* `upgrade_status` (string)
