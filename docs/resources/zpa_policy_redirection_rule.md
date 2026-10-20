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
  name        = "Example"
  description = "Example"
  action      = "REDIRECT_DEFAULT"

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs         = "id"
      rhs         = "zpn_client_type_branch_connector"
    }
  }
}
```

## Example Usage - REDIRECT_PREFERRED

```terraform
data "zpa_service_edge_group" "this" {
  name = "Service_Edge_Group01"
}

resource "zpa_policy_redirection_rule" "this" {
  name        = "Example"
  description = "Example"
  action      = "REDIRECT_PREFERRED"

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs         = "id"
      rhs         = "zpn_client_type_branch_connector"
    }
  }
  service_edge_groups {
    id = [data.zpa_service_edge_group.this.id]
  }
}
```

## Example Usage - REDIRECT_ALWAYS

```terraform
data "zpa_service_edge_group" "this" {
  name = "Service_Edge_Group01"
}

resource "zpa_policy_redirection_rule" "this" {
  name        = "Example"
  description = "Example"
  action      = "REDIRECT_ALWAYS"

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      lhs         = "id"
      rhs         = "zpn_client_type_branch_connector"
    }
  }
  service_edge_groups {
    id = [data.zpa_service_edge_group.this.id]
  }
}
```

## Example Usage - IdP, SCIM group and country conditions

```terraform
data "zpa_idp_controller" "this" {
  name = "IdP_Name"
}

data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "IdP_Name"
}

data "zpa_service_edge_group" "this" {
  name = "Service_Edge_Group01"
}

resource "zpa_policy_redirection_rule" "this" {
  name     = "Example"
  action   = "REDIRECT_PREFERRED"
  operator = "AND"

  conditions {
    operator = "OR"
    operands {
      object_type = "IDP"
      lhs         = "id"
      rhs         = data.zpa_idp_controller.this.id
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs         = data.zpa_idp_controller.this.id
      rhs         = data.zpa_scim_groups.engineering.id
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "COUNTRY_CODE"
      lhs         = "US"
      rhs         = "true"
    }
  }
  service_edge_groups {
    id = [data.zpa_service_edge_group.this.id]
  }
}
```

## Schema

### Required

- `name` (String) This is the name of the redirection policy rule.

### Optional

- `action` (String) This is for providing the rule action. Supported values: ``REDIRECT_DEFAULT``, ``REDIRECT_PREFERRED`` and ``REDIRECT_ALWAYS``
- `service_edge_groups` (Block List) The Private Service Edge groups to redirect to.
  * `id` (Set of String) The IDs of the service edge groups.

  ~> **NOTE** At least one service edge group is required when `action` is ``REDIRECT_PREFERRED`` or ``REDIRECT_ALWAYS``, and none may be set when `action` is ``REDIRECT_DEFAULT``. The provider returns an error before calling the API otherwise.

- `policy_set_id` - (String) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``
    ~> **NOTE** As of v3.2.0 the ``policy_set_id`` attribute is now optional, and will be automatically determined based on the policy type being configured. The attribute is being kept for backwards compatibility, but can be safely removed from existing configurations.
- `description` - (String) This is the description of the redirection policy rule.
- `operator` (String) Supported values: ``AND``, ``OR``
- `rule_order` (String, Deprecated)

    ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``zpa_policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block List)
Specifies the set of conditions for the policy rule.
  - `operator` (String, Required) Supported values: ``AND``, and ``OR``
  - `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

  - `operands` (Block List) - Operands block must be repeated if multiple per `object_type` conditions are to be added to the rule.
    - `object_type` (String, Required) This is for specifying the policy criteria. Supported values: `CLIENT_TYPE`, `IDP`, `SAML`, `SCIM`, `SCIM_GROUP`, `COUNTRY_CODE`. See [LHS and RHS Values](#lhs-and-rhs-values) for what each one expects.
        - `CLIENT_TYPE` (String) - The below options are the only ones supported in a redirection policy rule.
            - ``zpn_client_type_machine_tunnel``
            - ``zpn_client_type_edge_connector``
            - ``zpn_client_type_zapp``
            - ``zpn_client_type_zapp_partner``
            - ``zpn_client_type_branch_connector``

        - `COUNTRY_CODE` (String) - Use a standard 2 letter `ISO3166 Alpha2` Country codes. See list [here](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)
    - `lhs` (String, Required) The key of the criterion: the string ``id`` or the attribute ID of the resource being associated with the rule. Refer to the chart below for further details.
    - `rhs` (String) The value of the criterion: the ID of a resource or a fixed string value. Refer to the chart below for further details.
    - `rhs_list` (Set of String) A list of values for the same `lhs`; each value is sent as its own operand. Ignored when `rhs` is set.
    - `idp_id` (String) The IdP ID. Required when `object_type` is `SCIM`.
    - `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

    ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy redirection rule can be imported by using `<POLICY REDIRECTION RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_redirection_rule.example <rule_id>
```

## LHS and RHS Values

| Object Type | LHS | RHS |
|----------|-----------|----------|
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types) | ``"id"`` | ``zpn_client_type_machine_tunnel``, ``zpn_client_type_edge_connector``, ``zpn_client_type_zapp``, ``zpn_client_type_zapp_partner``, ``zpn_client_type_branch_connector`` |
| [IDP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_idp_controller) | ``"id"`` | ``identity_provider_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id`` | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id`` (also set `idp_id`) | ``attribute_value_to_match`` |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``idp_id`` | ``scim_group_id`` |
| [COUNTRY_CODE](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes) | [2 Letter ISO3166 Alpha2](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes) | ``"true"`` |
