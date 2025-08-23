---
page_title: "zpa_user_portal_controller Data Source - terraform-provider-zpa"
subcategory: "User Portal Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-user-portals
  API documentation https://help.zscaler.com/zpa/about-user-portals
  Creates and manages ZPA User Portal details.
---

# zpa_user_portal_controller (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-user-portals)
* [API documentation](https://help.zscaler.com/zpa/about-user-portals)

The **zpa_user_portal_controller** data source to get information about a user portal in the Zscaler Private Access cloud.

## Example Usage - With Customer Own Certificate

```hcl

data "zpa_ba_certificate" "this" {
  name = "example.acme.com"
}

resource "zpa_user_portal_controller" "this" {
  name                      = "UserPortal01"
  description               = "UserPortal01"
  enabled                   = true
  user_notification         = "User_Portal_Terraform_01"
  user_notification_enabled = true
  certificate_id            = data.zpa_ba_certificate.this.id
  domain                    = "portal01"
}
```

## Example Usage - With Zscaler Managed Certificate

```hcl
data "zpa_user_portal_controller" "this" {
  name                      = "UserPortal01"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the user portal controller to be exported.

### Optional

* `id` - (Optional) The ID of the user portal controller to be exported.
* `microtenant_id` - (Optional) Microtenant ID for the user portal controller.

### Read-Only

* `certificate_id` - (String) Certificate ID for the User Portal Controller
* `certificate_name` - (String) Certificate name for the User Portal Controller
* `creation_time` - (String) Creation time of the User Portal Controller
* `description` - (String) Description of the User Portal Controller
* `domain` - (String) Domain for the User Portal Controller
* `enabled` - (Boolean) Whether this User Portal Controller is enabled or not
* `ext_domain` - (String) External domain for the User Portal Controller
* `ext_domain_name` - (String) External domain name for the User Portal Controller
* `ext_domain_translation` - (String) External domain translation for the User Portal Controller
* `ext_label` - (String) External label for the User Portal Controller
* `getc_name` - (String) GETC name for the User Portal Controller
* `image_data` - (String) Image data for the User Portal Controller
* `modified_by` - (String) Modified by information for the User Portal Controller
* `modified_time` - (String) Modified time of the User Portal Controller
* `microtenant_name` - (String) Microtenant name for the User Portal Controller
* `user_notification` - (String) User notification message for the User Portal Controller
* `user_notification_enabled` - (Boolean) Whether user notifications are enabled for the User Portal Controller
* `managed_by_zs` - (Boolean) Whether the User Portal Controller is managed by Zscaler
