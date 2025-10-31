---
page_title: "zpa_user_portal_aup Resource - terraform-provider-zpa"
subcategory: "User Portal Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-user-portals
  API documentation https://help.zscaler.com/zpa/about-user-portals
  Creates and manages ZPA User Portal details.
---

# zpa_user_portal_aup (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-user-portals)
* [API documentation](https://help.zscaler.com/zpa/about-user-portals)

The **zpa_user_portal_aup** resource creates a user portal link in the Zscaler Private Access cloud.

## Example Usage

```hcl
resource "zpa_user_portal_aup" "this" {
  name        = "Org_AUP01"
  description = "Org_AUP01"
  enabled     = true
  aup         = "Org_AUP01"
  email       = "company@acme.com"
  phone_num   = "+1 123-1458"
}
```

## Schema

### Required

- `name` (String) - Name of the User Portal AUP

### Optional

* `description` - (String) Description of the User Portal AUP
* `enabled` - (Boolean) Whether this User Portal AUP is enabled or not
* `aup` - (String)
* `email` - (String)
* `phone_num` - (String)

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zpa_user_portal_aup** can be imported by using `<USER PORTAL ID>` or `<USER PORTAL NAME>` as the import ID.

For example:

```shell
terraform import zpa_user_portal_aup.example <portal_id>
```

or

```shell
terraform import zpa_user_portal_aup.example <portal_name>
```