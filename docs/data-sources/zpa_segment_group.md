---
subcategory: "Segment Group"
layout: "zscaler"
page_title: "ZPA: segment_group"
description: |-
  Get information about Segment Groups in Zscaler Private Access cloud.
---

# Data Source: zpa_segment_group

Use the **zpa_segment_group** data source to get information about a machine group created in the Zscaler Private Access cloud. This data source can then be referenced in an application segment or Access Policy rule.

-> **NOTE:** Segment Groups association is only supported in an Access Policy rule.

## Zenith Community - ZPA Segment Group

[![ZPA Terraform provider Video Series Ep6 - Segment Group](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_segment_groups.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEfCAI/video-zpa-terraform-provider-video-series-ep6-zpa-segment-group)

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

* `config_space` - (string)
* `creation_time` - (string)
* `description` - (string)
* `enabled` - (bool)
* `microtenant_id` (string)
* `microtenant_name` (string)
