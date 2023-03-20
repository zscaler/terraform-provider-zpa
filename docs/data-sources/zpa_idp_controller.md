---
subcategory: "Identity Provider"
layout: "zscaler"
page_title: "ZPA: idp_controller"
description: |-
  Get information about an Identity Provider in Zscaler Private Access cloud.
---

# Data Source: zpa_idp_controller

Use the **zpa_idp_controller** data source to get information about an Identity Provider created in the Zscaler Private Access cloud. This data source is required when creating:

1. Access policy Rules
2. Access policy timeout rules
3. Access policy forwarding rules
4. Access policy inspection rules
5. Access policy isolation rules

## Example Usage

```hcl
# ZPA IdP Controller Data Source
data "zpa_idp_controller" "example" {
 name = "idp_name"
}
```

```hcl
# ZPA IdP Controller Data Source
data "zpa_idp_controller" "example" {
 id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Identity Provider (IdP) to be exported.
* `id` - (Optional) The name of the Identity Provider (IdP) to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `auto_provision` (string)
* `creation_time` (string)
* `description` (string)
* `disable_saml_based_policy` (bool)
* `domain_list` (string)
* `enable_scim_based_policy` (bool)
* `enabled` (bool) Default value if null is True
* `idp_entity_id` (string)
* `login_name_attribute` (string)
* `login_url` (string)
* `login_hint` (bool)
* `force_auth` (bool)
* `enable_arbitrary_auth_domains` (string)
* `modified_time` (string)
* `modified_by` (string)
* `reauth_on_user_update` (bool)
* `redirect_binding` (bool)
* `scim_enabled` (bool)
* `scim_service_provider_endpoint` (string)
* `scim_shared_secret_exists` (bool)
* `sign_saml_request` (string)
* `sso_type` (string)
* `use_custom_sp_metadata` (bool)

* `user_metadata` (Computed)
  * `certificate_url` (string)
  * `sp_entity_id` (string)
  * `sp_metadata_url` (string)
  * `sp_post_url` (string)

* `admin_metadata` (Computed)
  * `certificate_url` (string)
  * `sp_entity_id` (string)
  * `sp_metadata_url` (string)
  * `sp_post_url` (string)

:warning: Notice that certificate and public_keys are omitted from the output.
