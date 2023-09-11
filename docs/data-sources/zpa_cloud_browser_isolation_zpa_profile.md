---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: cloud_browser_isolation_zpa_profile"
description: |-
  Get information about an Isolation Profile in Zscaler Private Access cloud.
---

# Data Source: zpa_cloud_browser_isolation_zpa_profile

Use the **zpa_cloud_browser_isolation_zpa_profile** data source to get information about an isolation profile in the Zscaler Private Access cloud. This data source is required when configuring an isolation policy rule resource

## Example Usage

```hcl
data "zpa_cloud_browser_isolation_zpa_profile" "this" {
    name = "ZPA_Profile"
}
```

## Argument Reference

* `name` - (Required) This field defines the name of the isolation profile.
* `id` - (Optional) This field defines the id of the isolation profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string)
* `enabled` - (string)
* `cbi_tenant_id` - (string)
* `cbi_profile_id` - (string)
* `cbi_url` - (string)
* `creation_time` - (string)
* `modified_by` - (string)
* `modified_time` - (string)
