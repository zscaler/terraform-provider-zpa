---
subcategory: "Segment Group"
layout: "zpa"
page_title: "ZPA: segment_group"
description: |-
  Creates a ZPA Segment Group resource
---
# zpa_segment_group (Resource)

The **zpa_segment_group** resource creates a segment group in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule or application segment resource.

## Example Usage

```hcl
# ZPA Segment Group resource
resource "zpa_segment_group" "example" {
  name = "Example"
  description = "Example"
  enabled = true
  policy_migrated = true
  tcp_keep_alive_enabled = true
}
```

### Required

* `name` - (Required) Name of the app group.

## Attributes Reference

* `description` (String) Description of the app group.
* `enabled` (Optional) Whether this app group is enabled or not.
* `config_space` (String)
* `policy_migrated` (Boolean)
* `tcp_keep_alive_enabled` (Number)

`applications`

* `id` (Number) The ID of Application Segment resources.

## Import

Application Segment Group can be imported; use `<SEGMENT GROUP ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_forwarding_rule.example 216196257331290863
```
