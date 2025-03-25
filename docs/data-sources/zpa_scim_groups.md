---
page_title: "zpa_scim_groups Data Source - terraform-provider-zpa"
subcategory: "SCIM Groups"
layout: "zscaler"
page_title: "ZPA: scim_groups"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-scim-groups
  API documentation https://help.zscaler.com/zpa/obtaining-scim-group-details-using-api
  Get information about SCIM Group from an Identity Provider (IdP) in the Zscaler Private Access cloud.
---

# zpa_scim_groups (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-scim-groups)
* [API documentation](https://help.zscaler.com/zpa/obtaining-scim-group-details-using-api)

Use the **zpa_scim_groups** data source to get information about a SCIM Group from an Identity Provider (IdP). This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
# ZPA SCIM Groups Data Source
data "zpa_scim_groups" "engineering" {
    name = "Engineering"
    idp_name = "idp_name"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) Name. The name of the scim group to be exported.
* `idp_name` - (Required) Name. The name of the IdP where the scim group must be exported from.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `creation_time` - (string)
* `idp_id` - (string) The ID of the IdP corresponding to the SAML attribute.
* `idp_group_id`(string)
* `modified_time` (string)
