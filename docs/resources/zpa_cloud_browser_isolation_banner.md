---
page_title: "zpa_cloud_browser_isolation_banner Resource - terraform-provider-zpa"
subcategory: "Cloud Browser Isolation"
description: |-
  Official documentation https://help.zscaler.com/isolation/adding-banner-theme-isolation-end-user-notification-zpa
  Creates and manages Cloud Browser Isolation Banner.
---

# zpa_cloud_browser_isolation_banner (Resource)

* [Official documentation](https://help.zscaler.com/isolation/adding-banner-theme-isolation-end-user-notification-zpa)

The **zpa_cloud_browser_isolation_banner** resource creates a Cloud Browser Isolation banner. This resource is required as part of the attribute `banner_id` when creating an Cloud Browser Isolation External Profile ``zpa_cloud_browser_isolation_external_profile``

## Example Usage

```terraform
resource "zpa_cloud_browser_isolation_banner" "this" {
  name = "CBI_Banner_Example"
  primary_color = "#0076BE"
  text_color = "#FFFFFF"
  notification_title = "Heads up, you’ve been redirected to Browser Isolation!"
  notification_text = "The website you were trying to access is now rendered in a fully isolated environment to protect you from malicious content."
  banner =  true
  persist = true
  logo = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAYQAAABQCAMAAAAuu
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (Required) The name of the CBI banner to be exported.
- `primary_color` - (Required) - The Banner Primary Color code in hexadecimal way to represent the color of the banner in RGB format
- `text_color` - (Required) - The Banner Text Color code in hexadecimal way to represent the color of the text in RGB format
- `notification_title` - (Required) The Banner Notification Title
- `notification_text` - (Required) The Banner Notification Text
- `logo` - (Required) - The Logo Image (.jpeg or .png; Maximum file size is 100KB.)

### Optional

In addition to all arguments above, the following attributes are exported:

- `banner` - (Optional) - Show Welcome Notification
- `Persist` - (Optional) - Persist the default banner
- `is_default` - (Optional) - Use the default banner

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Application Segment can be imported by using `<BANNER ID>` or `<BANNER NAME>` as the import ID.

```shell
terraform import zpa_cloud_browser_isolation_banner.example <banner_id>
```

or

```shell
terraform import zpa_cloud_browser_isolation_banner.example <banner_name>
```
