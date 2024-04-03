---
page_title: "zpa_emergency_access_user Resource - terraform-provider-zpa"
subcategory: "Emergency Access"
description: |-
  Official documentation https://help.zscaler.com/zpa/configuring-emergency-access
  API documentation https://help.zscaler.com/zpa/configuring-emergency-access-users-using-api
  Creates and manages emergency access users.
---

# zpa_emergency_access_user (Resource)

* [Official documentation](https://help.zscaler.com/zpa/configuring-emergency-access)
* [API documentation](https://help.zscaler.com/zpa/configuring-emergency-access-users-using-api)

The **zpa_emergency_access_user** Create emergency access users with permissions limited to privileged approvals in the specified IdP that is enabled for emergency access.

## Example Usage

```terraform
resource "zpa_emergency_access_user" "this" {
    email_id = "usertest@example.com"
    first_name = "User"
    last_name = "Test"
    user_id = "usertest"
}
```

## Schema

### Required

The following arguments are supported:

- `email_id` - (Required) The email address of the emergency access user, as provided by the admin
- `first_name` - (Required) The first name of the emergency access user.
- `last_name` - (Required) The last name of the emergency access user, as provided by the admin
- `user_id` - (Required) The unique identifier of the emergency access user.

## Import

The `zpa_emergency_access_user` do not support resource import.

