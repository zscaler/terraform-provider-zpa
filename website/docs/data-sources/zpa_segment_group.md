---
subcategory: "Segment Group"
layout: "zpa"
page_title: "ZPA: segment_group"
description: |-
  Gets a ZPA Segment Group details.

---
# zpa_segment_group

The **zpa_segment_group** data source provides details about a specific Segment Group created in the Zscaler Private Access.
This data source is required when creating:

1. Access policy rule

## Example Usage

```hcl
# ZPA Server Group Data Source
data "zpa_segment_group" "example" {
 name = "segment_group_name"
}

output "zpa_segment_group" {
  value = data.zpa_segment_group.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the asegment group to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Read-Only

* `config_space` - (String)
* `creation_time` - (Number)
* `description` - (String)
* `enabled` - (Boolean)
* `policy_migrated` - (Boolean)
* `tcp_keep_alive_enabled` - (Number)

`applications` - (List of Object)

* `bypass_type` - (String)
* `config_space` - (String)
* `default_idle_timeout` - (Number)
* `default_max_age` - (Number)
* `description` - (String)
* `domain_name` - (String)
* `domain_names`  - List of String)
* `double_encrypt` - (Boolean)
* `enabled` - (Boolean)
* `health_check_type` - (String)
* `id` - (Number)
* `ip_anchored` - (Boolean)
* `log_features` - (List of String)
* `name` - (String)
* `passive_health_enabled` - (Boolean)
* `tcp_port_ranges` - (List of String)
* `tcp_ports_in`  - (List of String)
* `udp_port_ranges` - (List of String)

`server_groups` - (List of Object)

* `config_space` - (String)
* `creation_time` - (Number)
* `description` - (String)
* `dynamic_discovery` - (Boolean)
* `enabled` - (Boolean)
* `id` (Number)
* `modified_time` - (Number)
* `modifiedby` - (Number)
* `name` - (String)
