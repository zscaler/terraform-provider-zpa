---
page_title: "zpa_policy_credential_rule Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller V2"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-capabilities-policy
  API documentation https://help.zscaler.com/zpa/configuring-privileged-policies-using-api
  Creates and manages ZPA Policy Credential Access Rule.
---

# zpa_policy_credential_rule (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-capabilities-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-privileged-policies-using-api)

The **zpa_policy_credential_rule** resource creates a policy credential rule in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
resource "zpa_application_segment_pra" "this" {
  name             = "Example"
  description      = "Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["22", "22", "3389", "3389"]
  domain_names     = ["ssh_pra.example.com", "rdp_pra.example.com"]
  segment_group_id = zpa_segment_group.this.id
  common_apps_dto {
    apps_config {
      name                 = "rdp_pra"
      domain               = "rdp_pra.example.com"
      application_protocol = "RDP"
      connection_security  = "ANY"
      application_port     = "3389"
      enabled              = true
      app_types            = ["SECURE_REMOTE_ACCESS"]
    }
    apps_config {
      name                 = "ssh_pra"
      domain               = "ssh_pra.example.com"
      application_protocol = "SSH"
      application_port     = "22"
      enabled              = true
      app_types            = ["SECURE_REMOTE_ACCESS"]
    }
  }
}

resource "zpa_segment_group" "this" {
  name        = "Example"
  description = "Example"
  enabled     = true
}

data "zpa_ba_certificate" "this1" {
  name = "pra01.example.com"
}

resource "zpa_pra_portal_controller" "this" {
  name                      = "pra01.example.com"
  description               = "pra01.example.com"
  enabled                   = true
  domain                    = "pra01.example.com"
  certificate_id            = data.zpa_ba_certificate.this.id
  user_notification         = "Created with Terraform"
  user_notification_enabled = true
}

locals {
  pra_application_ids = {
    for app_dto in flatten([for common_apps in zpa_application_segment_pra.this.common_apps_dto : common_apps.apps_config]) :
    app_dto.name => app_dto.id
  }
  pra_application_id_ssh_pra = lookup(local.pra_application_ids, "ssh_pra", "")
  pra_application_id_rdp_pra = lookup(local.pra_application_ids, "rdp_pra", "")
}

resource "zpa_pra_console_controller" "ssh_pra" {
  name        = "ssh_console"
  description = "Created with Terraform"
  enabled     = true
  pra_application {
    id = local.pra_application_id_ssh_pra
  }
  pra_portals {
    id = [zpa_pra_portal_controller.this1.id]
  }
  depends_on = [ zpa_application_segment_pra.this ]
}

resource "zpa_pra_console_controller" "rdp_pra" {
  name        = "rdp_console"
  description = "Created with Terraform"
  enabled     = true
  pra_application {
    id = local.pra_application_id_rdp_pra
  }
  pra_portals {
    id = [zpa_pra_portal_controller.this1.id]
  }
  depends_on = [ zpa_application_segment_pra.this ]
}

resource "zpa_pra_credential_controller" "this" {
    name = "John Carrow"
    description = "Created with Terraform"
    credential_type = "USERNAME_PASSWORD"
    user_domain = "acme.com"
    username = "jcarrow"
    password = ""
}

data "zpa_idp_controller" "this" {
	name = "Idp_Users"
}

data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_Idp_Users"
    idp_name = "Idp_Users"
}

data "zpa_saml_attribute" "group_user" {
    name = "GroupName_Idp_Users"
    idp_name = "Idp_Users"
}

data "zpa_scim_groups" "a000" {
    name = "A000"
    idp_name = "Idp_Users"
}

data "zpa_scim_groups" "b000" {
    name = "B000"
    idp_name = "Idp_Users"
}

resource "zpa_policy_credential_rule" "this" {
  name          = "Example_Credential_Rule"
  description   = "Example_Credential_Rule"
  action        = "INJECT_CREDENTIALS"

  credential {
    id = zpa_pra_credential_controller.this.id
  }

  conditions {
    operator = "OR"
    operands {
      object_type = "CONSOLE"
      values      = [ zpa_pra_console_controller.rdp_pra.id, zpa_pra_console_controller.ssh_pra.id ]
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
- `action` (String) This is for providing the rule action. Supported value: ``INJECT_CREDENTIALS``
- `credential` - The Privileged Remote Access application segment resource
    - `id` - (String) The unique identifier of the privileged credential.
    
### Optional

- `description` (String) This is the description of the access policy rule.
- `rule_order` (String, Deprecated)

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)
- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set) - This is for providing the set of conditions for the policy
    - `operator` (String) - Supported values are: `AND` or `OR`
    - `operands` (Block Set) - This signifies the various policy criteria. Supported Values: `object_type`, `values`
        - `object_type` (String) The object type of the operand. Supported values: `CONSOLE`
        - `values` (List of Strings) The list of values for the specified object type (e.g., PRA Console IDs).

  ⚠️ **WARNING:**: The first condition block specifying the `object_type` / `CONSOLE` is mandatory. This block refers to the `zpa_pra_console_controller` resource.

- `conditions` (Block Set)  - This is for providing the set of conditions for the policy
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
    - `operands` (Block Set)  - This signifies the various policy criteria. Supported Values: `object_type`, `entry_values`
        - `object_type` (String) This is for specifying the policy critiera. Supported values: `SCIM`
        - `entry_values` (Block Set)
            - `lhs` - (String) -  The SCIM Attribute Header ID. [See Documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header)
            - `rhs` - (String) - 	The SCIM Attribute value to match

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

Policy access credential can be imported by using `<POLICY CREDENTIAL ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_credential_rule.example <policy_credential_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS| VALUES
|----------|-----------|----------|----------
| [CONSOLE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_pra_console_controller) |   |  | ``pra_console_id``
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
