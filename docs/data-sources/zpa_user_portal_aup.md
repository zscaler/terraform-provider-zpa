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

The **zpa_user_portal_aup** data source to get information about a user portal AUP (Acceptance User Policy) in the Zscaler Private Access cloud.

## Example Usage - By Name

```hcl
data "zpa_user_portal_aup" "this" {
  name        = "server1.example.com"
}
```

## Example Usage - By ID

```hcl
data "zpa_user_portal_aup" "this" {
  name        = "server1.example.com"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the user portal aup to be exported.

### Optional

* `id` - (Optional) The ID of the user portal aup to be exported.
* `microtenant_id` - (Optional) Microtenant ID for the user portal aup.

### Read-Only

* `description` - (String) Description of the User Portal AUP
* `enabled` - (Boolean) Whether this User Portal AUP is enabled or not
* `aup` - (String)
* `email` - (String)
* `phone_num` - (String)
* `microtenant_name` - (Striing)
* `microtenant_id` - (String)
* `creation_time` - (String) Creation time of the User Portal Link
* `modified_by` - (String) Modified by information for the User Portal Link
* `modified_time` - (String) Modified time of the User Portal Link