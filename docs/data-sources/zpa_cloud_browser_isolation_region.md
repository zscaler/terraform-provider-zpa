---
page_title: "cloud_browser_isolation_region Data Source - terraform-provider-zpa"
subcategory: "Cloud Browser Isolation"
description: |-
  Get information about Cloud Browser Isolation Regions.
---

# zpa_cloud_browser_isolation_region (Data Source)

Use the **zpa_cloud_browser_isolation_region** data source to get information about Cloud Browser Isolation regions such as ID and Name. This data source information is required as part of the attribute `region_ids` when creating an Cloud Browser Isolation External Profile ``zpa_cloud_browser_isolation_external_profile``

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
# Retrieve CBI Region ID and Name
data "zpa_cloud_browser_isolation_region" "this" {
    name = "Singapore"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the CBI region to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (string) - ID information of the CBI region
