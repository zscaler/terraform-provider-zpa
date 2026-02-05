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

### Basic Usage

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

### Using SCIM Attributes in Policy Rules

```terraform
# Fetch SCIM attribute header
data "zpa_scim_attribute_header" "display_name" {
  name     = "DisplayName"
  idp_name = "IdP_Name"
}

# Output available values (optional, useful for debugging)
output "display_name_values" {
  value = data.zpa_scim_attribute_header.display_name.values
}

# Use in policy access rule
resource "zpa_policy_access_rule" "scim_rule" {
  name        = "SCIM-based Access Rule"
  description = "Allow access based on SCIM DisplayName attribute"
  action      = "ALLOW"
  operator    = "AND"

  conditions {
    operator = "OR"
    operands {
      object_type = "SCIM"
      idp_id      = data.zpa_scim_attribute_header.display_name.idp_id
      lhs         = data.zpa_scim_attribute_header.display_name.id
      # The rhs value must match one of the values from the values list
      rhs         = "John Smith"
    }
  }
}
```

**Note**: When using SCIM attributes in policy rules, the `rhs` value must exactly match one of the values available in the `values` attribute. Use the output to see all available values for the SCIM attribute.

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the scim attribute header to be exported.
* `idp_name` - (Required) The name of the scim attribute header that must be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `canonical_values` - (List of String) Canonical values for the SCIM attribute.
* `case_sensitive` - (Boolean) Whether the attribute is case-sensitive.
* `creation_time` - (String) The time when the SCIM attribute was created.
* `data_type` - (String) The data type of the SCIM attribute (e.g., "String", "Boolean").
* `description` - (String) Description of the SCIM attribute.
* `id` - (String) The unique identifier of the SCIM attribute header.
* `idp_id` - (String) The ID of the IdP corresponding to the SCIM attribute.
* `idp_name` - (String) The name of the IdP corresponding to the SCIM attribute.
* `modified_by` - (String) The ID of the user who last modified the SCIM attribute.
* `modified_time` - (String) The time when the SCIM attribute was last modified.
* `multivalued` - (Boolean) Whether the attribute can have multiple values.
* `mutability` - (String) Indicates whether the attribute can be modified.
* `required` - (Boolean) Whether the attribute is required.
* `returned` - (String) Indicates when the attribute is returned (e.g., "default", "always", "never").
* `schema_uri` - (String) The schema URI for the SCIM attribute.
* `uniqueness` - (Boolean) Whether the attribute values must be unique.
* `values` - (Set of String) **Important**: List of all available values for this SCIM attribute. Use these values when referencing the attribute in policy rules. The `rhs` field in policy conditions must match one of these values exactly.
