---
page_title: "zpa_private_cloud_controller Data Source - terraform-provider-zpa"
subcategory: "Private Cloud Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-private-cloud-controllers
  API documentation https://help.zscaler.com/zpa/about-private-cloud-controllers
  Get information about ZPA Private Cloud Controller in Zscaler Private Access cloud.
---

# zpa_private_cloud_controller (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-private-cloud-controllers)
* [API documentation](https://help.zscaler.com/zpa/about-private-cloud-controllers)

The **zpa_private_cloud_controller** data source to get information about a private cloud controller in the Zscaler Private Access cloud.

## Example Usage - Search by Name

```terraform
data "zpa_private_cloud_controller" "foo" {
  name = "DataCenter"
}
```

## Example Usage - Search by ID

```terraform
data "zpa_private_cloud_controller" "foo" {
  id = "123456789"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the private cloud controller to be exported.

### Optional

* `id` - (Optional) The ID of the private cloud controller to be exported.
* `microtenant_id` - (Optional) Microtenant ID for the private cloud controller.

### Read-Only

* `application_start_time` - (String) Application start time of the Private Cloud Controller
* `control_channel_status` - (String) Control channel status of the Private Cloud Controller
* `creation_time` - (String) Creation time of the Private Cloud Controller
* `ctrl_broker_name` - (String) Control broker name of the Private Cloud Controller
* `current_version` - (String) Current version of the Private Cloud Controller
* `description` - (String) Description of the Private Cloud Controller
* `enabled` - (Boolean) Whether this Private Cloud Controller is enabled or not
* `expected_sarge_version` - (String) Expected Sarge version of the Private Cloud Controller
* `expected_upgrade_time` - (String) Expected upgrade time of the Private Cloud Controller
* `expected_version` - (String) Expected version of the Private Cloud Controller
* `fingerprint` - (String) Fingerprint of the Private Cloud Controller
* `ip_acl` - (List of String) IP ACL list of the Private Cloud Controller
* `issued_cert_id` - (String) Issued certificate ID of the Private Cloud Controller
* `last_broker_connect_time` - (String) Last broker connect time of the Private Cloud Controller
* `last_broker_connect_time_duration` - (String) Last broker connect time duration of the Private Cloud Controller
* `last_broker_disconnect_time` - (String) Last broker disconnect time of the Private Cloud Controller
* `last_broker_disconnect_time_duration` - (String) Last broker disconnect time duration of the Private Cloud Controller
* `last_os_upgrade_time` - (String) Last OS upgrade time of the Private Cloud Controller
* `last_sarge_upgrade_time` - (String) Last Sarge upgrade time of the Private Cloud Controller
* `last_upgrade_time` - (String) Last upgrade time of the Private Cloud Controller
* `latitude` - (String) Latitude of the Private Cloud Controller
* `listen_ips` - (List of String) Listen IPs of the Private Cloud Controller
* `location` - (String) Location of the Private Cloud Controller
* `longitude` - (String) Longitude of the Private Cloud Controller
* `master_last_sync_time` - (String) Master last sync time of the Private Cloud Controller
* `modified_by` - (String) Modified by information for the Private Cloud Controller
* `modified_time` - (String) Modified time of the Private Cloud Controller
* `provisioning_key_id` - (String) Provisioning key ID of the Private Cloud Controller
* `provisioning_key_name` - (String) Provisioning key name of the Private Cloud Controller
* `os_upgrade_enabled` - (Boolean) Whether OS upgrade is enabled for the Private Cloud Controller
* `os_upgrade_status` - (String) OS upgrade status of the Private Cloud Controller
* `platform` - (String) Platform of the Private Cloud Controller
* `platform_detail` - (String) Platform detail of the Private Cloud Controller
* `platform_version` - (String) Platform version of the Private Cloud Controller
* `previous_version` - (String) Previous version of the Private Cloud Controller
* `private_ip` - (String) Private IP of the Private Cloud Controller
* `public_ip` - (String) Public IP of the Private Cloud Controller
* `publish_ips` - (List of String) Publish IPs of the Private Cloud Controller
* `read_only` - (Boolean) Whether the Private Cloud Controller is read-only
* `restriction_type` - (String) Restriction type of the Private Cloud Controller
* `runtime` - (String) Runtime of the Private Cloud Controller
* `sarge_upgrade_attempt` - (String) Sarge upgrade attempt of the Private Cloud Controller
* `sarge_upgrade_status` - (String) Sarge upgrade status of the Private Cloud Controller
* `sarge_version` - (String) Sarge version of the Private Cloud Controller
* `microtenant_name` - (String) Microtenant name for the Private Cloud Controller
* `shard_last_sync_time` - (String) Shard last sync time of the Private Cloud Controller
* `enrollment_cert` - (Map of String) Enrollment certificate of the Private Cloud Controller
* `private_cloud_controller_group_id` - (String) Private Cloud Controller group ID
* `private_cloud_controller_group_name` - (String) Private Cloud Controller group name
* `private_cloud_controller_version` - (Map of String) Private Cloud Controller version information
* `site_sp_dns_name` - (String) Site SP DNS name of the Private Cloud Controller
* `upgrade_attempt` - (String) Upgrade attempt of the Private Cloud Controller
* `upgrade_status` - (String) Upgrade status of the Private Cloud Controller
* `userdb_last_sync_time` - (String) User database last sync time of the Private Cloud Controller
* `zpn_sub_module_upgrade_list` - (List of String) ZPN sub-module upgrade list of the Private Cloud Controller
* `zscaler_managed` - (Boolean) Whether the Private Cloud Controller is managed by Zscaler
