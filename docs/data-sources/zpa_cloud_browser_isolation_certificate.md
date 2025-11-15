---
page_title: "zpa_cloud_browser_isolation_certificate Data Source - terraform-provider-zpa"
subcategory: "Cloud Browser Isolation"
description: |-
  Official documentation https://help.zscaler.com/isolation/about-custom-root-certificates-cloud-browser-isolation
  Get information about CBI Certificate for the customer based on the specified ID.
---

# zpa_cloud_browser_isolation_certificate (Data Source)

* [Official documentation](https://help.zscaler.com/isolation/adding-banner-theme-isolation-end-user-notification-zpa)

Use the **zpa_cloud_browser_isolation_certificate** data source to get information about Cloud Browser Isolation Certificate. This data source information is required as part of the attribute `certificate_ids` when creating an Cloud Browser Isolation External Profile ``zpa_cloud_browser_isolation_external_profile``

## Example Usage

```terraform
# Retrieve CBI Certificate ID and Name
data "zpa_cloud_browser_isolation_certificate" "this" {
  name = "Zscaler Root Certificate"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the CBI certificate to be exported.
* `id` - (Optional) The id of the CBI certificate to be exported.

### Read-Only

* N/A

~> **Warning**: Notice that certificate and public_keys are omitted from the output.
