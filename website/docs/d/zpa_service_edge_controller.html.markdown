---
layout: "zscaler"
page_title: "Zscaler Private Access (ZPA): service_edge_controller"
sidebar_current: "docs-datasource-zpa-service_edge_controller"
description: |-
  Get information about Service Edge Controller in Zscaler Private Access cloud.
---

# zpa_service_edge_controller

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

* `enabled` - (Computed) Whether this Service Edge Controller is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `description` (Computed) - Description of the App Connector.
* `app_connector_group_name` (Computed) - Expected values: UNKNOWN/ZPN_STATUS_AUTHENTICATED(1)/ZPN_STATUS_DISCONNECTED
* `latitude` - (Computed) Latitude of the Service Edge Controller. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (Computed) Longitude of the Service Edge Controller. Integer or decimal. With values in the range of `-180` to `180`
* `location` - (Computed) Location of the Service Edge Controller.
* `application_start_time` (Computed)
* `app_connector_group_id` (Computed)
* `control_channel_status` (Computed)
* `creation_time` (Computed)
* `modified_by` (Computed)
* `modified_time` (Computed)
* `ctrl_broker_name` (Computed)
* `current_version` (Computed)
* `expected_upgrade_time` (Computed)
* `expected_version` (Computed)
* `figerprint` (Computed)
* `ip_acl` (Computed)
* `issued_cert_id` (Computed)
* `last_broker_connect_time` (Computed)
* `last_broker_connect_time_duration` (Computed)
* `last_broker_disconnect_time` (Computed)
* `last_broker_disconnect_time_duration` (Computed)
* `last_upgrade_time` (Computed)
* `provisioning_key_id` (Computed)
* `provisioning_key_name` (Computed)
* `platform` (Computed)
* `previous_version` (Computed)
* `private_ip` (Computed)
* `public_ip` (Computed)
* `sarge_version` (Computed)
* `enrollment_cert` (Computed)
* `upgrade_attempt` (Computed)
* `upgrade_status` (Computed)
