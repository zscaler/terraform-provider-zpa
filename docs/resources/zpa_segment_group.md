---
page_title: "zpa_segment_group Resource - terraform-provider-zpa"
subcategory: "Segment Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-segment-groups
  API documentation https://help.zscaler.com/zpa/configuring-segment-groups-using-api
  Creates and manages ZPA Segment Group resource
---

# zpa_segment_group (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-segment-groups)
* [API documentation](https://help.zscaler.com/zpa/configuring-segment-groups-using-api)

The **zpa_segment_group** resource creates a segment group in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule or application segment resource.

~> **NOTE: The `applications` attribute is deprecated and will be removed in a future major release.** An application segment always belongs to exactly one segment group, and that association can only be changed through the application segment. Manage membership with the `segment_group_id` attribute of [`zpa_application_segment`](zpa_application_segment.md) (see [Associating Application Segments](#example-usage---associating-application-segments) below), not from the segment group:

* Updating a segment group (for example its name or description) never changes its applications.
* Adding or removing applications in the `applications` block of an existing segment group is rejected at plan time.
* Configurations that list the segment group's current applications keep working, but show a deprecation warning. Remove the `applications` block to clear it.

## Zenith Community - ZPA Segment Group

[![ZPA Terraform provider Video Series Ep6 - Segment Group](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_segment_groups.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEfCAI/video-zpa-terraform-provider-video-series-ep6-zpa-segment-group)

## Example Usage

```terraform
# ZPA Segment Group resource
resource "zpa_segment_group" "test_segment_group" {
  name                   = "test1-segment-group"
  description            = "test1-segment-group"
  enabled                = true
}
```

## Example Usage - Associating Application Segments

Application segments are added to a segment group by referencing the segment group in the application segment, using `segment_group_id`. Do not list them in the segment group's `applications` block. To move an application segment to another segment group, change its `segment_group_id`; the application segment is updated in place.

```terraform
resource "zpa_segment_group" "this" {
  name        = "Example Segment Group"
  description = "Example Segment Group"
  enabled     = true
}

data "zpa_server_group" "this" {
  name = "Example Server Group"
}

resource "zpa_application_segment" "this" {
  name             = "Example Application Segment"
  description      = "Example Application Segment"
  enabled          = true
  domain_names     = ["app.example.com"]
  segment_group_id = zpa_segment_group.this.id
  tcp_port_ranges  = ["443", "443"]
  server_groups {
    id = [data.zpa_server_group.this.id]
  }
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (String) Name of the segment group.

### Optional

In addition to all arguments above, the following attributes are exported:

* `description` (String) Description of the segment group.
* `enabled` (Optional) Whether this segment group is enabled or not.
* `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.
* `applications` - (Deprecated) The application segments in the segment group. Application segments are associated with a segment group through the `segment_group_id` attribute of `zpa_application_segment`; an application segment always belongs to exactly one segment group, so membership cannot be changed from the segment group. Adding or removing applications in this block on an existing segment group is rejected at plan time. This attribute will be removed in a future major release; remove it from the configuration.

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
