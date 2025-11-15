---
page_title: "zpa_isolation_profile Data Source - terraform-provider-zpa"
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: isolation_profile"
description: |-
  Official documentation https://help.zscaler.com/isolation/creating-isolation-profiles-zpa
  API documentation https://help.zscaler.com/zpa/obtaining-isolation-profile-details-using-api
  Get information about an Isolation Profile in Zscaler Private Access cloud.
---

# zpa_isolation_profile (Data Source)

* [Official documentation](https://help.zscaler.com/isolation/creating-isolation-profiles-zpa)
* [API documentation](https://help.zscaler.com/zpa/obtaining-isolation-profile-details-using-api)

Use the **zpa_isolation_profile** data source to get information about an isolation profile in the Zscaler Private Access cloud. This data source is required when configuring an isolation policy rule resource

## Example Usage

```terraform
data "zpa_isolation_profile" "isolation_profile" {
    name = "zpa_isolation_profile"
}
```

## Schema

### Required

* `name` - (Required) This field defines the name of the isolation profile.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) This field defines the id of the isolation profile.
* `description` - (string)
* `enabled` - (string)
* `isolation_profile_id` - (string)
* `isolation_tenant_id` - (string)
* `isolation_url` - (string)
* `creation_time` - (string)
* `modified_by` - (string)
* `modified_time` - (string)
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
