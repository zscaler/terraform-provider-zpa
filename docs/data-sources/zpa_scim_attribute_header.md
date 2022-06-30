---
subcategory: "SCIM Attribute Header"
layout: "zscaler"
page_title: "ZPA: scim_attribute_header"
description: |-
  Get information about SCIM attributes from an Identity Provider (IdP) in the Zscaler Private Access cloud.
---

# Data Source: zpa_scim_attribute_header

Use the **zpa_scim_attribute_header** data source to get information about a SCIM attribute from an Identity Provider (IdP). This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Inspection Policy.

## Example Usage

```hcl
# ZPA SCIM Attribute Header Data Source
data "zpa_scim_attribute_header" "givenName" {
  name     = "name.givenName"
  idp_name = "IdP_Name"
}

data "zpa_scim_attribute_header" "familyName" {
  name     = "name.familyName"
  idp_name = "IdP_Name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the scim attribute header to be exported.
* `idp_name` - (Required) The name of the scim attribute header that must be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `canonical_values` - (string)
* `case_sensitive` - (bool)
* `creation_time` - (string)
* `data_type` - (string)
* `description` - (string)
* `id` - (string)
* `idp_id` - (string) The ID of the IdP corresponding to the SAML attribute.
* `modified_by`(string)
* `modified_time` (string)
* `multivalued` (bool)
* `mutability` (string)
* `required` (bool)
* `returned` (string)
* `schema_uri` (string)
* `uniqueness` (bool)
