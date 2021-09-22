---
# generated by https://github.com/hashicorp/terraform-plugin-docs
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

- **name** (String) Name of the app group.
- **enabled** (Boolean) Whether this app group is enabled or not.

### Optional

- **applications** (Block List) (see [below for nested schema](#nestedblock--applications))
- **config_space** (String)
- **description** (String) Description of the app group.
- **policy_migrated** (Boolean)

### Read-Only

- **id** (String) The ID of this resource.
- **tcp_keep_alive_enabled** (Number)

<a id="nestedblock--applications"></a>
### Nested Schema for `applications`

Optional:

- **id** (Number) The ID of this resource.

