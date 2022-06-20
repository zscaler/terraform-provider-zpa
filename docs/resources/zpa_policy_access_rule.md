---
subcategory: "Policy Set Controller"
layout: "zscaler"
page_title: "ZPA: policy_access_rule"
description: |-
  Creates and manages ZPA Policy Access Rule.
---

# Resource: zpa_policy_access_rule

The **zpa_policy_access_rule** resource creates and manages policy access rule in the Zscaler Private Access cloud.

## Example Usage

```hcl
#Create Policy Access Rule
resource "zpa_policy_access_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "ALLOW"
  operator                      = "AND"
  policy_set_id                 = data.zpa_policy_type.access_policy.id

  conditions {
    negated   = false
    operator  = "OR"
    operands {
      object_type = "APP"
      lhs = "id"
      rhs = zpa_application_segment.test_app_segment.id
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

# Get Global Access Policy ID
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

# Get IdP ID
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

# Get SCIM Group attribute ID
data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "IdP_Name"
}
```

### Required

* `name` - (Required) This is the name of the policy rule.
* `policy_set_id` - (Required)

## Attributes Reference

* `action` (Optional) This is for providing the rule action.
* `action_id` (Optional) This field defines the description of the server.
* `custom_msg` (Optional) This is for providing a customer message for the user.
* `description` (Optional) This is the description of the access policy rule.
* `operator` (Optional)
* `policy_type` (Optional)
* `rule_order` (Optional)

* `conditions` - (Optional)
  * `negated` - (Optional)
  * `operator` (Optional)
  * `operands`
    * `name` (Optional)
    * `lhs` (Optional)
    * `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.
    * `idp_id` (Optional)
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `CLOUD_CONNECTOR_GROUP`. `TRUSTED_NETWORK`, and `CLIENT_TYPE`.
    * `CLIENT_TYPE` (Optional) - The below options are the only ones supported in an access policy rule.
      * `zpn_client_type_exporter`
      * `zpn_client_type_browser_isolation`
      * `zpn_client_type_machine_tunnel`
      * `zpn_client_type_ip_anchoring`
      * `zpn_client_type_edge_connector`
      * `zpn_client_type_zapp`

* `app_connector_groups`
  * `id` - (Optional) The ID of an app connector group resource

* `app_server_groups`
  * `id` - (Optional) The ID of a server group resource

## Import

Policy access rule can be imported by using `<POLICY ACCESS RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example <policy_access_rule_id>
```
