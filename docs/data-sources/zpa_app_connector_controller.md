---
page_title: "zpa_app_connector_controller Data Source - terraform-provider-zpa"
subcategory: "App Connector Controller"
  Official documentation https://help.zscaler.com/zpa/about-connectors
  documentation https://help.zscaler.com/zpa/managing-app-connectors-using-api
description: |-
  Get information about ZPA App Connector in Zscaler Private Access cloud.
---

# zpa_app_connector_controller (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-connectors)
* [API documentation](https://help.zscaler.com/zpa/managing-app-connectors-using-api)

Use the **zpa_app_connector_controller** data source to get information about a app connector created in the Zscaler Private Access cloud. This data source can then be referenced in an App Connector Group.

## Example Usage

```terraform
# ZPA App Connector Data Source
data "zpa_app_connector_controller" "example" {
  name = "AWS-VPC100-App-Connector"
}
```

```terraform
# ZPA App Connector Data Source
data "zpa_app_connector_controller" "example" {
  id = "123456789"
}
```

## Schema

### Required

The following values are returned:

- `name` - (String) Name of the App Connector Group.
- `id` - (String) Name of the App Connector Group.

### Read-Only

In addition to all arguments above, the following attributes are exported:

The following values are ignored in PUT/POST calls. Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.

- `description` (String) - Description of the App Connector.
- `app_connector_group_name` (String) - Expected values: `UNKNOWN/ZPN_STATUS_AUTHENTICATED(1)` or `ZPN_STATUS_DISCONNECTED`
- `latitude` (String) - Latitude of the App Connector. Integer or decimal. With values in the range of `-90` to `90`
- `longitude` - (String) - Longitude of the App Connector. Integer or decimal. With values in the range of `-180` to `180`
- `enabled` - (String) - Whether this App Connector is enabled or not. Default value: `true`. Supported values: `true`, `false`
- `location` (String) - Location of the App Connector.
- `application_start_time` (String) The start time of the App Connector.
- `app_connector_group_id` (String) The unique identifier of the App Connector Group.
- `control_channel_status` (String) The status of the control channel.
- `creation_time` (String) The time the App Connector is created.
- `modified_by` (String) The unique identifier of the tenant who modified the App Connector.
- `modified_time` (String) The time the App Connector is modified.
- `ctrl_broker_name` (String) The name of the Control Public Service Edge.
- `current_version` (String) The current version of the App Connector. 
- `expected_upgrade_time` (String) The expected upgrade time for the App Connector. 
- `expected_version` (String) The expected version of the App Connector.
- `figerprint` (String) The hardware fingerprint associated with the App Connector.
- `ip_acl` (String) The IP Access List (IP ACL) to allow App Connectors on a specific IP or subnet.
- `issued_cert_id` (String) The unique identifier of the issued certificate.
- `last_broker_connect_time` (String) The time the ZPA Public Service Edge last connected.
- `last_broker_connect_time_duration` (String) The duration of time when the ZPA Public Service Edge last connected.
- `last_broker_disconnect_time` (String) The time the ZPA Public Service Edge last disconnected.
- `last_broker_disconnect_time_duration` (String) The duration of time when the ZPA Public Service Edge last disconnected.
- `last_upgrade_time` (String) The time the App Connector last upgraded.
- `provisioning_key_id` (String) The unique identifier of the provisioning key.
- `provisioning_key_name` (String) The name of the provisioning key.
- `platform` (String) The host OS the App Connector is deployed on.
- `platform_detail` (String) The platform the App Connector is deployed on
- `previous_version` (String) The previous version of the App Connector.
- `private_ip` (String) The private IP of the App Connector.
- `public_ip` (String) The public IP of the App Connector.
- `runtime_os`(String) The run time OS on which the App Connector is running.
- `sarge_version` (String) The manager version of the App Connector.
- `enrollment_cert` (String) The enrollment certificate for the App Connector.
- `upgrade_attempt` (String) The number of attempts the App Connector takes to upgrade. 
- `upgrade_status` (String) The status of the App Connector upgrade.
- `microtenant_id` (string) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as `0` when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.
- `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

~> :warning: Notice that certificate and public_keys are omitted from the output.
