---
page_title: "zpa_policy_access_rule_v2 Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller V2"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-access-policy
  API documentation https://help.zscaler.com/zpa/configuring-access-policies-using-api#postV2
  Creates and manages ZPA Policy Access Rule via API v2 endpoints.
---

# zpa_policy_access_rule_v2 (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-access-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api#postV2)

The **zpa_policy_access_rule_v2** resource creates and manages policy access rule in the Zscaler Private Access cloud using a new v2 API endpoint.

  ⚠️ **NOTE**: This resource is recommended if your configuration requires the association of more than 1000 resource criteria per rule.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
# Retrieve Policy Types
# Retrieve Identity Provider ID
data "zpa_idp_controller" "this" {
	name = "Idp_Name"
}

# Retrieve SAML Attribute ID
data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_Users"
    idp_name = "Idp_Name"
}

# Retrieve SAML Attribute ID
data "zpa_saml_attribute" "group_user" {
    name = "GroupName_Users"
    idp_name = "Idp_Name"
}

# Retrieve SCIM Group ID
data "zpa_scim_groups" "a000" {
    name = "A000"
    idp_name = "Idp_Name"
}

# Retrieve SCIM Group ID
data "zpa_scim_groups" "b000" {
    name = "B000"
    idp_name = "Idp_Name"
}

# Create Segment Group
resource "zpa_segment_group" "this" {
   name = "Example"
   description = "Example"
   enabled = true
 }

# Create Policy Access Rule V2
resource "zpa_policy_access_rule_v2" "this" {
  name          = "Example"
  description   = "Example"
  action        = "ALLOW"

  conditions {
    operator = "OR"
    operands {
      object_type = "APP_GROUP"
      values      = [zpa_segment_group.this.id]
    }
  }

  conditions {
    operator = "OR"
    operands {
      object_type = "SAML"
      entry_values {
        rhs = "user1@acme.com"
        lhs = data.zpa_saml_attribute.email_user_sso.id
      }
      entry_values {
        rhs = "A000"
        lhs = data.zpa_saml_attribute.group_user.id
      }
    }
    operands {
      object_type = "SCIM_GROUP"
      entry_values {
        rhs = data.zpa_scim_groups.a000.id
        lhs = data.zpa_idp_controller.this.id
      }
      entry_values {
        rhs = data.zpa_scim_groups.b000.id
        lhs = data.zpa_idp_controller.this.id
      }
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "PLATFORM"
      entry_values {
        rhs = "true"
        lhs = "linux"
      }
      entry_values {
        rhs = "true"
        lhs = "android"
      }
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "COUNTRY_CODE"
      entry_values {
        lhs = "CA"
        rhs = "true"
      }
      entry_values {
        lhs = "US"
        rhs = "true"
      }
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "RISK_FACTOR_TYPE"
      entry_values {
        lhs = "ZIA"
        rhs = "UNKNOWN"
      }
      entry_values {
        lhs = "ZIA"
        rhs = "LOW"
      }
      entry_values {
        lhs = "ZIA"
        rhs = "MEDIUM"
      }
      entry_values {
        lhs = "ZIA"
        rhs = "HIGH"
      }
      entry_values {
        lhs = "ZIA"
        rhs = "CRITICAL"
      }
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "CHROME_ENTERPRISE"
      entry_values {
        lhs = "managed"
        rhs = "true"
      }
      entry_values {
        lhs = "managed"
        rhs = "false"
      }
    }
  }
}
```

## Example Usage - Configure Extranet Access Rule

```hcl
data "zpa_location_controller" "this" {
  name        = "ExtranetLocation01 | zscalerbeta.net"
  zia_er_name = "NewExtranet 8432"
}

data "zpa_location_group_controller" "this" {
  location_name = "ExtranetLocation01"
  zia_er_name   = "NewExtranet 8432"
}

data "zpa_extranet_resource_partner" "this" {
  name = "NewExtranet 8432"
}

