---
subcategory: "Emergency Access"
layout: "zscaler"
page_title: "ZPA: emergency_access_user"
description: |-
  Creates and manages emergency access users.
---

# Resource: zpa_emergency_access_user

The **zpa_emergency_access_user** Create emergency access users with permissions limited to privileged approvals in the specified IdP that is enabled for emergency access.

## Example Usage

```hcl
resource "zpa_emergency_access_user" "this" {
    email_id = "usertest@example.com"
    first_name = "User"
    last_name = "Test"
    user_id = "usertest"
}
```

## Argument Reference

The following arguments are supported:

* `email_id` - (Required) The email address of the emergency access user, as provided by the admin
* `first_name` - (Required) The first name of the emergency access user.
* `last_name` - (Required) The last name of the emergency access user, as provided by the admin
* `user_id` - (Required) The unique identifier of the emergency access user.

## Import

The `zpa_emergency_access_user` do not support resource import.

