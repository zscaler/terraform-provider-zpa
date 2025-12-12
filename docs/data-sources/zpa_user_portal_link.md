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

The **zpa_user_portal_link** data source to get information about a user portal link in the Zscaler Private Access cloud.

## Example Usage

```hcl

data "zpa_user_portal_link" "this" {
  name        = "server1.example.com"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the user portal link to be exported.

### Optional

* `id` - (Optional) The ID of the user portal link to be exported.
* `microtenant_id` - (Optional) Microtenant ID for the user portal link.

### Read-Only

* `application_id` - (String) Application ID for the User Portal Link
* `creation_time` - (String) Creation time of the User Portal Link
* `description` - (String) Description of the User Portal Link
* `enabled` - (Boolean) Whether this User Portal Link is enabled or not
* `icon_text` - (String) Icon text for the User Portal Link
* `link` - (String) Link URL for the User Portal Link
* `link_path` - (String) Link path for the User Portal Link
* `modified_by` - (String) Modified by information for the User Portal Link
* `modified_time` - (String) Modified time of the User Portal Link
* `protocol` - (String) Protocol for the User Portal Link
* `user_portal_id` - (String) User Portal ID for the User Portal Link
* `user_portals` - (List) List of User Portals associated with this link
  * `id` - (String) ID of the user portal
  * `name` - (String) Name of the user portal
  * `enabled` - (Boolean) Whether the user portal is enabled