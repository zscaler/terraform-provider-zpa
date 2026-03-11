---
subcategory: "Tag Controller"
layout: "zscaler"
page_title: "ZPA: tag_group"
description: |-
  Gets details of a ZPA tag group.
---

# zpa_tag_group (Data Source)

Use the **zpa_tag_group** data source to get information about a tag group in Zscaler Private Access (ZPA).

~> NOTE: This an Early Access feature.

## Example Usage

```hcl
data "zpa_tag_group" "this" {
  name = "Example Tag Group"
}

output "zpa_tag_group" {
  value = data.zpa_tag_group.this
}
```

```hcl
data "zpa_tag_group" "this" {
  id = "123456789"
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Optional) The ID of the tag group.
- `name` - (Optional) The name of the tag group.
- `microtenant_id` - (Optional) The ID of the microtenant.

~> **NOTE:** `id` and `name` are mutually exclusive but at least one must be provided.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `description` - The description of the tag group.
- `tags` - A list of tags in this tag group. Each element contains:
  - `origin` - The origin of the tag.
  - `namespace` - A list containing the tag namespace details:
    - `id` - The ID of the namespace.
    - `name` - The name of the namespace.
    - `enabled` - Whether the namespace is enabled.
  - `tag_key` - A list containing the tag key details:
    - `id` - The ID of the tag key.
    - `name` - The name of the tag key.
    - `enabled` - Whether the tag key is enabled.
  - `tag_value` - A list containing the tag value details:
    - `id` - The ID of the tag value.
    - `name` - The name of the tag value.
- `microtenant_id` - The ID of the microtenant.
- `microtenant_name` - The name of the microtenant.
