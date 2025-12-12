---
page_title: "zpa_extranet_resource_partner Data Source - terraform-provider-zpa"
subcategory: "Enrollment Certificate"
description: |-
  Official documentation https://help.zscaler.com/zpa/viewing-extranet-dashboard
  API documentation https://help.zscaler.com/zpa/viewing-extranet-dashboard
  Get information about partner extranet resources
---

# zpa_extranet_resource_partner (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/viewing-extranet-dashboard)
* [API documentation](https://help.zscaler.com/zpa/viewing-extranet-dashboard)

Use the **zpa_extranet_resource_partner** data source to get information about partner extranet resources in the Zscaler Private Access cloud. This data source is required when configuring resources such as: `zpa_server_group`, `zpa_application_segment`, `zpa_application_segmnent_pra`, `zpa_policy_access_rule_v2` in [Extranet mode](https://help.zscaler.com/zia/about-extranet)

## Example Usage

```terraform
data "zpa_extranet_resource_partner" "this" {
    name = "Extranet01"
}

data "zpa_extranet_resource_partner" "this" {
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