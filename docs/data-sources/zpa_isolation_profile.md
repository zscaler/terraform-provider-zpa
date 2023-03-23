---
subcategory: "Isolation"
layout: "zscaler"
page_title: "ZPA: isolation_profile"
description: |-
  Get information about an Isolation Profile in Zscaler Private Access cloud.
---

# Data Source: zpa_isolation_profile

Use the **zpa_isolation_profile** data source to get information about an isolation profile in the Zscaler Private Access cloud. This data source is required when configuring an isolation policy rule resource

## Example Usage

```hcl
data "zpa_isolation_profile" "isolation_profile" {
    name = "zpa_isolation_profile"
}
```

## Argument Reference

* `name` - (Required) This field defines the name of the isolation profile.
* `id` - (Optional) This field defines the id of the isolation profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string)
* `enabled` - (string)
* `isolation_profile_id` - (string)
* `isolation_tenant_id` - (string)
* `isolation_url` - (string)
* `creation_time` - (string)
* `modified_by` - (string)
* `modified_time` - (string)
