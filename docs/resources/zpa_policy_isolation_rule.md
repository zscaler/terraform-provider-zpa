---
page_title: "zpa_policy_isolation_rule Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-isolation-policy
  API documentation https://help.zscaler.com/zpa/configuring-isolation-policies-using-api
  Creates and manages ZPA Policy Access Isolation Rule.
---

# zpa_policy_isolation_rule (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-isolation-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-isolation-policies-using-api)

The **zpa_policy_isolation_rule** resource creates a policy isolation access rule in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
# Get Isolation Profile ID
data "zpa_isolation_profile" "isolation_profile" {
    name = "zpa_isolation_profile"
}

#Create Client Isolation Access Rule
resource "zpa_policy_isolation_rule" "this" {
  name                          = "Example_Isolation_Policy"
  description                   = "Example_Isolation_Policy"
  action                        = "ISOLATE"
  operator = "AND"
  zpn_isolation_profile_id = data.zpa_isolation_profile.isolation_profile.id

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs = "id"
      rhs = "zpn_client_type_exporter"
    }
  }
}
```

## Schema

### Required

- `name` (String) This is the name of the forwarding policy rule.
- `action` (String) This is for providing the rule action.
  * The supported actions for a policy isolation rule are: ``BYPASS_ISOLATE``, or ``ISOLATE``
- `zpn_isolation_profile_id` (String) Use [zpa_isolation_profile](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_isolation_profile) data source to retrieve the necessary Isolation profile ID ``zpn_isolation_profile_id``

### Optional

- `policy_set_id` - (String) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``
    ~> **NOTE** As of v3.2.0 the ``policy_set_id`` attribute is now optional, and will be automatically determined based on the policy type being configured. The attribute is being kept for backwards compatibility, but can be safely removed from existing configurations.
- `description` - (String) This is the description of the access policy rule.
- `operator` (String) Supported values: ``AND``, ``OR``
- `rule_order` (String, Deprecated)

    ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set)
Specifies the set of conditions for the policy rule.
  - `operator` (String) Supported values: ``AND``, and ``OR``
  - `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

  - `operands` (Block Set) - Operands block must be repeated if multiple per `object_type` conditions are to be added to the rule.
    - `lhs` (String) LHS must always carry the string value ``id`` or the attribute ID of the resource being associated with the rule.
    - `rhs` (String) RHS is either the ID attribute of a resource or fixed string value. Refer to the chart below for further details.
    - `idp_id` (String)
    - `object_type` (String) This is for specifying the policy critiera. Supported values: `APP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `CLOUD_CONNECTOR_GROUP`. `TRUSTED_NETWORK`, and `CLIENT_TYPE`.
    - `CLIENT_TYPE` (String) - The below options are the only ones supported in a timeout policy rule.
      - ``zpn_client_type_exporter`` "Web Browser"
    - `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

    ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

Policy Access Isolation Rule can be imported by using `<POLICY ISOLATION RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_isolation_rule.example <policy_isolation_rule_id>
```

## LHS and RHS Values

LHS and RHS values differ based on object types. Refer to the following table:

| Object Type | LHS| RHS
|----------|-----------|----------
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment) | ``"id"`` | ``application_segment_id`` |
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group) | ``"id"`` | ``segment_group_id``|
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types) | ``"id"`` | ``zpn_client_type_exporter`` |
| [PLATFORM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``mac``, ``ios``, ``windows``, ``android``, ``linux`` | ``"true"`` / ``"false"`` |
| [EDGE_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group) | ``"id"`` | ``edge_connector_id`` |
| [IDP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_idp_controller) | ``"id"`` | ``identity_provider_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | <Attribute_value_to_match> |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | <Attribute_value_to_match>  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | <Attribute_value_to_match>  |
| [MACHINE_GRP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group) | ``"id"`` | ``machine_group_id`` |
| [POSTURE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile) | ``posture_udid``  | ``"true"`` / ``"false"`` |
| [TRUSTED_NETWORK](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_trusted_network) | ``network_id``  | ``"true"`` |
