---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: cloud_browser_isolation_banner"
description: |-
  Get information about Cloud Browser Isolation Regions.
---

# Data Source: zpa_cloud_browser_isolation_banner

Use the **zpa_cloud_browser_isolation_banner** data source to get information about Cloud Browser Isolation banner. This data source information is required as part of the attribute `banner_id` when creating an Cloud Browser Isolation External Profile ``zpa_cloud_browser_isolation_external_profile``

## Example Usage

```hcl
# Retrieve CBI Region ID and Name
data "zpa_cloud_browser_isolation_banner" "this" {
    name = "Default"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the CBI banner to be exported.
* `id` - (Optional) The id of the CBI banner to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `primary_color` - (string) - The Banner Primary Color code in hexadecimal way to represent the color of the banner in RGB format
* `text_color` - (string) - The Banner Text Color code in hexadecimal way to represent the color of the text in RGB format
* `notification_title` - (string) The Banner Notification Title
* `notification_text` - (string) The Banner Notification Text
* `logo` - (string) - The Logo Image (.jpeg or .png; Maximum file size is 100KB.)
* `banner` - (bool) - Show Welcome Notification
* `Persist` - (bool) - Persist the default banner
* `is_default` - (bool) - Use the default banner
