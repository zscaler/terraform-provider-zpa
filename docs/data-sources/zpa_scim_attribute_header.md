---
page_title: "zpa_scim_attribute_header Data Source - terraform-provider-zpa"
subcategory: "SCIM Attribute Header"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-scim
  API documentation https://help.zscaler.com/zpa/obtaining-scim-attribute-details-using-api
  Get information about SCIM attributes from an Identity Provider (IdP) in the Zscaler Private Access cloud.
---

# zpa_scim_attribute_header (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-scim)
* [API documentation](https://help.zscaler.com/zpa/obtaining-scim-attribute-details-using-api)

Use the **zpa_scim_attribute_header** data source to get information about a SCIM attribute from an Identity Provider (IdP). This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Inspection Policy.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
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

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the scim attribute header to be exported.
* `idp_name` - (Required) The name of the scim attribute header that must be exported.

### Read-Only

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
