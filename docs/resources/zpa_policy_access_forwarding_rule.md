---
subcategory: "Policy Access Forwarding Rule"
layout: "zpa"
page_title: "ZPA: policy_forwarding_rule"
description: |-
  Creates and manages ZPA Policy Access Forwarding Rule.
---
# zpa_policy_forwarding_rule (Resource)

The **zpa_policy_forwarding_rule** resource creates a policy forwarding access rule in the Zscaler Private Access cloud.

## Example Usage

```hcl
#Create Client Forwarding Access Rule
resource "zpa_policy_forwarding_rule" "test_forwarding_rule" {
  name            = "test1-forwarding-rule"
  description     = "test1-forwarding-rule"
  action          = "BYPASS"
  operator        = "AND"
  policy_set_id   = data.zpa_policy_type.client_forwarding_policy.id

  conditions {
    negated   = false
    operator  = "OR"
    operands {
      object_type = "APP"
      lhs = "id"
      rhs = [ zpa_application_segment.test_app_segment.id ]
    }
  }
  conditions {
     negated = false
     operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = [ data.zpa_scim_groups.engineering.id ]
    }
  }
  depends_on = [
    data.zpa_policy_type.client_forwarding_policy,
    data.zpa_idp_controller.idp_name,
    data.zpa_scim_groups.engineering,
    zpa_application_segment.test_app_segment
  ]
}

# Get Global Policy Forwading ID
data "zpa_policy_type" "client_forwarding_policy" {
    policy_type = "CLIENT_FORWARDING_POLICY"
}

# Get IdP ID
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

# Get SCIM Group attribute ID
data "zpa_scim_groups" "engineering" {
  name = "Engineering"
  idp_name = "IdP_Name"
}

```hcl
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
}
```

### Required

* `name` - (Required) This is the name of the forwarding policy rule.
* `policy_set_id` - (Required)

## Attributes Reference

* `action` - (Optional) This is for providing the rule action.
  * The supported actions for a policy forwarding rule are: `BYPASS`, `INTERCEPT` or `INTERCEPT_ACCESSIBLE`
* `action_id` - (Optional) This field defines the description of the server.
* `bypass_default_rule` - (Optional)
* `custom_msg` - (Optional) This is for providing a customer message for the user.
* `description` - (Optional) This is the description of the access policy rule.
* `operator` (Optional)
* `policy_type` - (Optional)
  * The supported policy type values for a policy forwarding rule are: `CLIENT_FORWARDING_POLICY` and `BYPASS_POLICY`
* `rule_order` - (Optional)

* `conditions` - (Optional)
  * `negated` - (Optional)
  * `operator` (Optional)
  * `operands`
    * `name` (Optional)
    * `lhs` (Optional)
    * `rhs` (Optional) This denotes the value for the given object type. Its value depends upon the key.
    * `idp_id` (Optional)
    * `object_type` (Optional) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `CLOUD_CONNECTOR_GROUP`. `TRUSTED_NETWORK`, and `CLIENT_TYPE`.
    * `CLIENT_TYPE` (Optional) - The below options are the only ones supported in a timeout policy rule.
      * `zpn_client_type_exporter`
      * `zpn_client_type_browser_isolation`
      * `zpn_client_type_machine_tunnel`
      * `zpn_client_type_ip_anchoring`
      * `zpn_client_type_edge_connector`
      * `zpn_client_type_zapp`

## Import

Policy Access Forwarding Rule can be imported by using `<POLICY FORWARDING RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_forwarding_rule.example <policy_forwarding_rule_id>
```
