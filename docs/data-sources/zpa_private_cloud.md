---
page_title: "zpa_private_cloud Data Source - terraform-provider-zpa"
subcategory: "Private Clouds"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-private-cloud-controllers
  API documentation https://help.zscaler.com/zpa/about-private-cloud-controllers
  Get information about ZPA Private Cloud in Zscaler Private Access cloud.
---

# zpa_private_cloud (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-private-cloud-controllers)
* [API documentation](https://help.zscaler.com/zpa/about-private-cloud-controllers)

The **zpa_private_cloud** data source to get information about a private cloud in the Zscaler Private Access cloud.

## Example Usage

```terraform
# ZPA Private Cloud Data Source
data "zpa_private_cloud" "foo" {
  name = "PrivateCloud01"
}
```

```terraform
# ZPA Private Cloud Data Source
data "zpa_private_cloud" "foo" {
  id = "123456789"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the private cloud to be exported.

### Optional

* `id` - (Optional) The ID of the private cloud to be exported.
* `microtenant_id` - (Optional) Microtenant ID for the private cloud.

### Read-Only

* `description` - (String) Description of the Private Cloud
* `enabled` - (Boolean) Whether this Private Cloud is enabled or not
* `re_enroll_period` - (String) The re-enrollment period for the Private Cloud
* `fire_drill_enabled` - (Boolean) Whether fire drill is enabled for the Private Cloud
* `sitec_preferred` - (Boolean) Whether the Site Controller is preferred
* `remote_lss` - (Boolean) Whether remote Log Streaming Service (LSS) is enabled
* `read_only` - (Boolean) Whether the Private Cloud is read-only
* `zscaler_managed` - (Boolean) Whether the Private Cloud is managed by Zscaler
* `microtenant_name` - (String) Microtenant name for the Private Cloud
* `creation_time` - (String) The time the Private Cloud was created
* `modified_by` - (String) The ID of the user that last modified the Private Cloud
* `modified_time` - (String) The time the Private Cloud was last modified

* `assistant_groups_ids` - (List) The list of Assistant (App Connector) Group IDs associated with the Private Cloud.
    * `id` - (String) The unique identifier of the Assistant Group.
    * `name` - (String) The name of the Assistant Group.
    * `enabled` - (Boolean) Whether the Assistant Group is enabled.

* `site_controller_group_ids` - (List) The list of Site Controller Group IDs associated with the Private Cloud.
    * `id` - (String) The unique identifier of the Site Controller Group.
    * `name` - (String) The name of the Site Controller Group.
    * `enabled` - (Boolean) Whether the Site Controller Group is enabled.

* `siem_ids` - (List) The list of SIEM IDs associated with the Private Cloud.
    * `id` - (String) The unique identifier of the SIEM.
    * `name` - (String) The name of the SIEM.
    * `enabled` - (Boolean) Whether the SIEM is enabled.

* `private_exporter_group_ids` - (List) The list of Private Exporter Group IDs associated with the Private Cloud.
    * `id` - (String) The unique identifier of the Private Exporter Group.
    * `name` - (String) The name of the Private Exporter Group.
    * `enabled` - (Boolean) Whether the Private Exporter Group is enabled.

* `private_broker_group_ids` - (List) The list of Private Broker Group IDs associated with the Private Cloud.
    * `id` - (String) The unique identifier of the Private Broker Group.
    * `name` - (String) The name of the Private Broker Group.
    * `enabled` - (Boolean) Whether the Private Broker Group is enabled.

* `zpn_fire_drill_site` - (List) The fire drill site configuration for the Private Cloud.
    * `id` - (String) The unique identifier of the fire drill site.
    * `microtenant_id` - (String) The microtenant ID of the fire drill site.
    * `microtenant_name` - (String) The microtenant name of the fire drill site.
    * `fire_drill_interval` - (String) The fire drill interval.
    * `fire_drill_interval_time_unit` - (String) The fire drill interval time unit.
    * `creation_time` - (String) The time the fire drill site was created.
    * `modified_by` - (String) The ID of the user that last modified the fire drill site.
    * `modified_time` - (String) The time the fire drill site was last modified.
