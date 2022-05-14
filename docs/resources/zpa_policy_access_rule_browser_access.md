---
subcategory: "Policy Access Rule Browser Access"
page_title: "ZPA: policyset_rule"
description: |-
  Creates and manages ZPA Policy Access Rule for Browser Access.
---

# zpa_policyset_rule (Resource)

The **zpa_policyset_rule** resource creates and manages policy access rule to support Browser Access  in the Zscaler Private Access cloud.

## Example Usage

```hcl
resource "zpa_policy_access_rule" "browser_access_rule" {
  name                          = "test1-ba-policy-access"
  description                   = "test1-ba-policy-access"
  action                        = "ALLOW"
  operator                      = "AND"
  policy_set_id                 = data.zpa_policy_type.access_policy.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      object_type = "APP"
      lhs = "id"
      rhs = [zpa_application_segment.test_app_segment.id]
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

# ZPA Application Segment resource
resource "zpa_application_segment" "test_app_segment" {
    name              = "test1-app-segment"
    description       = "test1-app-segment"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
    tcp_port_ranges   = ["8080", "8080"]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.test_segment_group.id
    server_groups {
        id = [ zpa_server_group.test_server_group.id]
    }
    depends_on = [ zpa_server_group.test_server_group, zpa_server_group.test_segment_group]
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

Policy Access Rule for Browser Access can be imported by using`<POLICY ACCESS RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example <policy_access_rule_id>
```
