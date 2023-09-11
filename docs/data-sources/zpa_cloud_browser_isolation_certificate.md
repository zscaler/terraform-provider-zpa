---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: cloud_browser_isolation_banner"
description: |-
  Get information about CBI Certificate for the customer based on the specified ID.
---

# Data Source: zpa_cloud_browser_isolation_certificate

Use the **zpa_cloud_browser_isolation_certificate** data source to get information about Cloud Browser Isolation Certificate. This data source information is required as part of the attribute `certificate_ids` when creating an Cloud Browser Isolation External Profile ``zpa_cloud_browser_isolation_external_profile``

## Example Usage

```hcl
# Retrieve CBI Certificate ID and Name
data "zpa_cloud_browser_isolation_certificate" "this" {
  name = "Zscaler Root Certificate"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the CBI certificate to be exported.
* `id` - (Optional) The id of the CBI certificate to be exported.

## Attribute Reference

* N/A

:warning: Notice that certificate and public_keys are omitted from the output.