resource "zpa_policy_access_rule_v2" "this" {
  name             = "Extranet_Rule01"
  description      = "Extranet_Rule01"
  action           = "ALLOW"
  custom_msg       = "Test"
  operator         = "AND"
  extranet_enabled = true

  extranet_dto {
    zpn_er_id = data.zpa_extranet_resource_partner.this.id

    location_dto {
      id = data.zpa_location_controller.this.id
    }

    location_group_dto {
      id = data.zpa_location_group_controller.this.id
    }
  }
}
```

## Example Usage - Configuration Location Rule

```hcl
data "zpa_location_controller_summary" "this" {
  name = "BD_CC01_US | NONE | zscalerbeta.net"
}

resource "zpa_policy_access_rule_v2" "this" {
  name        = "ExampleLocationRule"
  description = "ExampleLocationRule"
  action      = "ALLOW"

  conditions {
    operator = "OR"
    operands {
      object_type = "LOCATION"
      values      = [data.zpa_location_controller_summary.this.id]
    }
  }
}
```

## Example Usage - Chrome Enterprise and Chrome Posture Profile

```hcl
data "zpa_managed_browser_profile" "this" {
  name = "Profile01"
}


resource "zpa_policy_access_rule_v2" "this" {
  name        = "Example_v2_100_test"
  description = "Example_v2_100_test"
  action      = "ALLOW"
  custom_msg  = "Test"
  operator    = "AND"

  conditions {
    operator = "OR"
    operands {
      object_type = "CHROME_ENTERPRISE"
      entry_values {
        lhs = "managed"
        rhs = "true"
      }
    }
    operands {
      object_type = "CHROME_POSTURE_PROFILE"
      values      = [data.zpa_managed_browser_profile.this.id]
    }
  }
}
```

## Example Usage - Multiple POSTURE Groups (Separate Condition Blocks)

The API does **not** allow the same `object_type` to appear more than once within a
single `conditions` block (it returns `duplicate.operand.found`). To evaluate more
than one group of the same criteria type — for example two independent groups of
posture profiles, each with its own operator — declare a **separate `conditions`
block per group**. Every block is joined to the others with `AND` (see
[Operator and Condition Behavior](#operator-and-condition-behavior)), while the
operator *inside* each block (`AND`/`OR`) controls how that group's profiles combine.

```hcl
resource "zpa_policy_access_rule_v2" "this" {
  name        = "Example_Posture_Groups"
  description = "Example_Posture_Groups"
  action      = "ALLOW"

  # Group 1: both posture profiles must be VERIFIED (operator = AND)
  conditions {
    operator = "AND"
    operands {
      object_type = "POSTURE"
      entry_values {
        lhs = "cfab2ee9-9bf4-4482-9dcc-dadf7311c49b" # posture_udid
        rhs = "true"
      }
      entry_values {
        lhs = "c6e095ef-9fba-49b9-98b3-f6a7add2faff" # posture_udid
        rhs = "true"
      }
    }
  }

  # Group 2: either posture profile may be VERIFIED (operator = OR)
  conditions {
    operator = "OR"
    operands {
      object_type = "POSTURE"
      entry_values {
        lhs = "661f8a33-2b4d-440f-be66-be4fad12e285" # posture_udid
        rhs = "true"
      }
      entry_values {
        lhs = "3c3aae3c-cd1d-452d-a58c-7e8ccc429e8e" # posture_udid
        rhs = "true"
      }
    }
  }
}
```

## Operator and Condition Behavior

A policy rule evaluates its criteria at three levels. Understanding these levels
explains which operators are allowed where, and why some structures that Terraform
and the API accept cannot be edited in the ZPA Admin GUI.

1. **Rule level — the top-level `operator` argument.** This joins the separate
   `conditions` blocks to each other. **Only `AND` is supported at this level**;
   sending `OR` returns the API error `operator.or.not.supported.policy.rule`
   (`"Operator:OR is not supported at policy rule level."`). The argument is
   optional and computed — if you omit it, ZPA sets it to `AND` automatically, so
   there is no need to set it explicitly. This is why the GUI shows a fixed `AND`
   between criteria groups with no toggle.

2. **Block level — the `operator` inside each `conditions` block.** This joins the
   `operands` *within that one block* and may be `AND` or `OR`. It is managed per
   block, so different blocks can use different operators.

3. **Operand level — multiple `entry_values` (or `values`) inside one operand.**
   Entries within a single operand are **always OR'd** together (any matching value
   satisfies the operand), regardless of the block-level operator.

### Rules and constraints

- **No duplicate `object_type` inside one block.** The same `object_type` may not
  appear in more than one `operands` block within a single `conditions` block; the
  API rejects it with `duplicate.operand.found`. Aggregate all `entry_values`/`values`
  under one operand, or, when you need two independent groups of the same type with
  different operators, split them into separate `conditions` blocks (see
  [Multiple POSTURE Groups](#example-usage---multiple-posture-groups-separate-condition-blocks)).
- **Different `object_type` values may share a block.** A single `conditions` block
  can contain operands of different object types (e.g. `SAML` + `SCIM_GROUP`, or
  `CHROME_ENTERPRISE` + `CHROME_POSTURE_PROFILE`) combined by the block-level operator.
- **GUI compatibility.** Rules that combine *different* object types within one block
  using the `OR` operator are valid via Terraform and the API, but the ZPA Admin GUI's
  criteria builder cannot always represent them. If such a rule is opened and re-saved
  in the GUI, the GUI may silently drop the operands it cannot render. When you manage
  a rule with Terraform, continue managing it with Terraform and avoid re-saving it in
  the GUI.

## Schema

### Required

- `name` (String) This is the name of the policy rule.

### Optional

- `description` (String) This is the description of the access policy rule.
- `action` (String) This is for providing the rule action. Supported values: ``ALLOW``, ``DENY``, and ``REQUIRE_APPROVAL``
- `custom_msg` (String) This is for providing a customer message for the user.
- `device_posture_failure_notification_enabled` (boolean) Enable Device Posture notification on failure.
- `operator` (String) The rule-level operator that joins the separate `conditions` blocks. **Only `AND` is supported**; `OR` is rejected by the API (`operator.or.not.supported.policy.rule`). This argument is optional and computed — if omitted, it defaults to `AND`. See [Operator and Condition Behavior](#operator-and-condition-behavior).
- `extranet_enabled` (boolean) Indiciates if the application is designated for Extranet Application Support (true) or not (false). Extranet applications connect to a partner site or offshore development center that is not directly available on your organization’s network.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

- `app_connector_groups` (Block Set)
  - `id` (String) The ID of an app connector group resource

- `app_server_groups` (Block Set)
  - `id` (String) The ID of a server group resource

- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set) - This is for providing the set of conditions for the policy. Each block is joined to the others with the rule-level `AND` operator, and the same `object_type` cannot be repeated within a single block. See [Operator and Condition Behavior](#operator-and-condition-behavior) for the full operator semantics and constraints.
    - `operator` (String) - The block-level operator that joins the `operands` within this block. Supported values are: `AND` or `OR`.
    - `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
        - `object_type` (String) The object type of the operand. Supported values: `APP`, `APP_GROUP`, `BRANCH_CONNECTOR_GROUP`, `CLIENT_TYPE`, `EDGE_CONNECTOR_GROUP`, `MACHINE_GRP`, `LOCATION`.
        - `values` (Block List) The list of values for the specified object type (e.g., application segment ID and/or segment group ID.).

