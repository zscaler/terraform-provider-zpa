---
page_title: "zpa_policy_portal_access_rule Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller V2"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-capabilities-policy
  API documentation https://help.zscaler.com/zpa/configuring-privileged-policies-using-api
  Creates and manages ZPA Privileged Portal Policy Capabilities Rule.
---
# zpa_policy_portal_access_rule (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-capabilities-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-policies-using-api)

The **zpa_policy_portal_access_rule** resource creates a Privileged Portal Policy Capabilities Rule in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
data "zpa_saml_attribute" "email_user_sso" {
  name     = "Email_BD_Okta_Users"
  idp_name = "BD_Okta_Users"
}

data "zpa_idp_controller" "this" {
  name = "BD_Okta_Users"
}

data "zpa_scim_groups" "a000" {
  idp_name = "BD_Okta_Users"
  name     = "A000"
}

data "zpa_scim_groups" "b000" {
  idp_name = "BD_Okta_Users"
  name     = "B000"
}

resource "zpa_policy_portal_access_rule" "example" {
  name        = "Portal Access Rule"
  description = "Allow portal access with specific capabilities"
  action      = "CHECK_PRIVILEGED_PORTAL_CAPABILITIES"

  privileged_portal_capabilities {
    delete_file             = true
    access_uninspected_file = true
    request_approvals       = true
    review_approvals        = true
  }

  conditions {
    operator = "OR"
    operands {
      object_type = "PRIVILEGE_PORTAL"
      values      = ["216196257331387235"]
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "COUNTRY_CODE"
      entry_values {
        lhs = "BR"
        rhs = "true"
      }
      entry_values {
        lhs = "CA"
        rhs = "true"
      }
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
    - `delete_file` - (Boolean) Allows a User to delete files to reclaim space. Allowing deletion will prevent auditing of the file. Supported values: `true` or `false`
    - `access_uninspected_file` - (Boolean) Allows a User like an Admin to see all files marked Uninspected from other users in the tenant. Supported values: `true` or `false`
    - `request_approvals` - (Boolean) Indicates the request approvals is enabled or disabled. Supported values: `true` or `false`
    - `review_approvals` - (Boolean) Indicates the review approvals is enabled or disabled. Supported values: `true` or `false`

### Optional

- `description` (String) This is the description of the access policy rule.
- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set)  - This is for providing the set of conditions for the policy. Separate condition blocks for each object type is required.
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Optional) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
        - `object_type` (String) The object type of the operand. Supported values: `PRIVILEGE_PORTAL`.
        - `values` (Block List) The list of values for the specified object type (e.g., User Portal IDs).
            **NOTE** Use the resource or data source `zpa_user_portal_controller` to retrieve the User Portal ID information.

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
            - `rhs` - (String) - 	The SCIM Attribute value to match

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy criteria. Supported values: `COUNTRY_CODE`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  2 Letter Country in ``ISO 3166 Alpha2 Code`` [Lear More](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)
            - `rhs` - (String) - Supported values: `"true"` or `"false"`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

Policy access capability can be imported by using `<RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_portal_access_rule.example <rule_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS| VALUES
|----------|-----------|----------|----------
| [PRIVILEGE_PORTAL](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_user_portal_controller)  | `NA` | `NA` | ``user_portal_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
| [COUNTRY_CODE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_platforms) | [2 Letter ISO3166 Alpha2](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)  | ``"true"`` |
