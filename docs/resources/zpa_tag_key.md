---
subcategory: "Tag Controller"
layout: "zscaler"
page_title: "ZPA: tag_key"
description: |-
  Creates and manages ZPA tag keys within a tag namespace.
---

# zpa_tag_key (Resource)

The **zpa_tag_key** resource creates and manages tag keys within a tag namespace in Zscaler Private Access (ZPA).

~> NOTE: This an Early Access feature.

## Example Usage

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
```

## Argument Reference

The following arguments are supported:

### Required

- `name` - (Required) Name of the tag key.
- `namespace_id` - (Required, ForceNew) The ID of the tag namespace this key belongs to.

### Optional

- `description` - (Optional) Description of the tag key.
- `enabled` - (Optional) Whether this tag key is enabled.
- `tag_values` - (Optional) List of tag values associated with this tag key. Each `tag_values` block supports:
  - `name` - (Required) Name of the tag value.
- `microtenant_id` - (Optional) The ID of the microtenant the resource is to be associated with.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the tag key.
- `tag_values` - The list of tag values. Each element contains:
  - `id` - The ID of the tag value.
  - `name` - The name of the tag value.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_tag_key** can be imported by using the composite format `<NAMESPACE ID>/<TAG KEY ID>` or `<NAMESPACE ID>/<TAG KEY NAME>` as the import ID.

For example:

```shell
terraform import zpa_tag_key.example <namespace_id>/<tag_key_id>
```

or

```shell
terraform import zpa_tag_key.example <namespace_id>/<tag_key_name>
```
