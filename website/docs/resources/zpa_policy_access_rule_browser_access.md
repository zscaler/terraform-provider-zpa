---
subcategory: "Policy Access Rule Browser Access"
page_title: "ZPA: policyset_rule"
description: |-
  Creates a ZPA Policy Access Rule for Browser Access.
---

# zpa_policyset_rule (Resource)

The **zpa_policyset_rule** resource creates a policy access rule in the Zscaler Private Access cloud.

## Example Usage

```hcl
resource "zpa_policy_access_rule" "browser_access_rule" {
  name                          = "Browser Access Corporate Services"
  description                   = "Browser Access Corporate Services"
  action                        = "ALLOW"
  operator                      = "AND"
  policy_set_id                 = data.zpa_policy_type.access_policy.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "APP"
      lhs = "id"
      rhs = [zpa_application_segment.as_corporate_services.id]
    }
  }

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
    negated = false
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [data.zpa_scim_groups.engineering.id]
    }
  }
}

// Get Global Access Policy ID
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

// Get IdP ID
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

// Get SCIM Group attribute ID
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP_Name"
}
```

```hcl
// Create Application Segment
resource "zpa_application_segment" "as_corporate_services" {
    name = "Corporate Services"
    description = "Corporate Services"
    enabled = true
    health_reporting = "ON_ACCESS"
    bypass_type = "NEVER"
    tcp_port_ranges = ["1", "52", "54", "65535"]
    domain_names = ["*.acme.com"]
    segment_group_id = zpa_segment_group.sg_corporate_services.id
    server_groups {
        id = [zpa_server_group.corporate_server_group.id]
    }
}
```

```hcl
// Create Segment Group
 resource "zpa_segment_group" "sg_corporate_services" {
   name = "Corporate Services"
   description = "Corporate Services"
   enabled = true
   policy_migrated = true
 }
```

```hcl
// Create Server Group
resource "zpa_server_group" corporate_server_group" {
  name =  "Corporate Services Group"
  description = "Corporate Services Group"
  enabled = true
  dynamic_discovery = true
  app_connector_groups {
    id = [data.zpa_app_connector_group.aws_connector.id]
  }
}

// Get App Connector Group ID
data "zpa_app_connector_group" "aws_connector" {
  name = "AWS-Connector"
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required)

## Attributes Reference

* `action` (String) This is for providing the rule action.
* `action_id` (String) This field defines the description of the server.
* `bypass_default_rule` (Boolean)
* `custom_msg` (String) This is for providing a customer message for the user.
* `description` (String) This is the description of the access policy rule.
* `operator` (String)
* `policy_type` (String)
* `priority` (String)
* `reauth_default_rule` (Boolean)
* `reauth_idle_timeout` (String)
* `reauth_timeout` (String)
* `rule_order` (String)

`conditions` - (Optional)

* `negated` (Optional)
* `idp_id` (Optional)
* `operator` (String)
* `name` (Optional)
* `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `CLOUD_CONNECTOR_GROUP`. TRUSTED_NETWORK is only supported for CLIENT_TYPE

`operands`

* `lhs` (Optional)
* `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.

`app_connector_groups`

* `id` - (Optional) The ID of this resource.

`app_server_groups`

* `id` - (Optional) The ID of this resource.

## Import

Application Segment can be imported; use `<POLICY Access RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example 216196257331290863
```
