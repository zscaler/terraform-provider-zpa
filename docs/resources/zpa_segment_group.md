---
subcategory: "Segment Group"
layout: "zscaler"
page_title: "ZPA): segment_group"
description: |-
  Creates and manages ZPA Segment Group resource
---

# Resource: zpa_segment_group

The **zpa_segment_group** resource creates a segment group in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule or application segment resource.

[![ZPA Terraform provider Video Series Ep6 - Segment Group](../../images/zpa_segment_groups.svg)](https://community.zscaler.com/t/video-zpa-terraform-provider-video-series-ep-6-zpa-segment-group/18808)

## Example Usage

```hcl
# ZPA Segment Group resource
resource "zpa_segment_group" "test_segment_group" {
  name                   = "test1-segment-group"
  description            = "test1-segment-group"
  enabled                = true
  tcp_keep_alive_enabled = "1"
}
```

## Attributes Reference

### Required

* `name` - (Required) Name of the segment group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` (Optional) Description of the segment group.
* `enabled` (Optional) Whether this segment group is enabled or not.
* `config_space` (Optional)
* `tcp_keep_alive_enabled` (Optional)

## Import

**segment_group** can be imported by using `<SEGMENT GROUP ID>` or `<SEGMENT GROUP NAME>` as the import ID.

For example:

```shell
terraform import zpa_segment_group.example <segment_group_id>
```

or

```shell
terraform import zpa_segment_group.example <segment_group_name>
```
