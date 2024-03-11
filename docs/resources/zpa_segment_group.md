---
subcategory: "Segment Group"
layout: "zscaler"
page_title: "ZPA): segment_group"
description: |-
  Creates and manages ZPA Segment Group resource
---

# Resource: zpa_segment_group

The **zpa_segment_group** resource creates a segment group in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule or application segment resource.

## Zenith Community - ZPA Segment Group

[![ZPA Terraform provider Video Series Ep6 - Segment Group](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_segment_groups.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEfCAI/video-zpa-terraform-provider-video-series-ep6-zpa-segment-group)

## Example Usage

```hcl
# ZPA Segment Group resource
resource "zpa_segment_group" "test_segment_group" {
  name                   = "test1-segment-group"
  description            = "test1-segment-group"
  enabled                = true
}
```

## Attributes Reference

### Required

* `name` - (Required) Name of the segment group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` (Optional) Description of the segment group.
* `enabled` (Optional) Whether this segment group is enabled or not.
* `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**segment_group** can be imported by using `<SEGMENT GROUP ID>` or `<SEGMENT GROUP NAME>` as the import ID.

For example:

```shell
terraform import zpa_segment_group.example <segment_group_id>
```

or

```shell
terraform import zpa_segment_group.example <segment_group_name>
```
