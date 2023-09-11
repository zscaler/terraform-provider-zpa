---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: cloud_browser_isolation_region"
description: |-
  Get information about Cloud Browser Isolation Regions.
---

# Data Source: zpa_cloud_browser_isolation_region

Use the **zpa_cloud_browser_isolation_region** data source to get information about Cloud Browser Isolation regions such as ID and Name. This data source information is required as part of the attribute `region_ids` when creating an Cloud Browser Isolation External Profile ``zpa_cloud_browser_isolation_external_profile``

## Example Usage

```hcl
# Retrieve CBI Region ID and Name
data "zpa_cloud_browser_isolation_region" "this" {
    name = "Singapore"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the CBI region to be exported.
* `id` - (Optional) The id of the CBI region to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (string) - ID information of the CBI region
* `name` - (string) - Name of the CBI region
