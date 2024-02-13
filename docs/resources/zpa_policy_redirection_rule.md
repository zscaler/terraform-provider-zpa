---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: zpa_policy_redirection_rule"
description: |-
  Creates and manages ZPA Policy Access Redirection Rule.
---

# Resource: zpa_policy_redirection_rule

The **zpa_policy_redirection_rule** resource creates and manages policy access redirection rule in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```hcl
# Get Redirection Access Policy ID
data "zpa_policy_type" "this" {
  policy_type = "REDIRECTION_POLICY"
}

# Get Service Edge Group ID
data "zpa_service_edge_group" "this" {
 name = "Example"
}

#Create Policy Access Rule
resource "zpa_policy_redirection_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "REDIRECT_ALWAYS"
  operator = "AND"
  policy_set_id = data.zpa_policy_type.this.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs         = "id"
      rhs         = "zpn_client_type_branch_connector"
    }
    operands {
      object_type = "CLIENT_TYPE"
      lhs         = "id"
      rhs         = "zpn_client_type_edge_connector"
   }
  }
  service_edge_groups {
    id = [ data.zpa_service_edge_group.this.id ]
  }
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``

## Attributes Reference

* `action` (Optional) This is for providing the rule action. Supported values: ``REDIRECT_DEFAULT``, ``REDIRECT_PREFERRED``, and ``REDIRECT_ALWAYS``
* `description` (Optional) This is the description of the access policy rule.
* `operator` (Optional) Supported values: ``AND``, ``OR``
* `rule_order` - (Deprecated)

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

* `conditions` - (Optional)
  * `negated` - (Optional) Supported values: ``true`` or ``false``
  * `operator` (Optional) Supported values: ``AND``, and ``OR``

  * `operands` (Optional) - Operands block must be repeated if multiple per `object_type` conditions are to be added to the rule.
    * `name` (Optional)
    * `lhs` (Optional) LHS must always carry the string value ``id`` or the attribute ID of the resource being associated with the rule.
    * `rhs` (Optional) RHS is either the ID attribute of a resource or fixed string value. Refer to the chart below for further details.
    * `idp_id` (Optional)
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `CLIENT_TYPE`
    * `CLIENT_TYPE` (Optional) - The below options are the only ones supported in an access policy rule.
      * `zpn_client_type_machine_tunnel`
      * `zpn_client_type_edge_connector`
      * `zpn_client_type_zapp`
      * `zpn_client_type_branch_connector`

* `service_edge_groups`
  * `id` - (Optional) The ID of an service edge group resource

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy access rule can be imported by using `<POLICY ACCESS REDIRECTION RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_redirection_rule.example <policy_access_rule_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS
|----------|-----------|----------
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types) | ``"id"`` | ``zpn_client_type_machine_tunnel``, ``zpn_client_type_edge_connector``, ``zpn_client_type_zapp``, ``zpn_client_type_branch_connector``  |
