---
page_title: "zpa_service_edge_controller Data Source - terraform-provider-zpa"
subcategory: "Service Edge Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-zpa-private-service-edges
  API documentation https://help.zscaler.com/zpa/managing-zpa-private-service-edges-using-api
  Get information about Service Edge Controller in Zscaler Private Access cloud.
---

# Data Source: zpa_service_edge_controller

* [Official documentation](https://help.zscaler.com/zpa/about-zpa-private-service-edges)
* [API documentation](https://help.zscaler.com/zpa/managing-zpa-private-service-edges-using-api)

Use the **zpa_service_edge_controller** data source to get information about a service edge controller in the Zscaler Private Access cloud. This data source can then be referenced in a Service Edge Group and Provisioning Key.

## Example Usage

```terraform
# ZPA Service Edge Controller Data Source
data "zpa_service_edge_controller" "example" {
  name = "On-Prem-PSE"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the service edge controller to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the service edge controllerto be exported.
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
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
