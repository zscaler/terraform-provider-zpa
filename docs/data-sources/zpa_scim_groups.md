---
subcategory: "SCIM Groups"
layout: "zscaler"
page_title: "ZPA: scim_groups"
description: |-
  Get information about SCIM Group from an Identity Provider (IdP) in the Zscaler Private Access cloud.
---

# Data Source: zpa_scim_groups

Use the **zpa_scim_groups** data source to get information about a SCIM Group from an Identity Provider (IdP). This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

## Example Usage

```hcl
# ZPA SCIM Groups Data Source
data "zpa_scim_groups" "engineering" {
    name = "Engineering"
    idp_name = "idp_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the scim group to be exported.
* `idp_name` - (Required) Name. The name of the IdP where the scim group must be exported from.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `creation_time` - (string)
* `idp_id` - (string) The ID of the IdP corresponding to the SAML attribute.
* `idp_group_id`(string)
* `modified_time` (string)
