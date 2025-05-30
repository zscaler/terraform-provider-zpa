---
page_title: "zpa_policy_inspection_rule_v2 Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller V2"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-security-policy
  API documentation https://help.zscaler.com/zpa/configuring-appprotection-policies-using-api
  Creates and manages ZPA Policy Access Inspection Rule via API v2 endpoints.
---

# zpa_policy_inspection_rule_v2 (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-security-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-appprotection-policies-using-api)

The **zpa_policy_inspection_rule_v2** resource creates and manages policy access inspection rule in the Zscaler Private Access cloud using a new v2 API endpoint.

  ⚠️ **NOTE**: This resource is recommended if your configuration requires the association of more than 1000 resource criteria per rule.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
data "zpa_inspection_predefined_controls" "this" {
  name    = "Failed to parse request body"
  version = "OWASP_CRS/3.3.0"
}

data "zpa_inspection_all_predefined_controls" "default_predefined_controls" {
  version    = "OWASP_CRS/3.3.0"
  group_name = "preprocessors"
}

resource "zpa_inspection_profile" "this" {
  name              = "Example"
  description       = "Example"
  paranoia_level    = "1"
  dynamic "predefined_controls" {
    for_each = data.zpa_inspection_all_predefined_controls.default_predefined_controls.list
    content {
      id           = predefined_controls.value.id
      action       = predefined_controls.value.action == "" ? predefined_controls.value.default_action : predefined_controls.value.action
      action_value = predefined_controls.value.action_value
    }
  }
  predefined_controls {
    id     = data.zpa_inspection_predefined_controls.this.id
    action = "BLOCK"
  }
}

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

# Create Policy Access Isolation Rule V2
resource "zpa_policy_inspection_rule_v2" "this" {
  name                      = "Example"
  description               = "Example"
  action                    = "INSPECT"
  zpn_inspection_profile_id = zpa_inspection_profile.this.id

  conditions {
    operator = "OR"
    operands {
      object_type = "CLIENT_TYPE"
      values      = ["zpn_client_type_exporter"]
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
}
```

## Schema

### Required

- `name` - (String) This is the name of the policy rule.
- `action` (String) This is for providing the rule action. Supported values: ``INSPECT`` and ``BYPASS_INSPECT``.
- `zpn_inspection_profile_id` (String) An inspection profile is required if the `action` is set to `INSPECT`

### Optional

- `description` (String) This is the description of the access policy rule.
- `rule_order` (String, Deprecated)

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set)  Specifies the set of conditions for the policy rule. Separate condition blocks for each object type is required.
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
        - `object_type` (String) The object type of the operand. Supported values: `APP`, `APP_GROUP`, `CLIENT_TYPE`, `EDGE_CONNECTOR_GROUP`, `MACHINE_GRP`
        - `values` (List of Strings) The list of values for the specified object type (e.g., application segment ID and/or segment group ID)

- `conditions` (Block Set)  Specifies the set of conditions for the policy rule.
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (String) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `PLATFORM`
        - `entry_values` (Block Set) 
            - `lhs` - (String) -  Supported values: `android`, `ios`, `linux`, `mac`, `windows`
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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy access inspection rule can be imported by using `<RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_inspection_rule_v2.example <rule_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS| VALUES
|----------|-----------|----------|----------
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment)  |   |  | ``application_segment_id`` |
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group)  |   |  | ``segment_group_id``|
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types)  |   |  |  ``zpn_client_type_zappl``, ``zpn_client_type_exporter``, ``zpn_client_type_browser_isolation``, ``zpn_client_type_ip_anchoring``, ``zpn_client_type_edge_connector``, ``zpn_client_type_branch_connector``,  ``zpn_client_type_zapp_partner``, ``zpn_client_type_zapp``  |
| [EDGE_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group)  |   |  |  ``<edge_connector_id>`` |
| [MACHINE_GRP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group)   |   |  | ``machine_group_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
| [PLATFORM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``mac``, ``ios``, ``windows``, ``android``, ``linux`` | ``"true"`` / ``"false"`` |
| [POSTURE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile) | ``posture_udid``  | ``"true"`` / ``"false"`` |