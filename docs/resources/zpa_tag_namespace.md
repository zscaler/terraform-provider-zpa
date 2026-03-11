---
subcategory: "Tag Controller"
layout: "zscaler"
page_title: "ZPA: tag_namespace"
description: |-
  Creates and manages ZPA tag namespaces.
---

# zpa_tag_namespace (Resource)

The **zpa_tag_namespace** resource creates and manages tag namespaces in Zscaler Private Access (ZPA).

~> NOTE: This an Early Access feature.

## Example Usage

```hcl
resource "zpa_tag_namespace" "this" {
  name        = "Example Namespace"
  description = "An example tag namespace"
  enabled     = true
}
```

## Argument Reference

The following arguments are supported:

### Required

- `name` - (Required) Name of the tag namespace.

### Optional

- `description` - (Optional) Description of the tag namespace.
- `enabled` - (Optional) Whether this tag namespace is enabled.
- `microtenant_id` - (Optional) The ID of the microtenant the resource is to be associated with.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the tag namespace.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_tag_namespace** can be imported by using `<TAG NAMESPACE ID>` or `<TAG NAMESPACE NAME>` as the import ID.

For example:

```shell
terraform import zpa_tag_namespace.example <tag_namespace_id>
```

or

```shell
terraform import zpa_tag_namespace.example <tag_namespace_name>
```
