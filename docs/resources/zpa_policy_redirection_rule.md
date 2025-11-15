---
page_title: "zpa_policy_redirection_rule Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-redirection-policy
  API documentation https://help.zscaler.com/zpa/configuring-redirection-policies-using-api
  Creates and manages ZPA Policy Access Redirection Rule.
---

# zpa_policy_redirection_rule (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-redirection-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-redirection-policies-using-api)

The **zpa_policy_redirection_rule** resource creates a policy redirection access rule in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage - REDIRECT_DEFAULT

```terraform
resource "zpa_policy_redirection_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "REDIRECT_DEFAULT"

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      values      = ["zpn_client_type_branch_connector"]
    }
  }
}
```

## Example Usage - REDIRECT_PREFERRED

```terraform
data "zpa_service_edge_group" "this" {
    name = "Service_Edge_Group01
}

resource "zpa_policy_redirection_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "REDIRECT_PREFERRED"

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      values      = ["zpn_client_type_branch_connector"]
    }
  }
  service_edge_groups {
    id = [ data.zpa_service_edge_group.this.id ]
  }
}
```

## Example Usage - REDIRECT_ALWAYS

```terraform
data "zpa_service_edge_group" "this" {
    name = "Service_Edge_Group01
}

resource "zpa_policy_redirection_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "REDIRECT_ALWAYS"

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      values      = ["zpn_client_type_branch_connector"]
    }
  }
  service_edge_groups {
    id = [ data.zpa_service_edge_group.this.id ]
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
    - `object_type` (String) This is for specifying the policy critiera. Supported values: `CLIENT_TYPE`, `COUNTRY_CODE`.
        - `CLIENT_TYPE` (String) - The below options are the only ones supported in a timeout policy rule.
            - ``zpn_client_type_machine_tunnel``
            - ``zpn_client_type_edge_connector``
            - ``zpn_client_type_zapp``
            - ``zpn_client_type_zapp_partner``
            - ``zpn_client_type_branch_connector``

        - `COUNTRY_CODE` (String) - Use a standard 2 letter `ISO3166 Alpha2` Country codes. See list [here](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)
    - `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

    ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

Policy Access Isolation Rule can be imported by using `<POLICY REDIRECTION RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_isolation_rule.example <rule_id>
```
