---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_access_rule"
description: |-
  Creates and manages ZPA Policy Access Rule with SAML Attribute conditions.
---

# Resource: zpa_policy_access_rule

The **zpa_policy_access_rule** resource creates and manages a policy access rule with SAML attribute conditions in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```hcl
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_IdP_Name"
}

resource "zpa_policy_access_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "ALLOW"
  operator = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id

  conditions {
     negated    = false
     operator   = "OR"
    operands {
      object_type = "SAML"
      lhs = data.zpa_saml_attribute.email_user_sso.id
      rhs = "user1@acme.com"
      idp_id = data.zpa_idp_controller.idp_name.id
    }
  }
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``

## Attributes Reference

* `action` (Optional) This is for providing the rule action. Supported values: ``ALLOW``, ``DENY``
* `custom_msg` (Optional) This is for providing a customer message for the user.
* `description` (Optional) This is the description of the access policy rule.
* `operator` (Optional) Supported values: ``AND``, ``OR``
* `policy_type` (Optional) Supported values: ``ACCESS_POLICY`` or ``GLOBAL_POLICY``
* `rule_order` - (Deprecated)

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

* `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

* `conditions` - (Optional)
  * `negated` - (Optional) Supported values: ``true`` or ``false``
  * `operator` (Optional) Supported values: ``AND``, and ``OR``
  * `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

  * `operands` (Optional) - Operands block must be repeated if multiple per `object_type` conditions are to be added to the rule.
    * `name` (Optional)
    * `lhs` (Optional)
    * `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.
    * `idp_id` (Optional) The ID information from an existing IDP.
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `SAML`. Use [zpa_saml_attribute](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) data source to retrieve the SAML attribute ``id``.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is NOT supported within the `operands` block when the `object_type` is set to `SAML`. IDP Information is controller at the parent tenant level.

* `app_connector_groups`
  * `id` - (Optional) The ID of an app connector group resource

* `app_server_groups`
  * `id` - (Optional) The ID of a server group resource

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy Access Rule for Browser Access can be imported by using`<POLICY ACCESS RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example <policy_access_rule_id>
```

## LHS and RHS Values

LHS and RHS values differ based on object types. Refer to the following table:

| Object Type | LHS| RHS
|----------|-----------|----------
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment) | "id" | <application_segment_ID> |
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group) | "id" | <segment_group_ID> |
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment_browser_access) | "id" | zpn_client_type_zappl or zpn_client_type_exporter |
| [EDGE_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group) | "id" | <edge_connector_ID> |
| [IDP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_idp_controller) | "id" | <identity_provider_ID> |
| [MACHINE_GRP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group) | "id" | <machine_group_ID> |
| [POSTURE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile) | <posture_udid>  | "true" / "false" |
| [TRUSTED_NETWORK](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_trusted_network) | <network_id>  | "true" |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | <saml_attribute_id>  | <Attribute_value_to_match> |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | <scim_attribute_id>  | <Attribute_value_to_match>  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | <scim_group_attribute_id>  | <Attribute_value_to_match>  |
