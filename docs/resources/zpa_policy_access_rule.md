---
page_title: "zpa_policy_access_rule Resource - terraform-provider-zpa"
subcategory: "Policy Set Controller"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-access-policy
  API documentation https://help.zscaler.com/zpa/configuring-access-policies-using-api
  Creates and manages ZPA Policy Access Rule.
---

# zpa_policy_access_rule (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-access-policy)
* [API documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api)

The **zpa_policy_access_rule** resource creates and manages policy access rule in the Zscaler Private Access cloud.

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

## Example Usage

```terraform
# Get IdP ID
data "zpa_idp_controller" "idp_name" {
 name = "IdP_Name"
}

# Get SCIM Group attribute ID
data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "IdP_Name"
}

#Create Policy Access Rule
resource "zpa_policy_access_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "ALLOW"
  operator                      = "AND"

  conditions {
    operator  = "OR"
    operands {
      object_type = "APP"
      lhs = "id"
      rhs = zpa_application_segment.this.id
    }
  }
  conditions {
    operator = "OR"
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.idp_name.id
      rhs = data.zpa_scim_groups.engineering.id
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

## Schema

### Required

- `name` (String) This is the name of the policy rule.

### Optional

- `policy_set_id` - (String) Use [zpa_policy_type](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_policy_type) data source to retrieve the necessary policy Set ID ``policy_set_id``
    ~> **NOTE** As of v3.2.0 the ``policy_set_id`` attribute is now optional, and will be automatically determined based on the policy type being configured. The attribute is being kept for backwards compatibility, but can be safely removed from existing configurations.
    
- `action` (String) This is for providing the rule action. Supported values: ``ALLOW``, ``DENY``, and ``REQUIRE_APPROVAL``
- `custom_msg` (String) This is for providing a customer message for the user.
- `description` (String) This is the description of the access policy rule.
- `operator` (String) Supported values: ``AND``, ``OR``
- `rule_order` (String, Deprecated)

  ⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of the new resource  [``policy_access_rule_reorder``](zpa_policy_access_rule_reorder.md)

- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `conditions` (Block Set) 
  - `operator` (String) Supported values: ``AND``, and ``OR``
  - `microtenant_id` (Optional) The ID of the microtenant the resource is to be associated with.

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

  - `operands` (Block Set)  - Operands block must be repeated if multiple per `object_type` conditions are to be added to the rule.
    - `lhs` (String) LHS must always carry the string value ``id`` or the attribute ID of the resource being associated with the rule.
    - `rhs` (String) RHS is either the ID attribute of a resource or fixed string value. Refer to the chart below for further details.
    - `idp_id` (String)
    - `object_type` (String) This is for specifying the policy critiera. Supported values: `APP`, `APP_GROUP`, `SAML`, `IDP`, `CLIENT_TYPE`, `TRUSTED_NETWORK`, `POSTURE`, `SCIM`, `SCIM_GROUP`, and `EDGE_CONNECTOR_GROUP`. `TRUSTED_NETWORK`, `CLIENT_TYPE`, `PLATFORM`, `COUNTRY_CODE`.
    - `CLIENT_TYPE` (Optional) - The below options are the only ones supported in an access policy rule.
      - `"zpn_client_type_exporter"` - "Web Browser"`
      - `"zpn_client_type_exporter_noauth"` - "Web Browser Unauthenticated"`
      - `"zpn_client_type_machine_tunnel"` - "Machine Tunnel"
      - `"zpn_client_type_edge_connector"` - "Cloud Connector"
      - `"zpn_client_type_zia_inspection"` - "ZIA Inspection"
      - `"zpn_client_type_vdi"` - "Client Connector for VDI"
      - `"zpn_client_type_zapp"` - "Client Connector"
      - `"zpn_client_type_slogger"` - "ZPA LSS"
      - `"zpn_client_type_browser_isolation"` - "Cloud Browser"
      - `"zpn_client_type_ip_anchoring"` - "ZIA Service Edge"
      - `"zpn_client_type_zapp_partner"` - "Client Connector Partner"
      - `"zpn_client_type_branch_connector"` - "Branch Connector"

  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is not supported within the `operands` block when the `object_type` is set to `SAML`, `SCIM`, `SCIM_GROUP`, `IDP`, `POSTURE`, `TRUSTED_NETWORK` . ZPA automatically assumes the posture profile ID that belongs to the parent tenant.

- `app_connector_groups` (Block Set)
  * `id` (String) The ID of an app connector group resource

- `app_server_groups` (Block Set)
  * `id` (String) The ID of a server group resource

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy access rule can be imported by using `<POLICY ACCESS RULE ID>` as the import ID.

For example:

```shell
terraform import zpa_policy_access_rule.example <policy_access_rule_id>
```

## LHS and RHS Values

| Object Type | LHS| RHS
|----------|-----------|----------
| [APP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_application_segment) | ``"id"`` | ``application_segment_id`` |
| [APP_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_segment_group) | ``"id"`` | ``segment_group_id``|
| [CLIENT_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_client_types) | ``"id"`` | ``zpn_client_type_zappl``, ``zpn_client_type_exporter``, ``zpn_client_type_browser_isolation``, ``zpn_client_type_ip_anchoring``, ``zpn_client_type_edge_connector``, ``zpn_client_type_branch_connector``,  ``zpn_client_type_zapp_partner``, ``zpn_client_type_zapp``  |
| [EDGE_CONNECTOR_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_cloud_connector_group) | ``"id"`` | ``<edge_connector_id>`` |
| [IDP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_idp_controller) | ``"id"`` | ``identity_provider_id`` |
| [SAML](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_saml_attribute) | ``saml_attribute_id``  | ``attribute_value_to_match`` |
| [SCIM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_attribute_header) | ``scim_attribute_id``  | ``attribute_value_to_match``  |
| [SCIM_GROUP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_scim_groups) | ``scim_group_attribute_id``  | ``attribute_value_to_match``  |
| [PLATFORM](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``mac``, ``ios``, ``windows``, ``android``, ``linux`` | ``"true"`` / ``"false"`` |
| [MACHINE_GRP](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_machine_group) | ``"id"`` | ``machine_group_id`` |
| [POSTURE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_posture_profile) | ``posture_udid``  | ``"true"`` / ``"false"`` |
| [TRUSTED_NETWORK](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_trusted_network) | ``network_id``  | ``"true"`` |
| [COUNTRY_CODE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/data-sources/zpa_access_policy_platforms) | [2 Letter ISO3166 Alpha2](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes)  | ``"true"`` / ``"false"`` |
| [RISK_FACTOR_TYPE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``ZIA``  | ``"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"`` |
| [CHROME_ENTERPRISE](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule) | ``managed``  | ``"true" / "false"`` |
