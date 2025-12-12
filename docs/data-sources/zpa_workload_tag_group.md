---
page_title: "zpa_workload_tag_group Data Source - terraform-provider-zpa"
subcategory: "Workload Tag Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/configuring-resource-groups
  API documentation https://help.zscaler.com/zpa/configuring-resource-groups
  Get information about partner extranet resources
---

# zpa_workload_tag_group (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/configuring-resource-groups)
* [API documentation](https://help.zscaler.com/zpa/configuring-resource-groups)

Use the **zpa_workload_tag_group** data source to get information about workload tag group in the Zscaler Private Access cloud. This data source can be used when configuring `zpa_policy_access_rule` or `zpa_policy_access_rule_v2`, `object_type` is `WORKLOAD_TAG_GROUP`

## Example Usage

```terraform
data "zpa_workload_tag_group" "this" {
    name = "Group01"
}

data "zpa_workload_tag_group" "this" {
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
