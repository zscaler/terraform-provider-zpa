---
page_title: "zpa_user_portal_link Resource - terraform-provider-zpa"
subcategory: "User Portal Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-user-portals
  API documentation https://help.zscaler.com/zpa/about-user-portals
  Creates and manages ZPA User Portal details.
---

# zpa_user_portal_link (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-user-portals)
* [API documentation](https://help.zscaler.com/zpa/about-user-portals)

The **zpa_user_portal_link** resource creates a user portal link in the Zscaler Private Access cloud.

## Example Usage

```hcl

resource "zpa_user_portal_link" "this" {
  name        = "server1.example.com"
  description = "server1.example.com"
  enabled     = true
  link        = "server1.example.com"
  icon_text   = ""
  protocol    = "https://"
  user_portals {
    id = [zpa_user_portal_controller.this.id]
  }
}
```

## Schema

### Required

- `name` (String) - Name of the User Portal Link

### Optional

- `id` (String) - The ID of the User Portal Link
- `description` (String) - Description of the User Portal Link
- `enabled` (Boolean) - Whether this User Portal Link is enabled or not
- `icon_text` (String) - Icon text for the User Portal Link
- `link` (String) - Link URL for the User Portal Link
- `link_path` (String) - Link path for the User Portal Link
- `protocol` (String) - Protocol for the User Portal Link
- `microtenant_id` (String) - Microtenant ID for the User Portal Link
- `user_portals` (List) - List of User Portals
  * `id` (Set of String) - List of User Portal IDs

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

**user_portal_controller** can be imported by using `<USER PORTAL ID>` or `<USER PORTAL NAME>` as the import ID.

For example:

```shell
terraform import zpa_user_portal_link.example <portal_id>
```

or

```shell
terraform import zpa_user_portal_link.example <portal_name>
```