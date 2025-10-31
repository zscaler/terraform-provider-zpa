---
page_title: "zpa_branch_connector_group Data Source - terraform-provider-zpa"
subcategory: "Location Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-branch-connector-groups
  API documentation https://help.zscaler.com/zpa/about-branch-connector-groups
  Get information about branch connector Group
---

# zpa_branch_connector_group (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-branch-connector-groups)
* [API documentation](https://help.zscaler.com/zpa/about-branch-connector-groups)

Use the **zpa_branch_connector_group** data source to get information about branch connector group resources from ZTW shared within the Zscaler Private Access cloud. This data source can be used when configuring `zpa_policy_access_rule` or `zpa_policy_access_rule_v2`, `zpa_policy_forwarding_rule`, `zpa_policy_forwarding_rule_v2`, where the `object_type` is `BRANCH_CONNECTOR_GROUP`

## Example Usage

```terraform
data "zpa_branch_connector_group" "this" {
    name = "GROUP01"
}

data "zpa_branch_connector_group" "this" {
    id = "123635465"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (String) The name of the partner extranet to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The id of the partner extranet to be exported.
* `enabled` - (bool) If the partner extranet iis enabled.