---
page_title: "zpa_cloud_browser_isolation_zpa_profile Data Source - terraform-provider-zpa"
subcategory: "Cloud Browser Isolation"
description: |-
  Get information about an Isolation Profile in Zscaler Private Access cloud.
---

# zpa_cloud_browser_isolation_zpa_profile (Data Source)

Use the **zpa_cloud_browser_isolation_zpa_profile** data source to get information about an isolation profile in the Zscaler Private Access cloud. This data source is required when configuring an isolation policy rule resource

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
data "zpa_cloud_browser_isolation_zpa_profile" "this" {
    name = "ZPA_Profile"
}
```

## Schema

### Required

* `name` - (String) This field defines the name of the isolation profile.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (String) This field defines the id of the isolation profile.
* `description` - (string)
* `enabled` - (string)
* `cbi_tenant_id` - (string)
* `cbi_profile_id` - (string)
* `cbi_url` - (string)
* `creation_time` - (string)
* `modified_by` - (string)
* `modified_time` - (string)
