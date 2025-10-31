---
page_title: "zpa_location_controller_summary_summary Data Source - terraform-provider-zpa"
subcategory: "Location Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/configuring-access-policies-using-api
  API documentation https://help.zscaler.com/zpa/policy-management#/mgmtconfig/v2/admin/customers/{customerId}/policySet/{policySetId}/rule-post
  Get information about location resource
---

# zpa_location_controller_summary (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api)
* [API documentation](https://help.zscaler.com/zpa/policy-management#/mgmtconfig/v2/admin/customers/{customerId}/policySet/{policySetId}/rule-post)

Use the **zpa_location_controller_summary** data source to get information about location resources from ZIA shared within the Zscaler Private Access cloud. This data source can be used when configuring `zpa_policy_access_rule` or `zpa_policy_access_rule_v2` where the `object_type` is `LOCATION`

## Example Usage

```terraform
data "zpa_location_controller_summary" "this" {
  name        = "ExtranetLocation01 | zscalerbeta.net"
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