- `conditions` (Block Set)  - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `PLATFORM`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  Supported values: `android`, `ios`, `linux`, `mac`, `windows`
            - `rhs` - (String) - Supported values: `"true"` or `"false"`

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `COUNTRY_CODE`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  2 Letter Country in ``ISO 3166 Alpha2 Code`` [Lear More](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)
            - `rhs` - (String) - Supported values: `"true"` or `"false"`

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (String) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `POSTURE`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  The Posture Profile `posture_udid` value. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile)
            - `rhs` - (String) - Supported values: `"true"` or `"false"`

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `TRUSTED_NETWORK`
        - `entry_values` (Block Set)
            - `lhs` (String) -  The Trusted Network `network_id` value. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_trusted_network)
            - `rhs` (String) - Supported values: `"true"` or `"false"`

- `conditions` - (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `SCIM_GROUP`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  ID of the Identity Provider
            - `rhs` - (String) - ID of the SCIM Group. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups)

    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `SCIM`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  The SCIM Attribute Header ID. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header)
            - `rhs` - (String) - 	The SCIM Attribute value to match

    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `SAML`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  The ID of the SAML Attribute value. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute)
            - `rhs` - (String) - The SAML attribute string i.e Group name, Department Name, Email address etc.

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `CHROME_ENTERPRISE`, `CHROME_POSTURE_PROFILE`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  Must be set to `managed`
            - `rhs` - (String) - Supported values: `"true"` or `"false"`
        - `values` (Block List) The list of ID values for each `CHROME_POSTURE_PROFILE`

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (String) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `CHROME_ENTERPRISE`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  `"managed"`
            - `rhs` - (String) - Supported values: `"true"` or `"false"`

    - `operands` (String) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `CHROME_POSTURE_PROFILE`
        - `values` (Block List) The list of values for the specified object type (e.g., managed browser profile ID `zpa_managed_browser_profile`).

