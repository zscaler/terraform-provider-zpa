---
subcategory: "Tag Controller"
layout: "zscaler"
page_title: "ZPA: tag_namespace"
description: |-
  Gets details of a ZPA tag namespace.
---

# zpa_tag_namespace (Data Source)

Use the **zpa_tag_namespace** data source to get information about a tag namespace in Zscaler Private Access (ZPA).

~> NOTE: This an Early Access feature.

## Example Usage

```hcl
data "zpa_tag_namespace" "this" {
  name = "Example Namespace"
}

output "zpa_tag_namespace" {
  value = data.zpa_tag_namespace.this
}
```

```hcl
data "zpa_tag_namespace" "this" {
  id = "123456789"
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Optional) The ID of the tag namespace.
- `name` - (Optional) The name of the tag namespace.
- `microtenant_id` - (Optional) The ID of the microtenant.

~> **NOTE:** `id` and `name` are mutually exclusive but at least one must be provided.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `description` - The description of the tag namespace.
- `enabled` - Whether the tag namespace is enabled.
- `origin` - The origin of the tag namespace.
- `type` - The type of the tag namespace.
- `microtenant_id` - The ID of the microtenant.
- `microtenant_name` - The name of the microtenant.
