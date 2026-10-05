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

## Example Usage - Groups with the same name under different IAM IdPs

With ZIdentity, several IAM IdPs can map to the same ZPA IdP, and each can provide a group with the same name. Use `iam_idp_name` or `iam_idp_id` to select the group from a specific IAM IdP:

```terraform
data "zpa_scim_groups" "engineering_okta" {
    name         = "Engineering"
    idp_name     = "idp_name"
    iam_idp_name = "Okta"
}

data "zpa_scim_groups" "engineering_entra" {
    name         = "Engineering"
    idp_name     = "idp_name"
    iam_idp_name = "Entra ID"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) Name. The name of the scim group to be exported.
* `idp_name` - (Required) Name. The name of the IdP where the scim group must be exported from.

### Optional

* `iam_idp_name` - (Optional) Name of the ZIdentity (IAM) IdP the group belongs to. Use it when groups with the same `name` exist under different IAM IdPs mapped to the same ZPA IdP. Requires `name`; conflicts with `iam_idp_id`.
* `iam_idp_id` - (Optional) ID of the ZIdentity (IAM) IdP the group belongs to. Same purpose as `iam_idp_name`. Requires `name`; conflicts with `iam_idp_name`.

-> **NOTE:** When `iam_idp_name` or `iam_idp_id` is set, the lookup fails if more than one group matches, or if none does, and the error lists the candidate groups (`id`, `name`, `iamIdpId`, `iamIdpName`). Without these arguments, the lookup behaves as before and returns the first group whose name matches.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `creation_time` - (string)
* `idp_id` - (string) The ID of the IdP corresponding to the SAML attribute.
* `idp_group_id`(string)
* `iam_idp_id` - (string) The ID of the ZIdentity (IAM) IdP the group belongs to, when returned by the API.
* `iam_idp_name` - (string) The name of the ZIdentity (IAM) IdP the group belongs to, when returned by the API.
* `modified_time` (string)