- `extranet_dto` (Block Set) - Extranet location and location group configuration
    - `zpn_er_id` (String) - The unique identifier of the extranet resource that is configured in ZIA. Use the data source `zpa_extranet_resource_partner` to retrieve the Extranet ID
        - `location_dto` (Block Set)
            - `id` - (String) -  Unique identifiers for the location
        - `location_group_dto` (Block Set)
            - `id` - (String) -  Unique identifiers for the location group

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy access rule can be imported by using `<RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule_v2.example <rule_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS| VALUES
|----------|-----------|----------|----------
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment)  | `NA` | `NA` | ``application_segment_id`` |
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group)  | `NA`  | `NA` | ``segment_group_id``|
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types)  | `NA`  | `NA` |  ``zpn_client_type_exporter``, ``zpn_client_type_exporter_noauth``, ``zpn_client_type_machine_tunnel``, ``zpn_client_type_edge_connector``, ``zpn_client_type_zia_inspection``, ``zpn_client_type_vdi``, ``zpn_client_type_zapp``, ``zpn_client_type_slogger``, ``zpn_client_type_zapp_partner``, ``zpn_client_type_browser_isolation``, ``zpn_client_type_ip_anchoring``, ``zpn_client_type_branch_connector`` |
| [EDGE_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group)  | `NA`  | `NA` |  ``<edge_connector_id>`` |
| [BRANCH_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group)  | `NA` | `NA` |  ``<branch_connector_id>`` |
| [LOCATION](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group)   | `NA` | `NA` | ``location_id`` |
| [MACHINE_GRP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group)   | `NA` | `NA` | ``machine_group_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
| [PLATFORM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``mac``, ``ios``, ``windows``, ``android``, ``linux`` | ``"true"`` |
| [POSTURE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile) | ``posture_udid``  | ``"true"`` / ``"false"`` |
| [TRUSTED_NETWORK](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_trusted_network) | ``network_id``  | ``"true"`` |
| [COUNTRY_CODE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_platforms) | [2 Letter ISO3166 Alpha2](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)  | ``"true"`` |
| [RISK_FACTOR_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``ZIA``  | ``"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"`` |
| [CHROME_ENTERPRISE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``managed``  | ``"true" / "false"`` |