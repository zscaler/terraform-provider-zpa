---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_access_rule_v2"
description: |-
  Creates and manages ZPA Policy Access Rule via API v2 endpoints.
---

# Resource: zpa_policy_access_rule_v2

The **zpa_policy_access_rule_v2** resource creates and manages policy access rule in the Zscaler Private Access cloud using a new v2 API endpoint.

  ⚠️ **NOTE**: This resource is recommended if your configuration requires the association of more than 1000 resource criteria per rule.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```hcl
# Retrieve Policy Types
data "zpa_policy_type" "this" {
  policy_type = "ACCESS_POLICY"
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
}
```

### Required

* `name` - (Required) This is the name of the policy rule.

## Attributes Reference

* `description` (Optional) This is the description of the access policy rule.
* `action` (Optional) This is for providing the rule action. Supported values: ``ALLOW``, ``DENY``, and ``REQUIRE_APPROVAL``
* `custom_msg` (Optional) This is for providing a customer message for the user.
* `rule_order` - (Deprecated)

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

* `app_connector_groups`
  * `id` - (Optional) The ID of an app connector group resource

* `app_server_groups`
  * `id` - (Optional) The ID of a server group resource

* `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

* `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
        * `object_type` (Optional) The object type of the operand. Supported values: `APP`, `APP_GROUP`, `BRANCH_CONNECTOR_GROUP`, `CLIENT_TYPE`, `EDGE_CONNECTOR_GROUP`, `MACHINE_GRP`, `LOCATION`.
        * `values` (Optional) The list of values for the specified object type (e.g., application segment ID and/or segment group ID.).

* `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        * `object_type` (Optional) This is for specifying the policy criteria. Supported values: `PLATFORM`
        * `entry_values` (Optional)
            * `lhs` - (Optional) -  Supported values: `android`, `ios`, `linux`, `mac`, `windows`
            * `rhs` - (Optional) - Supported values: `"true"` or `"false"`

* `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        * `object_type` (Optional) This is for specifying the policy criteria. Supported values: `COUNTRY_CODE`
        * `entry_values` (Optional)
            * `lhs` - (Optional) -  2 Letter Country in ``ISO 3166 Alpha2 Code`` [Lear More](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)
            * `rhs` - (Optional) - Supported values: `"true"` or `"false"`

* `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        * `object_type` (Optional) This is for specifying the policy criteria. Supported values: `POSTURE`
        * `entry_values` (Optional)
            * `lhs` - (Optional) -  The Posture Profile `posture_udid` value.
            * `rhs` - (Optional) - Supported values: `"true"` or `"false"`

* `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        * `object_type` (Optional) This is for specifying the policy criteria. Supported values: `TRUSTED_NETWORK`
        * `entry_values` (Optional)
            * `lhs` - (Optional) -  The Trusted Network `network_id` value.
            * `rhs` - (Optional) - Supported values: `"true"` or `"false"`

* `conditions` - (Optional) - This is for providing the set of conditions for the policy
    * `operator` (Optional) - Supported values are: `AND` or `OR`
    * `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        * `object_type` (Optional) This is for specifying the policy criteria. Supported values: `SAML`, `SCIM`, `SCIM_GROUP`
        * `entry_values` (Optional)
            * `lhs` - (Optional) -  2 Letter Country in ``ISO 3166 Alpha2 Code`` [Lear More](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)
            * `rhs` - (Optional) - Supported values: `"true"` or `"false"`

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
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment)  |   |  | ``application_segment_id`` |
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group)  |   |  | ``segment_group_id``|
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types)  |   |  |  ``zpn_client_type_zappl``, ``zpn_client_type_exporter``, ``zpn_client_type_browser_isolation``, ``zpn_client_type_ip_anchoring``, ``zpn_client_type_edge_connector``, ``zpn_client_type_branch_connector``,  ``zpn_client_type_zapp_partner``, ``zpn_client_type_zapp``  |
| [EDGE_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group)  |   |  |  ``<edge_connector_id>`` |
| [BRANCH_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group)  |   |  |  ``<branch_connector_id>`` |
| [LOCATION](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group)   |   |  | ``location_id`` |
| [MACHINE_GRP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group)   |   |  | ``machine_group_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
| [PLATFORM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``mac``, ``ios``, ``windows``, ``android``, ``linux`` | ``"true"`` / ``"false"`` |
| [POSTURE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile) | ``posture_udid``  | ``"true"`` / ``"false"`` |
| [TRUSTED_NETWORK](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_trusted_network) | ``network_id``  | ``"true"`` |
| [COUNTRY_CODE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_platforms) | [2 Letter ISO3166 Alpha2](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)  | ``"true"`` / ``"false"`` |
