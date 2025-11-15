---
page_title: "zpa_policy_capabilities_rule Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller V2"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-capabilities-policy
  API documentation https://help.zscaler.com/zpa/configuring-privileged-policies-using-api
  Creates and manages ZPA Policy Capabilities Rule.
---

# zpa_policy_capabilities_rule (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-capabilities-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-policies-using-api)

The **zpa_policy_capabilities_rule** resource creates a policy capabilities rule in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
data "zpa_idp_controller" "this" {
	name = "IdP_Users"
}

data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_IdP_Users"
    idp_name = "IdP_Users"
}

data "zpa_saml_attribute" "group_user" {
    name = "GroupName_IdP_Users"
    idp_name = "IdP_Users"
}

data "zpa_scim_groups" "a000" {
    name = "A000"
    idp_name = "IdP_Users"
}

data "zpa_scim_groups" "b000" {
    name = "B000"
    idp_name = "IdP_Users"
}

resource "zpa_policy_capabilities_rule" "this" {
    name = "Example"
    description = "Example"
    action = "CHECK_CAPABILITIES"
    privileged_capabilities {
        file_upload = true
        file_download = true
        inspect_file_upload = true
        clipboard_copy = true
        clipboard_paste = true
        record_session = true
    }
    conditions {
        operator = "OR"
            operands {
                object_type = "SAML"
                entry_values {
                    rhs = "user1@example.com"
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

- `name` (String) This is the name of the policy rule.
- `action` (Required) This is for providing the rule action. Supported value: ``CHECK_CAPABILITIES``

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)
- `privileged_capabilities` (Required) The Privileged Remote Access application segment resource
    - `clipboard_copy` - (Boolean) Indicates the PRA Clipboard Copy function. Supported values: `true` or `false`
    - `clipboard_paste` - (Boolean) Indicates the PRA Clipboard Paste function. Supported values: `true` or `false`
    - `file_download` - (Boolean) Indicates the PRA File Transfer capabilities that enables the File Download function. Supported values: `true` or `false`
    - `file_upload` - (Boolean) Indicates the PRA File Transfer capabilities that enables the File Upload function. Supported values: `true` or `false`
    - `inspect_file_download` - (Boolean) Inspects the file via ZIA sandbox (if you have set up the ZIA cloud and the Integrations settings) and downloads the file following the inspection. Supported values: `true` or `false`
    - `inspect_file_upload` - (Boolean) Inspects the file via ZIA sandbox (if you have set up the ZIA cloud and the Integrations settings) and uploads the file following the inspection. Supported values: `true` or `false`
    - `monitor_session` - (Boolean) Indicates the PRA Monitoring Capabilities to enable the PRA Session Monitoring function. Supported values: `true` or `false`
    - `record_session` - (Boolean) Indicates the PRA Session Recording capabilities to enable PRA Session Recording. Supported values: `true` or `false`
    - `share_session` - (Boolean) Indicates the PRA Session Control and Monitoring capabilities to enable PRA Session Monitoring. Supported values: `true` or `false`

### Optional

- `description` (String) This is the description of the access policy rule.
- `rule_order` (String, Deprecated)
- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
        - `object_type` (String) The object type of the operand. Supported values: `APP` for application segments and `APP_GROUP` for segment groups
        - `values` (List of Strings) The list of values for the specified object type (e.g., application segment ID and/or segment group ID.).

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `SAML`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  The ID of the SAML Attribute value. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute)
            - `rhs` - (String) - The SAML attribute string i.e Group name, Department Name, Email address etc.
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy critiera. Supported values: `SCIM_GROUP`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  The Identity Provider (IdP) ID
            - `rhs` - (String) - The SCIM Group unique identified (ID)
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy critiera. Supported values: `SCIM`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  The SCIM Attribute Header ID. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header)
            - `rhs` - (String) - 	The SCIM Attribute value to match

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

Policy access capability can be imported by using `<RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_capabilities_rule.example <rule_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS| VALUES
|----------|-----------|----------|----------
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment) |   |  | ``application_segment_id``
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group) |   |  | ``segment_group_id``
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
