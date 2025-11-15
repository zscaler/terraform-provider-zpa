---
page_title: "zpa_segment_group Data Source - terraform-provider-zpa"
subcategory: "Segment Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-segment-groups
  API documentation https://help.zscaler.com/zpa/configuring-segment-groups-using-api
  Get information about Segment Groups in Zscaler Private Access cloud.
---

# zpa_segment_group (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-segment-groups)
* [API documentation](https://help.zscaler.com/zpa/configuring-segment-groups-using-api)

Use the **zpa_segment_group** data source to get information about a machine group created in the Zscaler Private Access cloud. This data source can then be referenced in an application segment or Access Policy rule.

## Zenith Community - ZPA Segment Group

[![ZPA Terraform provider Video Series Ep6 - Segment Group](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_segment_groups.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEfCAI/video-zpa-terraform-provider-video-series-ep6-zpa-segment-group)

## Example Usage

```terraform
# ZPA Server Group Data Source
data "zpa_segment_group" "example" {
 name = "segment_group_name"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the segment group to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the segment group to be exported.
* `config_space` - (string)
* `creation_time` - (string)
* `description` - (string)
* `enabled` - (bool)
* `policy_migrated` - (bool)
* `tcp_keep_alive_enabled` - (string)
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

* `applications` - (Computed)
  * `bypass_type` - (string)
  * `config_space` - (string)
  * `default_idle_timeout` - (string)
  * `default_max_age` - (string)
  * `description` - (string)
  * `domain_name` - (string)
  * `domain_names`  - (string)
  * `double_encrypt` - (string)
  * `enabled` - (bool)
  * `health_check_type` - (string)
  * `id` - (string)
  * `ip_anchored` - (bool)
  * `name` - (string)
  * `passive_health_enabled` - (bool)
  * `tcp_port_ranges` - (string)
  * `tcp_ports_in`  - (string)
  * `udp_port_ranges` - (string)
  * `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
  * `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

* `server_groups` - (Computed)
  * `config_space` - (string)
  * `creation_time` - (string)
  * `description` - (string)
  * `dynamic_discovery` - (bool)
  * `enabled` - (bool)
  * `id` - (string)
  * `modified_time` - (string)
  * `modified_by` - (string)
  * `name` - (string)
  * `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
  * `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
