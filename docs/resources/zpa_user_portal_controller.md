---
page_title: "zpa_user_portal_controller Resource - terraform-provider-zpa"
subcategory: "User Portal Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-user-portals
  API documentation https://help.zscaler.com/zpa/about-user-portals
  Creates and manages ZPA User Portal details.
---

# zpa_user_portal_controller (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-user-portals)
* [API documentation](https://help.zscaler.com/zpa/about-user-portals)

The **zpa_user_portal_controller** resource creates a user portal in the Zscaler Private Access cloud.

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
resource "zpa_user_portal_controller" "this" {
  name                      = "UserPortal01"
  description               = "UserPortal01"
  enabled                   = true
  user_notification         = "UserPortal01"
  user_notification_enabled = true
  certificate_id            = ""
  ext_domain_translation    = "acme.io"
  ext_label                 = "portal01"
  ext_domain_name           = "acme-io.b.zscalerportal.net"
  ext_domain                = "acme.io"
  domain                    = "portal01-acme-io.b.zscalerportal.net"
}
```

## Schema

### Required

- `name` (String) - Name of the User Portal Controller

### Optional

- `id` (String) - The ID of the User Portal Controller
- `certificate_id` (String) - Certificate ID for the User Portal Controller
- `description` (String) - Description of the User Portal Controller
- `domain` (String) - Domain for the User Portal Controller
- `enabled` (Boolean) - Whether this User Portal Controller is enabled or not
- `ext_domain` (String) - External domain for the User Portal Controller
- `ext_domain_name` (String) - External domain name for the User Portal Controller
- `ext_domain_translation` (String) - External domain translation for the User Portal Controller
- `ext_label` (String) - External label for the User Portal Controller
- `microtenant_id` (String) - Microtenant ID for the User Portal Controller
- `user_notification` (String) - User notification message for the User Portal Controller
- `user_notification_enabled` (Boolean) - Whether user notifications are enabled for the User Portal Controller

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**user_portal_controller** can be imported by using `<USER PORTAL ID>` or `<USER PORTAL NAME>` as the import ID.

For example:

```shell
terraform import zpa_user_portal_controller.example <portal_id>
```

or

```shell
terraform import zpa_user_portal_controller.example <portal_name>
```
