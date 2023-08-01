---
subcategory: "App Connector Controller"
layout: "zscaler"
page_title: "ZPA: app_connector_controller"
description: |-
  Get information about ZPA App Connector in Zscaler Private Access cloud.
---

# Data Source: app_connector_controller

Use the **zpa_app_connector_controller** data source to get information about a app connector created in the Zscaler Private Access cloud. This data source can then be referenced in an App Connector Group.

## Example Usage

```hcl
# ZPA App Connector Data Source
data "zpa_app_connector" "example" {
  name = "AWS-VPC100-App-Connector"
}
```

```hcl
# ZPA App Connector Data Source
data "zpa_app_connector" "example" {
  id = "123456789"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the App Connector Group.
* `id` - (Optional) Name of the App Connector Group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

The following values are ignored in PUT/POST calls. Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.

* `description` (Computed) - Description of the App Connector.
* `app_connector_group_name` (Computed) - Expected values: UNKNOWN/ZPN_STATUS_AUTHENTICATED(1)/ZPN_STATUS_DISCONNECTED
* `latitude` (Computed) - Latitude of the App Connector. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (Computed) - Longitude of the App Connector. Integer or decimal. With values in the range of `-180` to `180`
* `enabled` - (Computed) - Whether this App Connector is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `location` (Computed) - Location of the App Connector.
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
* `microtenant_id` (Computed)
* `microtenant_name` (Computed)

:warning: Notice that certificate and public_keys are omitted from the output.
