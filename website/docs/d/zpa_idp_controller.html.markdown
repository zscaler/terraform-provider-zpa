---
layout: "zscaler"
page_title: "Zscaler Private Access (ZPA): idp_controller"
sidebar_current: "docs-datasource-zpa-idp-controller"
description: |-
  Get information about an Identity Provider in Zscaler Private Access cloud.
---

# zpa_idp_controller

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

* `auto_provision` (Computed)
* `creation_time` (Computed)
* `description` (Computed)
* `disable_saml_based_policy` (Computed)
* `domain_list` (Computed)
* `enable_scim_based_policy` (Computed)
* `enabled` (Computed) Default value if null is True
* `idp_entity_id` (Computed)
* `login_name_attribute` (Computed)
* `login_url` (Computed)
* `modified_time` (Computed)
* `modifiedby` (Computed)
* `reauth_on_user_update` (Computed)
* `redirect_binding` (Computed)
* `scim_enabled` (Computed)
* `scim_service_provider_endpoint` (Computed)
* `scim_shared_secret` (Computed)
* `scim_shared_secret_exists` (Computed)
* `sign_saml_request` (Computed)
* `sso_type` (Computed)
* `use_custom_sp_metadata` (Computed)

* `user_metadata` (Computed)
  * `certificate_url` (Computed)
  * `sp_entity_id` (Computed)
  * `sp_metadata_url` (Computed)
  * `sp_post_url` (Computed)

* `admin_metadata` (Computed)
  * `certificate_url` (Computed)
  * `sp_entity_id` (Computed)
  * `sp_metadata_url` (Computed)
  * `sp_post_url` (Computed)

* `certificates` (Computed)

:warning: Notice that certificate and public_keys are omitted from the output.
