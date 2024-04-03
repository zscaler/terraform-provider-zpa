---
page_title: "zpa_policy_access_rule Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-access-policy
  API documentation https://help.zscaler.com/zpa/configuring-access-policies-using-api
  Creates and manages ZPA Policy Access Rule with SCIM Group conditions.
---

# zpa_policy_access_rule (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-access-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api)

The **zpa_policy_access_rule** resource creates and manages a policy access rule with SCIM Group conditions in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP_Name"
}

resource "zpa_policy_access_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "ALLOW"
  operator                      = "AND"

  conditions {
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs         = data.zpa_idp_controller.idp_name.id
      rhs         = data.zpa_scim_groups.engineering.id
    }
  }
}
```

## Schema

### Required

- `name` (String) This is the name of the policy rule.

### Optional

* `policy_set_id` - (String) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``
    ~> **NOTE** As of v3.2.0 the ``policy_set_id`` attribute is now optional, and will be automatically determined based on the policy type being configured. The attribute is being kept for backwards compatibility, but can be safely removed from existing configurations.

- `action` (String) This is for providing the rule action. Supported values: ``ALLOW``, ``DENY``
- `custom_msg` (String) This is for providing a customer message for the user.
- `description` (String) This is the description of the access policy rule.
- `operator` (String) Supported values: ``AND``, ``OR``
- `policy_type` (String) Supported values: ``ACCESS_POLICY`` or ``GLOBAL_POLICY``
- `rule_order` (String, Deprecated)

    ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set) 
  - `operator` (String) Supported values: ``AND``, and ``OR``
  - `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

  - `operands` (Block Set) Operands block must be repeated if multiple per `object_type` conditions are to be added to the rule.`object_type` conditions are to be added to the rule.
    - `lhs` (String) The key for the object type.
    - `rhs` (String) This denotes the value for the given object type. Its value depends upon the key.
    - `idp_id` (String) The unique identifier of the IdP.
    - `object_type` (String) This is for specifying the policy critiera. Supported values: `SCIM_GROUP`. Use [zpa_scim_groups](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) data source to retrieve the SCIM Group ``id``.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is NOT supported within the `operands` block when the `object_type` is set to `SCIM_GROUP`. IDP Information is controller at the parent tenant level.

- `app_connector_groups` (Block Set)
  * `id` (String) The ID of an app connector group resource

- `app_server_groups` (Block Set)
  * `id` (String) The ID of a server group resource

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy Access Rule for Browser Access can be imported by using`<POLICY ACCESS RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example <policy_access_rule_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS
|----------|-----------|----------
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment) | ``"id"`` | ``application_segment_id`` |
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group) | ``"id"`` | ``segment_group_id``|
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types) | ``"id"`` | ``zpn_client_type_zappl``, ``zpn_client_type_exporter``, ``zpn_client_type_browser_isolation``, ``zpn_client_type_ip_anchoring``, ``zpn_client_type_edge_connector``, ``zpn_client_type_branch_connector``,  ``zpn_client_type_zapp_partner``, ``zpn_client_type_zapp``  |
| [EDGE_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group) | ``"id"`` | ``<edge_connector_id>`` |
| [IDP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_idp_controller) | ``"id"`` | ``identity_provider_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
| [PLATFORM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``mac``, ``ios``, ``windows``, ``android``, ``linux`` | ``"true"`` / ``"false"`` |
| [MACHINE_GRP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group) | ``"id"`` | ``machine_group_id`` |
| [POSTURE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile) | ``posture_udid``  | ``"true"`` / ``"false"`` |
| [TRUSTED_NETWORK](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_trusted_network) | ``network_id``  | ``"true"`` |
| [COUNTRY_CODE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_platforms) | [2 Letter ISO3166 Alpha2](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)  | ``"true"`` / ``"false"`` |
