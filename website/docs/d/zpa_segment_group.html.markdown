---
layout: "zscaler"
page_title: "Zscaler Private Access (ZPA): segment_group"
sidebar_current: "docs-datasource-zpa-segment-group"
description: |-
  Get information about Segment Groups in Zscaler Private Access cloud.
---

# zpa_segment_group

Use the **zpa_segment_group** data source to get information about a machine group created in the Zscaler Private Access cloud. This data source can then be referenced in an application segment or Access Policy rule.

-> **NOTE:** Segment Groups association is only supported in an Access Policy rule.

## Example Usage

```hcl
# ZPA Server Group Data Source
data "zpa_segment_group" "example" {
 name = "segment_group_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the segment group to be exported.
* `id` - (Optional) The ID of the segment group to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `config_space` - (Computed)
* `creation_time` - (Computed)
* `description` - (Computed)
* `enabled` - (Computed)
* `policy_migrated` - (Computed)
* `tcp_keep_alive_enabled` - (Computed)

* `applications` - (Computed)
  * `bypass_type` - (Computed)
  * `config_space` - (Computed)
  * `default_idle_timeout` - (Computed)
  * `default_max_age` - (Computed)
  * `description` - (Computed)
  * `domain_name` - (Computed)
  * `domain_names`  - (Computed)
  * `double_encrypt` - (Computed)
  * `enabled` - (Computed)
  * `health_check_type` - (Computed)
  * `id` - (Computed)
  * `ip_anchored` - (Computed)
  * `log_features` - (Computed)
  * `name` - (Computed)
  * `passive_health_enabled` - (Computed)
  * `tcp_port_ranges` - (Computed)
  * `tcp_ports_in`  - (Computed)
  * `udp_port_ranges` - (Computed)

* `server_groups` - (Computed)
  * `config_space` - (Computed)
  * `creation_time` - (Computed)
  * `description` - (Computed)
  * `dynamic_discovery` - (Computed)
  * `enabled` - (Computed)
  * `id` - (Computed)
  * `modified_time` - (Computed)
  * `modifiedby` - (Computed)
  * `name` - (Computed)
