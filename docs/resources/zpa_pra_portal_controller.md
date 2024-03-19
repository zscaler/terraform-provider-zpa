---
subcategory: "Privileged Remote Access"
layout: "zscaler"
page_title: "ZPA): pra_portal_controller"
description: |-
  Creates and manages ZPA privileged remote access portal
---

# Resource: zpa_pra_portal_controller

The **zpa_pra_portal_controller** resource creates a privileged remote access portal in the Zscaler Private Access cloud. This resource can then be referenced in an privileged remote access console resource.

## Example Usage

```hcl
# Retrieves Browser Access Certificate
data "zpa_ba_certificate" "this" {
 name = "portal.acme.com"
}

resource "zpa_pra_portal_controller" "this" {
  name = "portal.acme.com"
  description = "portal.acme.com"
  enabled = true
  domain = "portal.acme.com"
  certificate_id = data.zpa_ba_certificate.this.id
  user_notification = "Created with Terraform"
  user_notification_enabled = true
}
```

## Attributes Reference

### Required

* `name` - (Required) The name of the privileged portal.
* `domain` - (Required) The domain of the privileged portal.
* `certificate_id` - (Required) The unique identifier of the certificate.
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` (Optional) The description of the privileged portal.
* `enabled` (Optional) Whether or not the privileged portal is enabled.
* `user_notification` (Optional) The notification message displayed in the banner of the privileged portallink, if enabled.
* `user_notification_enabled` (Optional) Indicates if the Notification Banner is enabled (true) or disabled (false).
* `microtenant_id` (Optional) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**pra_portal_controller** can be imported by using `<PORTAL ID>` or `<PORTAL NAME>` as the import ID.

For example:

```shell
terraform import zpa_pra_portal_controller.this <portal_id>
```

or

```shell
terraform import zpa_pra_portal_controller.this <portal_name>
```
