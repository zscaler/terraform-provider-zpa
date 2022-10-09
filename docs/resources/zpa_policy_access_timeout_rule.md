---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_timeout_rule"
description: |-
  Creates and manages ZPA Policy Timeout Access Rule.
---

# Resource: zpa_policy_timeout_rule

The **zpa_policy_timeout_rule** resource creates a policy timeout rule in the Zscaler Private Access cloud.

## Example Usage

```hcl
# Get Global Timeout Policy ID
data "zpa_policy_type" "timeout_policy" {
    policy_type = "TIMEOUT_POLICY"
}

# Get IdP ID
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

# Get SCIM Group attribute ID
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP_Name"
}

resource "zpa_policy_timeout_rule" "this"  {
  name                          = "Example"
  description                   = "Example"
  action                        = "RE_AUTH"
  reauth_idle_timeout           = "600"
  reauth_timeout                = "172800"
  operator                      = "AND"
  policy_set_id                 = data.zpa_policy_type.timeout_policy.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs = "id"
      rhs = "zpn_client_type_exporter"
    }
  }
  conditions {
    negated  = false
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [data.zpa_scim_groups.engineering.id]
    }
  }
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``
* `reauth_timeout` (Required) This denotes the authentication timeout. Provides the timeout value in seconds. -1 value denotes Never.
* `reauth_idle_timeout` (Required) This denotes the idle connection timeout. Provides the timeout value in seconds. -1 value denotes Default.

## Attributes Reference

* `action` (Optional) This is for providing the rule action. Supported value: ``RE_AUTH``
* `custom_msg` (Optional) This is for providing a customer message for the user.
* `description` (Optional) This is the description of the access policy rule.
* `operator` (Optional) Supported values: ``AND``, and ``OR``
* `policy_type` (Optional) Supported values: ``TIMEOUT_POLICY`` or ``REAUTH_POLICY``

* `rule_order` (Optional)

* `conditions` - (Optional)
  * `negated` - (Optional) Supported values: ``true`` or ``false``
  * `operator` (Optional) Supported values: ``AND``, and ``OR``
  * `operands` (Optional) - Operands block must be repeated if multiple per `object_type` conditions are to be added to the rule.
    * `name` (Optional)
    * `lhs` (Optional) LHS must always carry the string value ``id`` or the attribute ID of the resource being associated with the rule.
    * `rhs` (Optional) RHS is either the ID attribute of a resource or fixed string value. Refer to the chart below for further details.
    * `idp_id` (Optional)
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `SAML`, `SCIM`, `SCIM_GROUP`, `IDP`, `CLIENT_TYPE`,  `POSTURE`
    * `CLIENT_TYPE` (Optional) - The below options are the only ones supported in a timeout policy rule.
      * `zpn_client_type_zapp`
      * `zpn_client_type_browser_isolation`
      * `zpn_client_type_exporter`

## Import

Policy access timeout can be imported by using `<POLICY TIMEOUT RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_timeout_rule.example <policy_timeout_rule_id>
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
