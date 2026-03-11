---
subcategory: "Tag Controller"
layout: "zscaler"
page_title: "ZPA: tag_group"
description: |-
  Creates and manages ZPA tag groups.
---

# zpa_tag_group (Resource)

The **zpa_tag_group** resource creates and manages tag groups in Zscaler Private Access (ZPA). A tag group associates tag values together.

~> NOTE: This an Early Access feature.

## Example Usage

### Basic Usage

```hcl
resource "zpa_tag_group" "this" {
  name        = "Example Tag Group"
  description = "An example tag group"
}
```

### With Tag Values

```hcl
resource "zpa_tag_namespace" "this" {
  name        = "Example Namespace"
  description = "An example tag namespace"
  enabled     = true
}

resource "zpa_tag_key" "this" {
  name         = "Environment"
  description  = "Environment tag key"
  enabled      = true
  namespace_id = zpa_tag_namespace.this.id

  tag_values {
    name = "Production"
  }

  tag_values {
    name = "Staging"
  }
}

resource "zpa_tag_group" "this" {
  name        = "Example Tag Group"
  description = "An example tag group"

  tags = [
    zpa_tag_key.this.tag_values[0].id,
    zpa_tag_key.this.tag_values[1].id,
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required

- `name` - (Required) Name of the tag group.

### Optional

- `description` - (Optional) Description of the tag group.
- `tags` - (Optional) Set of tag value IDs to associate with this tag group.
- `microtenant_id` - (Optional) The ID of the microtenant the resource is to be associated with.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the tag group.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_tag_group** can be imported by using `<TAG GROUP ID>` or `<TAG GROUP NAME>` as the import ID.

For example:

```shell
terraform import zpa_tag_group.example <tag_group_id>
```

or

```shell
terraform import zpa_tag_group.example <tag_group_name>
```
