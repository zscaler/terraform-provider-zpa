---
subcategory: "Tag Controller"
layout: "zscaler"
page_title: "ZPA: tag_key"
description: |-
  Gets details of a ZPA tag key within a namespace.
---

# zpa_tag_key (Data Source)

Use the **zpa_tag_key** data source to get information about a tag key in Zscaler Private Access (ZPA).

~> NOTE: This an Early Access feature.

## Example Usage

```hcl
data "zpa_tag_namespace" "this" {
  name = "Example Namespace"
}

data "zpa_tag_key" "this" {
  name         = "Environment"
  namespace_id = data.zpa_tag_namespace.this.id
}

output "zpa_tag_key" {
  value = data.zpa_tag_key.this
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Optional) The ID of the tag key.
- `name` - (Optional) The name of the tag key.
- `namespace_id` - (Required) The ID of the tag namespace this key belongs to.
- `microtenant_id` - (Optional) The ID of the microtenant.

~> **NOTE:** `id` and `name` are mutually exclusive but at least one must be provided.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `description` - The description of the tag key.
- `enabled` - Whether the tag key is enabled.
- `customer_id` - The customer ID.
- `origin` - The origin of the tag key.
- `type` - The type of the tag key.
- `tag_values` - A list of tag values. Each element contains:
  - `id` - The ID of the tag value.
  - `name` - The name of the tag value.
- `microtenant_id` - The ID of the microtenant.
- `microtenant_name` - The name of the microtenant.
