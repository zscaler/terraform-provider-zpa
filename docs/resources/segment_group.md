---
subcategory: "Segment Group"
layout: "zpa"
page_title: "ZPA: segment_group"
description: |-
  Creates a ZPA Segment Group resource.

---
# zpa_segment_group (Resource)

The **zpa_segment_group** resource creates a Segment Group in the Zscaler Private Access portal. This resource is required, when creating an application segment and can also be associated with an Access Policy.

## Example Usage

```hcl
# ZPA Server Group Data Source
  resource "zpa_segment_group" "example" {
   name = "Example"
   description = "Example"
   enabled = true
   policy_migrated = true
 }
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the app group
* `description` - (Optional) Description of the app group.
* `enabled` - (Required) Whether this app group is enabled or not. The accepted values are `true`, `false`
* `policy_migrated` - (Optional)
* `config_space` - (Optional)
* `tcp_keep_alive_enabled` - (Optional)

`applications` - (Optional) The App ID.

* `id` - (Optional)
