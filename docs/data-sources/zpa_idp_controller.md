---
subcategory: "IdP Controller"
layout: "zpa"
page_title: "ZPA: idp_controller"
description: |-
  Gets an Identity Provider (IdP) details.
---

# zpa_idp_controller

The **zpa_idp_controller** data source provides details about a specific Identity Provider created in the Zscaler Private Access cloud.
This data source is required when creating:

1. Access policy Rule
2. Access policy timeout rule
3. Access policy forwarding rule

## Example Usage

```hcl
# ZPA IdP Controller Data Source
data "zpa_idp_controller" "example" {
 name = "idp_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the Identity Provider (IdP) to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `auto_provision` (String)
* `creation_time` (String)
* `description` (String)
* `disable_saml_based_policy` (Boolean)
* `domain_list` (List of String)
* `enable_scim_based_policy` (Boolean)
* `enabled` (Boolean) Default value if null is True
* `idp_entity_id` (String)
* `login_name_attribute` (String)
* `login_url` (String)
* `modified_time` (String)
* `modifiedby` (String)
* `reauth_on_user_update` (Boolean)
* `redirect_binding` (Boolean)
* `scim_enabled` (Boolean)
* `scim_service_provider_endpoint` (String)
* `scim_shared_secret` (String)
* `scim_shared_secret_exists` (Boolean)
* `sign_saml_request` (String)
* `sso_type` (List of String)
* `use_custom_sp_metadata` (Boolean)

`user_metadata` (Set of Object)

* `certificate_url` (String)
* `sp_entity_id` (String)
* `sp_metadata_url` (String)
* `sp_post_url` (String)

`admin_metadata` (Set of Object)

* `certificate_url` (String)
* `sp_entity_id` (String)
* `sp_metadata_url` (String)
* `sp_post_url` (String)

* `certificates` (List of Object)
