---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: access_policy_platforms"
description: |-
  Get information about all platforms for the specified customer.
---

# Data Source: zpa_access_policy_platforms

Use the **zpa_access_policy_platforms** data source to get information about all platforms for the specified customer in the Zscaler Private Access cloud. This data source can be optionally used when defining the following policy types:
    - ``zpa_policy_access_rule``
    - ``zpa_policy_timeout_rule``
    - ``zpa_policy_forwarding_rule``
    - ``zpa_policy_isolation_rule``
    - ``zpa_policy_inspection_rule``

The ``object_type`` attribute must be defined as "PLATFORM" in the policy operand condition. To learn more see the To learn more see the [Getting Platform Types for a Customer](https://help.zscaler.com/zpa/configuring-access-policies-using-api#getPlatformTypes)

-> **NOTE** By Default the ZPA provider will return all platform types

## Example Usage

```hcl
data "zpa_access_policy_platforms" "this" {
}
```

## Argument Reference

The following values are returned:

* `"android" = "Android"`
* `"id" = "platforms"`
* `"ios" = "iOS"`
* `"linux" = "Linux"`
* `"mac" = "Mac"`
* `"windows" = "Windows"`

To learn more see the [Getting Platform Types for a Customer](https://help.zscaler.com/zpa/configuring-access-policies-using-api#getPlatformTypes)
