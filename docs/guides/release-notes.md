---
layout: "zscaler"
page_title: "Release Notes"
description: |-
  The Zscaler Private Access (ZPA) Provider Release Notes
---
# ZPA Provider: Release Notes

## USAGE

Track all ZPA Terraform provider's releases. New resources, features, and bug fixes will be tracked here.

---

``Last updated: v4.3.4``

---

## 4.3.4 (December, 3 2025)

### Notes

- Release date: **(December, 3 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #617](https://github.com/zscaler/terraform-provider-zpa/pull/617) - Fixed `zpa_application_segment_pra` and `zpa_application_segment_inspection` resources not storing `pra_app_id`/`inspect_app_id` in state, which prevented proper updates and deletes of app configurations.

## 4.3.3 (November, 11 2025)

### Notes

- Release date: **(November, 11 2025)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #612](https://github.com/zscaler/terraform-provider-zpa/pull/612) - Added `zpa_application_segment_weightedlb_config` resource and data source to manage weighted load balancer configuration for application segments, including automatic lookup by application and server group name.

- [PR #612](https://github.com/zscaler/terraform-provider-zpa/pull/612) - Added new attributes to `zpa_service_edge_group` attributes `exclusive_for_business_continuity`, and `city`

## 4.3.2 (November, 6 2025)

### Notes

- Release date: **(November, 6 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #611](https://github.com/zscaler/terraform-provider-zpa/pull/611) - Upgraded to [Zscaler SDK GO v3.8.3](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.8.3)

## 4.3.1 (November, 5 2025)

### Notes

- Release date: **(November, 5 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #610](https://github.com/zscaler/terraform-provider-zpa/pull/610) - Upgraded to [Zscaler SDK GO v3.8.2](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.8.2)

## 4.3.0 (October 31, 2025)

### Notes

- Release date: **(October 31, 2025)**
- Supported Terraform version: **v1.x**

### NEW - RESOURCES AND DATA SOURCES

The following new resources have been introduced:

- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added resource ``zpa_policy_portal_access_rule`` for Policy Portals access rule management
- [PR #588](https://github.com/zscaler/terraform-provider-zpa/pull/588) - Added resource and datasource``zpa_zia_cloud_config`` to configure ZIA Sandbox settings. This is required when configuring the following policy types: `zpa_policy_capabilities_rule` or `zpa_policy_portal_access_rule`
- [PR #588](https://github.com/zscaler/terraform-provider-zpa/pull/588) - Added resource and datasource``zpa_user_portal_aup`` for Acceptable Use Policy configuration
- [PR #588](https://github.com/zscaler/terraform-provider-zpa/pull/588) - Added resource and datasource``zpa_application_segment_multimatch_bulk`` to get application segments by domain that are incompatible with application segment Multimatch and Bulk updates application segment Multimatch in multiple applications segments

### NEW - DATA SOURCES

The following new data sources have been introduced:

- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource ``zpa_workload_tag_group``. This data source can be used when configuring `zpa_policy_access_rule` or `zpa_policy_access_rule_v2`, `object_type` is `WORKLOAD_TAG_GROUP`
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_risk_score_values`. This data source can be used when configuring policy types that support the `object_type` `RISK_SCORE`
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_managed_browser_profile`. This data source can be used when configuring `zpa_policy_access_rule_v2` or `zpa_policy_isolation_rule_v2` where the `object_type` is `CHROME_POSTURE_PROFILE`
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_browser_protection`. This data source can be used when configuring `zpa_policy_browser_protection_rule`.
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_branch_connector_group`. This data source can be used when configuring `zpa_policy_access_rule` or `zpa_policy_access_rule_v2`, `zpa_policy_forwarding_rule`, `zpa_policy_forwarding_rule_v2`, where the `object_type` is `BRANCH_CONNECTOR_GROUP`
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_extranet_resource_partner`. This data source is required when configuring resources such as: `zpa_server_group`, `zpa_application_segment`, `zpa_application_segmnent_pra`, `zpa_policy_access_rule_v2` in [Extranet mode](https://help.zscaler.com/zia/about-extranet)
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_location_controller`. This data source is required when configuring resources such as: `zpa_policy_access_rule_v2`, `zpa_server_group`
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_location_group_controller`. This data source is required when configuring resources such as: `zpa_policy_access_rule_v2`, `zpa_server_group`
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Added datasource `zpa_location_controller_summary`. This data source can be used when configuring `zpa_policy_access_rule` or `zpa_policy_access_rule_v2` where the `object_type` is `LOCATION`

### Enhancements

- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - The resource `zpa_server_group` now supports [Extranet mode](https://help.zscaler.com/zia/about-extranet) configuration.
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - The resource `zpa_policy_access_rule_v2` now supports [Extranet mode](https://help.zscaler.com/zia/about-extranet) configuration.
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - The resource `zpa_application_segment` now supports [Extranet mode](https://help.zscaler.com/zia/about-extranet) configuration.

### Bug Fixes

- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Set attribute `enabled` to default value `true` for the resource `zpa_app_connector_group` due to API setting value to `false` when configuring via API.
- [PR #608](https://github.com/zscaler/terraform-provider-zpa/pull/608) - Addressed edge case on `zpa_scim_groups` data source to handle group names containing special characther `@`.

## 4.2.6 (October 14, 2025)

### Notes

- Release date: **(October 14 2025)**
- Supported Terraform version: **v1.x**

### Enhancements

[PR #602](https://github.com/zscaler/terraform-provider-zpa/pull/602) - Implemented local caching for policy-controller resources for more efficiency on data source utilization.
[PR #602](https://github.com/zscaler/terraform-provider-zpa/pull/602) - Upgraded to [Zscaler-SDK-GO v3.7.5](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.7.5) to leverage better rate limiting and retry logic mechanism.

## 4.2.5 (October 8, 2025)

### Notes

- Release date: **(October 8 2025)**
- Supported Terraform version: **v1.x**

### Documentation

[PR #600](https://github.com/zscaler/terraform-provider-zpa/pull/600) - Updated index documentation for further clarity on ZPA Customer ID configuration
[PR #600](https://github.com/zscaler/terraform-provider-zpa/pull/600) - Added additional examples within the examples folder for resources and datasources `zpa_c2c_ip_ranges`, `zpa_private_cloud_controller`, `zpa_private_cloud_group`, `zpa_user_portal_controller`, `zpa_user_portal_link`

## 4.2.4 (September, 22 2025)

### Notes

- Release date: **(September, 22 2025)**
- Supported Terraform version: **v1.x**

### Enhancements

[PR #596](https://github.com/zscaler/terraform-provider-zpa/pull/596) - Introduced `GetByIdpAndAttributeID` in the data source `data_source_zpa_saml_attribute` to allow search IDP and Attribute ID. 
[PR #596](https://github.com/zscaler/terraform-provider-zpa/pull/596) - Made `id` attribute optional in some datasources to allow for search by ID.
[PR #596](https://github.com/zscaler/terraform-provider-zpa/pull/596) - Enhanced the Terraform schema in the following resources: `zpa_policy_inspection_rule` and `zpa_policy_isolation_rule`

## 4.2.3 (September, 22 2025)

### Notes

- Release date: **(September, 22 2025)**
- Supported Terraform version: **v1.x**

### Enhancements

[PR #596](https://github.com/zscaler/terraform-provider-zpa/pull/596) - Introduced `GetByIdpAndAttributeID` in the data source `data_source_zpa_saml_attribute` to allow search IDP and Attribute ID. 
[PR #596](https://github.com/zscaler/terraform-provider-zpa/pull/596) - Made `id` attribute optional in some datasources to allow for search by ID.
[PR #596](https://github.com/zscaler/terraform-provider-zpa/pull/596) - Enhanced the Terraform schema in the following resources: `zpa_policy_inspection_rule` and `zpa_policy_isolation_rule`

## 4.2.2 (September, 18 2025)

### Notes

- Release date: **(September, 18 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

[PR #593](https://github.com/zscaler/terraform-provider-zpa/pull/593) - Removed computed attribute `enabled` from `common_apps_dto.apps_config`. PRA application segments are always enabled. and can only be disabled by removing it entirely.

## 4.2.1 (September, 5 2025)

### Notes

- Release date: **(September, 5 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #590](https://github.com/zscaler/terraform-provider-zpa/pull/590) - Fixed `zpa_app_connector_group` resource attribute `lss_app_connector_group` due to API issues.

## 4.2.0 (August, 22 2025)

### Notes

- Release date: **(August, 222025)**
- Supported Terraform version: **v1.x**

### NEW - RESOURCES AND DATA SOURCES

The following new resources have been introduced:

- [PR #588](https://github.com/zscaler/terraform-provider-zia/pull/588) - Added and resource``zpa_c2c_ip_ranges`` - Added C2C IP Ranges
- [PR #588](https://github.com/zscaler/terraform-provider-zia/pull/588) - Added and resource``zpa_private_cloud_group`` - Added Private Cloud Group
- [PR #588](https://github.com/zscaler/terraform-provider-zia/pull/588) - Added and resource``zpa_user_portal_controller`` - Added User Portal Controller
- [PR #588](https://github.com/zscaler/terraform-provider-zia/pull/588) - Added and resource``zpa_user_portal_link`` - Added User Portal Link

### NEW - DATA SOURCES

The following new data sources have been introduced:

- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added and datasource``zpa_private_cloud_controller`` - Retrieves private cloud controller.

## 4.1.14 (July, 14 2025)

### Notes

- Release date: **(July, 14 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #583](https://github.com/zscaler/terraform-provider-zpa/pull/583) - Fixed attribute `servers` in the resource `zpa_server_groups` to prevent unexpected drifts.


## 4.1.13 (July, 2 2025)

### Notes

- Release date: **(July, 2 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #581](https://github.com/zscaler/terraform-provider-zpa/pull/581) - Fixed `zpa_policy_access_rule_reorder` resource to ensure it ignores `Zscaler Deception` order when it is not present.

## 4.1.12 (June, 27 2025)

### Notes

- Release date: **(June, 27 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #578](https://github.com/zscaler/terraform-provider-zpa/pull/578) - Added the attribute `share_to_microtenants` to resource `zpa_application_segment`. [Issue #577](https://github.com/zscaler/terraform-provider-zpa/issues/577)

## 4.1.11 (June, 20 2025)

### Notes

- Release date: **(June, 20 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #578](https://github.com/zscaler/terraform-provider-zpa/pull/578) - Added the attribute `share_to_microtenants` to resource `zpa_application_segment`. [Issue #577](https://github.com/zscaler/terraform-provider-zpa/issues/577)

## 4.1.10 (June, 19 2025)

### Notes

- Release date: **(June, 19 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #578](https://github.com/zscaler/terraform-provider-zpa/pull/578) - Added the attribute `share_to_microtenants` to resource `zpa_application_segment`. [Issue #577](https://github.com/zscaler/terraform-provider-zpa/issues/577)

## 4.1.9 (June, 11 2025)

### Notes

- Release date: **(June, 11 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #576](https://github.com/zscaler/terraform-provider-zpa/pull/576) - Updated the Policy Client Types: `zpn_client_type_zapp_partner`, `zpn_client_type_vdi`, `zpn_client_type_zia_inspection`

## 4.1.8 (June, 6 2025)

### Notes

- Release date: **(June, 6 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #348](https://github.com/zscaler/zscaler-sdk-go/pull/348) - Added new Policy Client Types: `zpn_client_type_zapp_partner`, `zpn_client_type_vdi`, `zpn_client_type_zia_inspection` to data source: `zpa_access_policy_client_types`

### Bug Fixes
[PR #574](https://github.com/zscaler/terraform-provider-zpa/pull/574) - Upgraded to [Zscaler GO SDK v3.4.4](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.4.3)
[PR #574](https://github.com/zscaler/terraform-provider-zpa/pull/574) - Resolved a crash in the Terraform provider when `use_legacy_client = true`. The SDK's `NewOneAPIClient()` function was performing OAuth2 authentication unconditionally, which caused the provider to hang or fail during legacy client initialization. The logic has been updated to skip authentication when the legacy client is in use.

## 4.1.7 (June, 5 2025)

### Notes

- Release date: **(June, 5 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #572](https://github.com/zscaler/terraform-provider-zpa/pull/572) - Upgraded to [Zscaler GO SDK v3.4.3](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.4.3)
[PR #572](https://github.com/zscaler/terraform-provider-zpa/pull/572) - Fixed pagination encoding due to API changes.

## 4.1.6 (May, 30 2025)

### Notes

- Release date: **(May, 30 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #571](https://github.com/zscaler/terraform-provider-zpa/pull/571) - Added new Access Policy Resource `zpa_policy_browser_protection_rule` to support ``USER_PORTAL`` `object_type`

### Documentation
[PR #571](https://github.com/zscaler/terraform-provider-zpa/pull/571) - Updated all access policy documentations and examples.

## 4.1.5 (May, 29 2025)

### Notes

- Release date: **(May, 29 2025)**
- Supported Terraform version: **v1.x**

### Enhancements
[PR #569](https://github.com/zscaler/terraform-provider-zpa/pull/569) - Added support for PRA User Portal with Zscaler Managed Certificate to resource `zpa_pra_portal_controller`
[PR #570](https://github.com/zscaler/terraform-provider-zpa/pull/570) - Added support for Zscaler Managed Certificate to resource `zpa_application_segment_browser_access`

## 4.1.4 (May, 22 2025)

### Notes

- Release date: **(May, 22 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #568](https://github.com/zscaler/terraform-provider-zpa/pull/568) - Fixed panic on the resource `zpa_policy_credential_rule` due to missing nil pointer on the attribute `credential` and `credential_pool`

## 4.1.3 (May, 14 2025)

### Notes

- Release date: **(May, 14 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #564](https://github.com/zscaler/terraform-provider-zpa/pull/564) - Fixed drift with attributes `reauth_idle_timeout` and `reauth_timeout` in the resource `zpa_policy_timeout_rule_v2`.
[PR #564](https://github.com/zscaler/terraform-provider-zpa/pull/564) - Modified behavior of `service_edges` block in `zpa_service_edge_group` resource:

### New Behavior
- **When omitted**: Terraform will ignore service edge membership completely (no drift detection)
- **When specified**: Terraform will enforce exact membership matching
  - You must include all required service edge IDs in the list
  - Any discrepancy between configuration and actual state will be reported as drift

### Important Notes
⚠️ **Deprecation Notice**: The `service_edges` block will be deprecated in a future release  
🔧 **External Management**: Service edge membership is typically managed outside Terraform  
💡 **Recommendation**: Only use this block if you require Terraform to explicitly manage membership

### Migration Guidance
If you're currently using this block but don't need strict membership control:
1. Remove the `service_edges` block from your configuration
2. Run `terraform apply` to update the state

## 4.1.2 (April, 30 2025)

### Notes

- Release date: **(April, 30 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #561](https://github.com/zscaler/terraform-provider-zpa/pull/561) - Fixed `zpa_server_group` panic during import process.
[PR #561](https://github.com/zscaler/terraform-provider-zpa/pull/561) - Fixed `zpa_service_edge_group` drift with nested blocks. Notice that the `trusted_networks` and `service_edges` blocks can only repeated once.

[PR #561](https://github.com/zscaler/terraform-provider-zpa/pull/561) - Enhanced nested flattening and expanding functions to provide more agnostic configuration when using dynamic blocks. The following resources are affected:
- `zpa_application_segment`
- `zpa_application_segment_pra`
- `zpa_application_segment_inspection`
- `zpa_application_segment_browser_access`
- `zpa_policy_access_rule`
- `zpa_policy_access_rule_v2`
- `zpa_policy_access_redirection_rule`
- `zpa_server_group`

## 4.1.1 (April, 29 2025)

### Notes

- Release date: **(April, 29 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #557](https://github.com/zscaler/terraform-provider-zpa/pull/557) - Fixed `zpa_server_group` panic during import process.
[PR #558](https://github.com/zscaler/terraform-provider-zpa/pull/558) - Fixed `zpa_service_edge_group` drift with nested blocks. Notice that the `trusted_networks` and `service_edges` blocks can only repeated once.

[PR #558](https://github.com/zscaler/terraform-provider-zpa/pull/558) - Enhanced nested flattening and expanding functions to provide more agnostic configuration when using dynamic blocks. The following resources are affected:
- `zpa_application_segment`
- `zpa_application_segment_pra`
- `zpa_application_segment_inspection`
- `zpa_application_segment_browser_access`
- `zpa_policy_access_rule`
- `zpa_policy_access_rule_v2`
- `zpa_policy_access_redirection_rule`
- `zpa_server_group`

## 4.1.0 (April, 23 2025)

### Notes

- Release date: **(April, 23 2025)**
- Supported Terraform version: **v1.x**

### NEW RESOURCE

[PR #552](https://github.com/zscaler/terraform-provider-zpa/pull/552) - The following new resource and data source have been introduced: `zpa_pra_credential_pool`. This resource creates a privileged remote access credential pool that can be referenced in an privileged credential access policy resource.

### Bug Fixes
[PR #552](https://github.com/zscaler/terraform-provider-zpa/pull/552) - Enhanced nested flattening and expanding functions to provide more agnostic configuration when using dynamic blocks. The following resources are affected:
- `zpa_application_segment`
- `zpa_application_segment_pra`
- `zpa_application_segment_inspection`
- `zpa_application_segment_browser_access`
- `zpa_policy_access_rule`
- `zpa_policy_access_rule_v2`
- `zpa_policy_access_redirection_rule`
- `zpa_server_group`

## 4.0.12 (April, 16 2025)

### Notes

- Release date: **(April, 16 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

### Bug Fixes
[PR #552](https://github.com/zscaler/terraform-provider-zpa/pull/552) - Fixed resource `zpa_application_segment_pra` to ensure proper PRA Application update and deletion process / cleanup.

## 4.0.11 (April, 14 2025)

### Notes

- Release date: **(April, 14 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
[PR #548](https://github.com/zscaler/zscaler-sdk-go/pull/548) - Set pointer in the `credential` block attribute in the ZPA `policysetcontrollerv2` resource.

## 4.0.10 (April, 11 2025)

### Notes

- Release date: **(April, 11 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #546](https://github.com/zscaler/terraform-provider-zpa/pull/546) - Fixed drift issue with the resources `zpa_policy_access_rule` and `zpa_policy_access_rule_v2` due to pre-populated attribute `custom_msg`.
- [PR #546](https://github.com/zscaler/terraform-provider-zpa/pull/546) - Fixed `flattenServiceEdgeSimple` and `flattenAppTrustedNetworksSimple` functions in the resource `zpa_service_edge_group` to prevent drifts due to block count ordering.
- [PR #546](https://github.com/zscaler/terraform-provider-zpa/pull/546) - Fixed documentation for `zpa_application_segment_pra` by removing the attribute `name` from within the `common_apps_dto.apps_config` block as it's not required.

## 4.0.9 (April, 8 2025)

### Notes

- Release date: **(April, 8 2025)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #544](https://github.com/zscaler/terraform-provider-zpa/pull/544) - Added new attribute `fqdn_dns_check` to all application segment resources.

### Bug Fixes

- [PR #544](https://github.com/zscaler/terraform-provider-zpa/pull/544) - Fixed `zpa_service_edge_group` attributes `service_edges` and `trusted_networks` to prevent drifts due to attribute ordering.

## 4.0.8 (March, 28 2025)

### Notes

- Release date: **(March, 28 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #543](https://github.com/zscaler/terraform-provider-zpa/pull/543) - Fixed `detachPRAConsoleFromPolicy` function within the resource `zpa_pra_console_controller` to ensure proper resource deletion flow.

## 4.0.7 (March, 25 2025)

### Notes

- Release date: **(March, 25 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #541](https://github.com/zscaler/terraform-provider-zpa/pull/541) - Upgraded to [Zscaler SDK GO v3.1.11](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.1.11) to fix url encoding.
- [PR #541](https://github.com/zscaler/terraform-provider-zpa/pull/541) - Fixed PRA Console, PRA Portal, PRA Credentials detachment functions to ensure proper removal from associated policies.

## 4.0.6 (March, 17 2025)

### Notes

- Release date: **(March, 17 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #537](https://github.com/zscaler/terraform-provider-zpa/pull/537) - Upgraded to [Zscaler SDK GO v3.1.10](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.1.8) to fix url encoding.

## 4.0.5 (March, 17 2025)

### Notes

- Release date: **(March, 17 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #536](https://github.com/zscaler/terraform-provider-zpa/pull/536) - Upgraded to [Zscaler SDK GO v3.1.9](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.1.8) to fix url encoding.


## 4.0.4 (March, 15 2025)

### Notes

- Release date: **(March, 15 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #534](https://github.com/zscaler/terraform-provider-zpa/pull/534) - Upgraded to [Zscaler SDK GO v3.1.8](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.1.8) to fix url encoding.

## 4.0.3 (March, 5 2025)

### Notes

- Release date: **(March, 5 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #529](https://github.com/zscaler/terraform-provider-zpa/pull/529) - Improved attributes `version_profile_name` and `version_profile_id` in the resource `zpa_app_connector_group`. Users can now use the attribute `version_profile_name` by providing the profile name and the provider will automatically retrieve and set the `version_profile_id`. The currently supported `version_profile_name` attribute values are:
  - ``Default``, ``Previous Default``, ``New Release``, ``Default - el8``, ``New Release - el8``, ``Previous Default - el8``

  **NOTE:** Users still leveraging the attribute `version_profile_id` directly can continue without impact to their existing configuration.

## 4.0.2 (February, 10 2025)

### Notes

- Release date: **(February, 10 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #522](https://github.com/zscaler/terraform-provider-zpa/pull/522) - Re-introduced attribute `cname` for the resource and data source `zpa_application_segment_browser_access` as a `Computed` only attribute.

### Internal Changes
- [PR #522](https://github.com/zscaler/terraform-provider-zpa/pull/522) - Updated `version.go` to `v4.0.2`
- [PR #522](https://github.com/zscaler/terraform-provider-zpa/pull/522) - Upgraded provider to [Zscaler SDK GO v3.1.4](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.1.4) to fix rate limiting override panic issues.

## 4.0.1 (January, 27 2025)

### Notes

- Release date: **(January, 27 2025)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #516](https://github.com/zscaler/terraform-provider-zpa/pull/516) - Removed `ConflictsWith` validation from `provider.go`. 

## 4.0.0 (January, 21 2025) - BREAKING CHANGES

### Notes

- Release date: **(January 21, 2025)**
- Supported Terraform version: **v1.x**

#### Enhancements - Zscaler OneAPI Support

[PR #515](https://github.com/zscaler/terraform-provider-zpa/pull/515): The ZPA Terraform Provider now offers support for [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) Oauth2 authentication through [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE** As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

⚠️ **WARNING**: Please refer to the [Index Page](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.md) page for details on authentication requirements prior to upgrading your provider configuration.

⚠️ **WARNING**: Attention Government customers. OneAPI and Zidentity is not currently supported for the following clouds: `GOV` and `GOVUS`. Refer to the [Legacy API Framework](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index) section for more information on how authenticate to these environments using the legacy method.

### NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS

#### NEW ATTRIBUTES
- [PR #515](https://github.com/zscaler/terraform-provider-zpa/pull/515) - Added new `object_type` `CHROME_ENTERPRISE` to the following ZPA access policy resources: `zpa_policy_access_rule`, and `zpa_policy_access_rule_v2`

## 3.332.0 (January, 8 2025)

### Notes

- Release date: **(January, 8 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #513](https://github.com/zscaler/terraform-provider-zpa/pull/513) - Upgraded Provider to SDK v2.74.2 to address Double encoding of special characters during GET operations.
- [PR #513](https://github.com/zscaler/terraform-provider-zpa/pull/513) -Fixed attribute `app_server_groups` on `zpa_policy_access_rule` resource to prevent innadivertent drifts during plan and apply. Issue [#512](https://github.com/zscaler/terraform-provider-zpa/pull/512)
- [PR #513](https://github.com/zscaler/terraform-provider-zpa/pull/513) - Deprecated previous `3.331.0` version due to missconfigured semversioning hash calculation.

**NOTE**  v3.331.0 and v3.332.0 was a versioning mistake due to backend issues and does not represent hundreds of new features Either version can be safely used without concerns on breaking changes. This will be corrected in the next major version release 4.0.0 upcoming in the next few weeks. We apologize for the confusion and  incovenience caused."

## 3.331.0 (January, 5 2025)

### Notes

- Release date: **(January, 5 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #509](https://github.com/zscaler/terraform-provider-zpa/pull/509) - Upgraded Provider to SDK v2.74.2 to address Double encoding of special characters during GET operations.

## 3.33.9 (October, 31 2024)

### Notes

- Release date: **(October, 31 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #500](https://github.com/zscaler/terraform-provider-zpa/pull/500) - Implemented a fix to the update function across all specialized application segment resources:
  - `zpa_application_segment_browser_access` - The fix now automatically includes the attributes `app_id` and `ba_app_id` in the payload during updates
  - `zpa_application_segment_inspection` - The fix now automatically includes the attributes `app_id` and `inspect_app_id` in the payload during updates
  - `zpa_application_segment_pra` - The fix now automatically includes the attributes `app_id` and `pra_app_id` in the payload during updates.
  **NOTE:** This update/fix is required to ensure the ZPA API can properly identify the Browser Access, Inspection App and PRA App, based on its specific custom ID. The fix also includes the removal of the `ForceNew` option previously included in the schema to force the resource replacement in case of changes. Issue [PR #498](https://github.com/zscaler/terraform-provider-zpa/pull/498)

## 3.33.8 (October, 29 2024)

### Notes

- Release date: **(October, 29 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #499](https://github.com/zscaler/terraform-provider-zpa/pull/499) - Fixed `zpa_application_segment_pra` import function and normalization of computed attributes.
- [PR #499](https://github.com/zscaler/terraform-provider-zpa/pull/499) - Fixed drift with attribute `health)check_type` in the resources `zpa_application_segment`, `zpa_application_segment_pra`, `zpa_application_segment_inspection` and `zpa_application_segment_browser_access`

### Enhancements
- [PR #499](https://github.com/zscaler/terraform-provider-zpa/pull/499) - Added new `zpa_application_segment` attribute `inspect_traffic_with_zia`

## 3.33.7 (October, 3 2024)

### Notes

- Release date: **(October, 3 2024)**
- Supported Terraform version: **v1.x**

### Enhancements
- [PR #496](https://github.com/zscaler/terraform-provider-zpa/pull/496) - Added new `object_type` `RISK_FACTOR_TYPE` to the following ZPA access policy resources: `zpa_policy_access_rule`, and `zpa_policy_access_rule_v2`

### Bug Fixes
- [PR #496](https://github.com/zscaler/terraform-provider-zpa/pull/496) - Fixed issue with attribute `tcp_port_range`/`udp_port_range` and `tcp_port_ranges`/`udp_port_ranges` within `zpa_application_segment`. The fix ensure that both port configuration formats are suported separately without mid-conversion in between. The fix also ensure the port configuration order is ignored during apply and update process. [Issue #490](https://github.com/zscaler/terraform-provider-zpa/issues/490).


### Internal Changes
- [PR #496](https://github.com/zscaler/terraform-provider-zpa/pull/496) Consolidated multiple functions supported common/cross-shared resources. The following new common functions were introduced for simplicity:
  - `expandCommonServerGroups`
  - `expandCommonAppConnectorGroups`
  - `expandCommonServiceEdgeGroups`
  - `flattenCommonAppConnectorGroups`
  - `flattenCommonAppServerGroups`
  - `flattenCommonServiceEdgeGroups`

## 3.33.6 (October, 1 2024)

### Notes

- Release date: **(October, 1 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #495](https://github.com/zscaler/terraform-provider-zpa/pull/495) - Fixed issue with attribute `tcp_port_range` and `udp_port_range` within the resource `zpa_application_segment` 

## 3.33.5 (September, 30 2024)

### Notes

- Release date: **(September, 30 2024)**
- Supported Terraform version: **v1.x**

### Enhancement
- [PR #493](https://github.com/zscaler/terraform-provider-zpa/pull/493) - Added plan stage validation for attributes `select_connector_close_to_app` and `bypass_type` in the resource `zpa_application_segment`.
- [PR #493](https://github.com/zscaler/terraform-provider-zpa/pull/493) - Added new attribute `use_in_dr_mode` in the resource `zpa_service_edge_group`.


## 3.33.4 (September, 23 2024)

### Notes

- Release date: **(September, 23 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #492](https://github.com/zscaler/terraform-provider-zpa/pull/492) - Fixed drift within attribute `zpa_policy_credential_access_rule`.
- [PR #492](https://github.com/zscaler/terraform-provider-zpa/pull/492) - Fixed detachement function within `zpa_segment_group`
  ~> **NOTE** This fix does not affect existing configurations.
  
## 3.33.3 (September, 18 2024)

### Notes

- Release date: **(September, 18 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #489](https://github.com/zscaler/terraform-provider-zpa/pull/489) - Fixed drift within attribute `common_apps_dto` and `pra_apps` in the `zpa_application_segment_pra` resource.
- [PR #489](https://github.com/zscaler/terraform-provider-zpa/pull/489) - Fixed drift within attribute `common_apps_dto` and `pra_apps` in the `zpa_application_segment_inspection` resource.
  ~> **NOTE** This fix does not affect existing configurations.
  
## 3.33.2 (September, 10 2024)

### Notes

- Release date: **(September, 10 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #486](https://github.com/zscaler/terraform-provider-zpa/pull/486) - Fixed drift related to common flattening and expand functions for all v2 Access policy resources.
  ~> **NOTE** This fix does not affect existing configurations using the `v2` policy type.
  
## 3.33.1 (September, 8 2024)

### Notes

- Release date: **(September, 8 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #484](https://github.com/zscaler/terraform-provider-zpa/pull/484) - Fixed drift within the resource `zpa_application_segment` related to the attribute `microtenant_id` when setting the microtenant ID via environment variable.

## 3.33.0 (September, 5 2024)

### Notes

- Release date: **(September, 5 2024)**
- Supported Terraform version: **v1.x**

### Enhancements
- [PR #483](https://github.com/zscaler/terraform-provider-zpa/pull/483) Updated `resourceSegmentGroupUpdate` function in the resource `zpa_segment_group` to use the new GO SDK function `UpdateV2`. The `UpdateV2` function offers a newly enhanced v2 ZPA API endpoint.

### Bug Fixes
- [PR #483](https://github.com/zscaler/terraform-provider-zpa/pull/483) - Fixed drift issue within all policy access rule v2 resources.
- [PR #483](https://github.com/zscaler/terraform-provider-zpa/pull/483) - Fixed drift within the resource `zpa_provisioning_key` related to the attribute `microtenant_id` when setting the microtenant ID via environment variable.

## 3.32.5 (August, 28 2024)

### Notes

- Release date: **(August, 28 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #481](https://github.com/zscaler/terraform-provider-zpa/pull/481) - Fixed drift issue within the attribute `tcp_port_ranges` and `udp_port_ranges` for the resource `zpa_application_segment` to ignore the order which the port numbers are configured.

## 3.32.4 (August, 26 2024)

### Notes

- Release date: **(August, 26 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #478](https://github.com/zscaler/terraform-provider-zpa/pull/478) - Fixed drift within the attribute `service_edge` for the resource `zpa_service_edge_group` to ignore the order of IDs

### Enhancements
- [PR #478](https://github.com/zscaler/terraform-provider-zpa/pull/478) Added new attributes to `privileged_approvals_enabled` to resource: `zpa_microtenant_controller`. The attribute indicates if Privileged Approvals is enabled (true) for the Microtenant.

## 3.32.3 (August, 22 2024)

### Notes

- Release date: **(August, 22 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #476](https://github.com/zscaler/terraform-provider-zpa/pull/476) - Fixed validation for the `match_style` attribute in the resource `zpa_application_segment`. 

## 3.32.2 (August, 16 2024)

### Notes

- Release date: **(August, 16 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #476](https://github.com/zscaler/terraform-provider-zpa/pull/476) - Fixed validation for the `match_style` attribute in the resource `zpa_application_segment`. 
  **NOTE**: Notice that `match_style` also known as [Multimatch](https://help.zscaler.com/zpa/using-app-segment-multimatch) cannot be configured when `ip_anchored` is enabled. Also, `match_style` is NOT supported for the following application segment resources: `zpa_application_segment_browser_access`, `zpa_application_segment_inspection` or `zpa_application_segment_pra`.

### Documentation
- [PR #476](https://github.com/zscaler/terraform-provider-zpa/pull/476) - Added documentation for resource and datasource: `zpa_service_edge_assistant_schedule`

## 3.32.1 (July, 31 2024)

### Notes

- Release date: **(July, 31 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #473](https://github.com/zscaler/terraform-provider-zpa/pull/473) Removed unsupported attributes `microtenant_id` and `microtenant_name` from `zpa_application_segment_inspection` resource and data source.
  ~> **NOTE** Although that's a safe change, it may cause a temporary drift in order to update the statefile. Microtenant is not currently supported for Inspection Application Segments

 - [PR #473](https://github.com/zscaler/terraform-provider-zpa/pull/473) Added missing `microtenant_id` attribute to nested block `common_apps_dto.apps_config` in the resource and data source `zpa_application_segment_pra`.
  ~> **NOTE** Although that's a safe change, it may cause a temporary drift in order to update the statefile.

 - [PR #473](https://github.com/zscaler/terraform-provider-zpa/pull/473) Added missing `microtenant_id` attribute to nested block `clientless_apps` in the resource and data source `zpa_application_segment_browser_access`.
  ~> **NOTE** Although that's a safe change, it may cause a temporary drift in order to update the statefile.

 - [PR #473](https://github.com/zscaler/terraform-provider-zpa/pull/473) Fixed drift related to common flattening and expand functions for all v2 Access policy resources.
  ~> **NOTE** This fix does not affect existing configurations using the `v2` policy type.

### Documentation
- [PR #473](https://github.com/zscaler/terraform-provider-zpa/pull/473) Added documentation examples for the following resources and datasources:
  - ``zpa_service_edge_assistant_schedule``
  - ``zpa_policy_credential_rule``

## 3.32.0 (July, 24 2024)

### Notes

- Release date: **(July, 24 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #473](https://github.com/zscaler/terraform-provider-zpa/pull/473) Fixed drift issues with the following resources:
  - `zpa_inspection_profile`
  - `zpa_policy_access_inspection_rule_v2`
  - `zpa_pra_approval`

### Documentation
- [PR #473](https://github.com/zscaler/terraform-provider-zpa/pull/473) Added documentation for the following resource:
  - `zpa_policy_redirection_rule`
  
## 3.31.0 (July, 11 2024)

### Notes

- Release date: **(July, 11 2024)**
- Supported Terraform version: **v1.x**

### Deprecations
- [PR #471](https://github.com/zscaler/terraform-provider-zpa/pull/468) The following attributes are not deprecated:
  - ``zpa_application_segment_browser_access``: `cname`, `hidden`, `local_name`, `path`
  - ``zpa_application_segment_pra``: `cname`, `segment_group_name`

### Bug Fixes
- [PR #471](https://github.com/zscaler/terraform-provider-zpa/pull/468) Fixed `zpa_application_segment_inspection` drift issues within `common_apps_dto` and `tcp_port_range`
- [PR #471](https://github.com/zscaler/terraform-provider-zpa/pull/468) Fixed `zpa_inspection_custom_controls` drift issues `protocol_type` attribute
- [PR #471](https://github.com/zscaler/terraform-provider-zpa/pull/468) Fixed `zpa_inspection_custom_controls`import issues.

### Enhancements
- [PR #471](https://github.com/zscaler/terraform-provider-zpa/pull/468) Added new attributes to `zpa_cloud_browser_isolation_external_profile`:
  - `flattened_pdf` - Enable to allow downloading of flattened files from isolation container to your local computer.

    **NOTE** `flattened_pdf` must be set to `false` when `upload_download` is set to `all`
- `security_controls` - The CBI security controls enabled for the profile
  - `copy_paste:` - Enable or disable copy & paste for local computer to isolation. Supported values are: `none` or `all`
  - `document_viewer:` - Enable or disable to view Microsoft Office files in isolation.
  - `local_render:` - Enables non-isolated hyperlinks to be opened on the user's native browser.
  - `upload_download` - Enable or disable file transfer from local computer to isolation. Supported values are: `none`, `all`, `upstream`

    **NOTE** `upload_download` must be set to `none` or `upstream` when `flattened_pdf` is set to `true`

  - `deep_link:` - Enter applications that are allowed to launch outside of the Isolation session
    - `enabled:` - Enable or disable to view Microsoft Office files in isolation.
    - `applications:` - List of deep link applications

  - `watermark:` - Enable to display a custom watermark on isolated web pages.
    - `enabled:` - Enable to display a custom watermark on isolated web pages.
    - `show_user_id:` - Display the user ID on watermark isolated web pages.
    - `show_timestamp:` - Display the timestamp on watermark isolated web pages.
    - `show_message:` - Enable custom message on watermark isolated web pages.
    - `message:` - Display custom message on watermark isolated web pages.

- `user_experience` - The CBI security controls enabled for the profile
  - `forward_to_zia:` - Enable to forward non-ZPA Internet traffic via ZIA.
    - `enabled:` - Enable to forward non-ZPA Internet traffic via ZIA.
    - `organization_id:` - Use the ZIA organization ID from the Company Profile section.
    - `cloud_name:` - The ZIA cloud name on which the organization exists i.e `zscalertwo`
    - `pac_file_url:` - Enable to have the PAC file be configured on the Isolated browser to forward traffic via ZIA.

- `debug_mode`- Enable to allow starting isolation sessions in debug mode to collect troubleshooting information.
  - `allowed:` - Enable to allow starting isolation sessions in debug mode to collect troubleshooting information.
  - `file_password:` - Set an optional password to debug files when this mode is enabled.

## 3.3.25 (July, 2 2024)

### Notes

- Release date: **(July, 2 2024)**
- Supported Terraform version: **v1.x**

### Deprecations
- [PR #468](https://github.com/zscaler/terraform-provider-zpa/pull/468) The following attributes are not deprecated:
  - ``zpa_application_segment_browser_access``: `cname`, `hidden`, `local_name`, `path`
  - ``zpa_application_segment_pra``: `cname`, `segment_group_name`

### Bug Fixes
- [PR #468](https://github.com/zscaler/terraform-provider-zpa/pull/468) Fixed `zpa_application_segment_inspection` drift issues within `common_apps_dto` and `tcp_port_range`

## 3.3.24 (June, 14 2024)

### Notes

- Release date: **(June, 14 2024)**
- Supported Terraform version: **v1.x**

### Internal Changes
- [PR #464](https://github.com/zscaler/terraform-provider-zpa/pull/464) Upgraded to [Zscaler-SDK-GO](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.61.0). The upgrade supports easier ZPA API Client instantiation for existing and new resources.
- [PR #464](https://github.com/zscaler/terraform-provider-zpa/pull/464) Upgraded ``releaser.yml`` to [GoReleaser v6](https://github.com/goreleaser/goreleaser-action/releases/tag/v6.0.0)

## 3.3.23 (May, 31 2024)

### Notes

- Release date: **(May, 31 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
Upgraded to Zscaler SDK GO v2.5.31 to address new ZPA error handling to retry on new `400` and `409` error format message:

```json
  "id" : "api.concurrent.access.error",
  "reason" : "Unable to modify the resource due to concurrent change requests. Try again"
```

## 3.3.22 (May, 24 2024)

### Notes

- Release date: **(May, 24 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes
- [PR #459](https://github.com/zscaler/terraform-provider-zpa/pull/459) Fixed panic issue with attribute `trusted_networks` within the resource `zpa_service_edge_group`.

## 3.3.21 (May, 18 2024)

### Notes

- Release date: **(May, 18 2024)**
- Supported Terraform version: **v1.x**

### ENHACEMENTS
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Added new data source `zpa_application_segment_by_type`. The data source allows for querying of application segments by type. The ``application_type`` attribute supports the following values: `BROWSER_ACCESS`, `INSPECT`, and `SECURE_REMOTE_ACCESS`

### Bug Fixes
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Fixed resource `zpa_service_edge_group` due to misconfiguration in the importing function.

### Internal Changes
- [PR #454](https://github.com/zscaler/terraform-provider-zpa/pull/454) - Added Support to arbitrary clouds for testing purposes

## 3.3.2 (May, 18 2024)

### Notes

- Release date: **(May, 18 2024)**
- Supported Terraform version: **v1.x**

### ENHACEMENTS
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Added new data source `zpa_application_segment_by_type`. The data source allows for querying of application segments by type. The ``application_type`` attribute supports the following values: `BROWSER_ACCESS`, `INSPECT`, and `SECURE_REMOTE_ACCESS`

### Bug Fixes
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Fixed resource `zpa_service_edge_group` due to misconfiguration in the importing function.

### Internal Changes
- [PR #454](https://github.com/zscaler/terraform-provider-zpa/pull/454) - Added Support to arbitrary clouds for testing purposes

## 3.3.1 (May, 18 2024)

### Notes

- Release date: **(May, 18 2024)**
- Supported Terraform version: **v1.x**

### ENHACEMENTS
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Added new data source `zpa_application_segment_by_type`. The data source allows for querying of application segments by type. The ``application_type`` attribute supports the following values: `BROWSER_ACCESS`, `INSPECT`, and `SECURE_REMOTE_ACCESS`

### Bug Fixes
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Fixed resource `zpa_service_edge_group` due to misconfiguration in the importing function.

### Internal Changes
- [PR #454](https://github.com/zscaler/terraform-provider-zpa/pull/454) - Added Support to arbitrary clouds for testing purposes

## 3.3.0 (May, 17 2024)

### Notes

- Release date: **(May, 17 2024)**
- Supported Terraform version: **v1.x**

### ENHACEMENTS
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Added new data source `zpa_application_segment_by_type`. The data source allows for querying of application segments by type. The ``application_type`` attribute supports the following values: `BROWSER_ACCESS`, `INSPECT`, and `SECURE_REMOTE_ACCESS`

### Bug Fixes
- [PR #455](https://github.com/zscaler/terraform-provider-zpa/pull/455) Fixed resource `zpa_service_edge_group` due to misconfiguration in the importing function.

### Internal Changes
- [PR #454](https://github.com/zscaler/terraform-provider-zpa/pull/454) - Added Support to arbitrary clouds for testing purposes

## 3.2.11 (May, 3 2024)

### Notes

- Release date: **(May, 3 2024)**
- Supported Terraform version: **v1.x**

### Internal Changes

- [PR #449](https://github.com/zscaler/terraform-provider-zpa/pull/449) - Added `CodeCov` Support to GitHub Workflow 

### Bug Fixes
- [PR #450](https://github.com/zscaler/terraform-provider-zpa/pull/450) - Implemented additional validation within the resource `zpa_policy_access_rule_reorder` to ensure it accounts for the potential existence of the `Zscaler Deception` rule. [Zscaler API Documentation](https://help.zscaler.com/zpa/configuring-access-policies-using-api#:~:text=Updating%20the%20rule,configured%20using%20Deception.) for further details.

⚠️ **WARNING:**: This change does not affect existing rule configurations, and is only applicable for tenants with the Zscaler Deception rule configured. If your tenant have this rule configured, please refer to the [provider documentation](https://registry.terraform.io/providers/zscaler/zpa/latest/docs/resources/zpa_policy_access_rule_reorder) for further examples on how you can address potential drift issues due to rule order missmatch. [Issue #445](https://github.com/zscaler/terraform-provider-zpa/issues/445)

### ENHACEMENTS
- [PR #450](https://github.com/zscaler/terraform-provider-zpa/pull/450) - The resource `zpa_service_edge_group` now supports the following new attributes:
  * `grace_distance_enabled`: Allows ZPA Private Service Edge Groups within the specified distance to be prioritized over a closer ZPA Public Service Edge.
  * `grace_distance_value`: Indicates the maximum distance in miles or kilometers to ZPA Private Service Edge groups that would override a ZPA Public Service Edge.
  * `grace_distance_value_unit`: Indicates the grace distance unit of measure in miles or kilometers. This value is only required if `grace_distance_enabled` is set to true. Support values are: `MILES` and `KMS`

### Documentation
- [PR #450](https://github.com/zscaler/terraform-provider-zpa/pull/450) - Updated documentation for `zpa_policy_access_rule_reorder` by removing deprecated `policy_set_id` attribute from the resource. Only the `policy_type` is required.
### Documentation
- [PR #450](https://github.com/zscaler/terraform-provider-zpa/pull/450) - Updated documentation for `zpa_service_edge_group` by including detailed description of the new attributes: `grace_distance_enabled`, `grace_distance_value`, `grace_distance_value_unit`.

## 3.2.1 (April, 8 2024)

### Notes

- Release date: **(April, 8  2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #442](https://github.com/zscaler/terraform-provider-zpa/pull/442) - Fixed `zpa_ba_certificate` resource and aligned with `zpa_application_segment_browser_access` `certificate_id` attribute. 

  !> **WARNING:** Notice that updating the ``cert_blob`` attribute in the `zpa_ba_certificate` will trigger a full replacement of both the certificate and the `zpa_application_segment_browser_access`  along with any access policy the application segment may be associated with.

## 3.2.0 (April, 3 2024)

### Notes

- Release date: **(April, 3 2024)**
- Supported Terraform version: **v1.x**

### NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES:

### NEW RESOURCES AND DATASOURCES:
* New datasource: `zpa_pra_approval_controller` retrieve Privileged Remote Access Approval [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New resource: `zpa_pra_approval_controller` manages Privileged Remote Access Approval [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New datasource: `zpa_pra_portal_controller` retrieve Privileged Remote Access Portal [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New resource: `zpa_pra_portal_controller` manages Privileged Remote Access Portal [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New datasource: `zpa_pra_credential_controller` retrieve Privileged Remote Access Credential [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New resource: `zpa_pra_credential_controller` manages Privileged Remote Access Credential [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New datasource: `zpa_pra_console_controller` retrieve Privileged Remote Access Console [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New resource: `zpa_pra_console_controller` manages Privileged Remote Access Console [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
* New Resources: Introduced new Policy Access resources that are managed via a new `v2` API endpoint:
  - `zpa_policy_access_rule_v2` manages access policy rule via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
  - `zpa_policy_forwarding_rule_v2` manages access policy forwarding rule via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
  - `zpa_policy_isolation_rule_v2` manages access policy isolation rule via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
  - `zpa_policy_inspection_rule_v2` manages access policy inspection rule via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
  - `zpa_policy_timeout_rule_v2` manages access policy timeout rule via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
  - `zpa_policy_redirection_rule` manages redirection access policy via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/425)
  - `zpa_policy_credential_rule` manages access policy credential rule via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
  - `zpa_policy_capabilities_rule` manages access policy capabilities rule via `v2` API endpoint [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
  
    ⚠️ **WARNING:**: Notice that any Access Policy `v2` is a new resource and uses a different HCL format structure. If you decide to migrate to the new v2 resources, notice that this is considered a breaking change and must be done carefully. This warning only applies for those with existing `v1` Access Policy HCL format structure.

[PR #434](https://github.com/zscaler/terraform-provider-zpa/pull/434)
* New resource: `zpa_emergency_access_user` manages Emergency Access Users 

### NEW PROPERTIES
* New Properties: The resource `zpa_ba_certificate` now displays the attributes `valid_from_in_epochsec` and `valid_to_in_epochsec` in human readable `RFC1123` format
* New Properties: The provider now includes support to `ZPATWO` cloud [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)

### DEPRECATIONS
* Deprecated attribute: The attributes `policy_migrated` and `tcp_keep_alive_enabled` are now deprecated for the resource `zpa_segment_group`. For the attribute `tcp_keep_alive_enabled` use the attribute `tcp_keep_alive` within the resource  `zpa_application_segment`", [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432). 
* Deprecated attribute: The attributes `negated` within all access policy rule resource types. [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432). 
* Deprecated attribute: The attributes `rule_order` within all access policy rule resource types. Please use the newly dedicated resource `zpa_policy_access_rule_reorder` [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432). 

### ENHACEMENTS
* Attribute `policy_set_id` is now optional across all access policy rule resources `v1` and `v2`. The provider will automatically set the `policy_set_id` according to the policy access resource being configured. This improvement removes the need to explicitly use the data source `zpa_policy_type` [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
* Added new `match_style` attribute to the `zpa_application_segment` resource [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432). Issue [#424](https://github.com/zscaler/terraform-provider-zpa/issues/424). To learn more about this attribute visit [Zscaler Help Portal](https://help.zscaler.com/zpa/using-app-segment-multimatch)
* Update `zpa_ba_certificate` documentation [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)
* Several ACC tests maintenance [PR #432](https://github.com/zscaler/terraform-provider-zpa/pull/432)

## 3.1.1 (February, 28 2024)

### Notes

- Release date: **(February, 28 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #423](https://github.com/zscaler/terraform-provider-zpa/pull/423) - Fixed drift issue within `zpa_application_segment_pra` resource

## 3.1.0 (January, 17 2024) - Unreleased

### Notes

- Release date: **(January, 17 2024)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #394](https://github.com/zscaler/terraform-provider-zpa/pull/394) - ✨ Added support for ZPA Certificate provisioning
- [PR #405](https://github.com/zscaler/terraform-provider-zpa/pull/405) - ✨ Added support for ZPA Assistant Schedule feature to configures Auto Delete for the specified disconnected App Connectors.
- [PR #389](https://github.com/zscaler/terraform-provider-zpa/pull/389) - ✨ Added support to New ZPA Bulk Reorder Policy Rule

### Fixes

- [PR #391](https://github.com/zscaler/terraform-provider-zpa/pull/391) - Removed `enrollment_cert_name` computed attribute from provisioning key resource

## 3.0.5 (November, xx 2023)

### Notes

- Release date: **(November, xx 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #388](https://github.com/zscaler/terraform-provider-zpa/pull/388) - Updated provider to zscaler-sdk-go v2.1.6 to support ZPA SCIM Group SortOrder and SortBy search criteria option
- [PR #389](https://github.com/zscaler/terraform-provider-zpa/pull/389) - Added support for new ZPA Access Policy Bulk Reorder Endpoint

## 3.0.4 (November, 6 2023)

### Notes

- Release date: **(November, 6 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #385](https://github.com/zscaler/terraform-provider-zpa/pull/385) - Fixed `microtenant_id` attribute for all access policy types.
  ⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.
- [PR #383](https://github.com/zscaler/terraform-provider-zpa/pull/383) - Fixed issues with hard-coded authentication within the provider block.

## 3.0.3 (October, 27 2023)

### Notes

- Release date: **(October, 27 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #375](https://github.com/zscaler/terraform-provider-zpa/pull/375) - Fixed drift issues in ``zpa_application_segment_pra`` and ``zpa_application_segment_inspection`` when setting up ``apps_config`` options.
- [PR #375](https://github.com/zscaler/terraform-provider-zpa/pull/375) - Upgrade to [Zscaler-SDK-GO v2.1.3](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.1.3). The upgrade allows searches for resources in which the name include 1 or more spaces.
- [PR #380](https://github.com/zscaler/terraform-provider-zpa/pull/380) - Fixed provider authentication to accept `ZPA_CLOUD` via environment variables.
- [PR #381](https://github.com/zscaler/terraform-provider-zpa/pull/381) - Included and fixed additional acceptance test cases for several resources and datasources

## 3.0.2 (September, 30 2023)

### Notes

- Release date: **(September, 30 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #374](https://github.com/zscaler/terraform-provider-zpa/pull/374) - Resource `zpa_lss_config_controller` now supports ability to configure granular access policies via the embbeded `policy_type` `SIEM_POLICY`.

### Fixes

- [PR #372](https://github.com/zscaler/terraform-provider-zpa/pull/372) - Provider HTTP Header now includes enhanced ``User-Agent`` information for troubleshooting assistance.
  - i.e ``User-Agent: (darwin arm64) Terraform/1.5.5 Provider/3.0.2 CustomerID/xxxxxxxxxxxxxxx``

## 3.0.1-beta (September, 21 2023)

### Notes

- Release date: **(September, 21 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #369](https://github.com/zscaler/terraform-provider-zpa/pull/369) - Added fix to resource `zpa_policy_access_rule_reorder` to support multiple policy types. The reorder operation is now supported for the following policy types:
  - ``ACCESS_POLICY or GLOBAL_POLICY``
  - ``TIMEOUT_POLICY or REAUTH_POLICY``
  - ``BYPASS_POLICY or CLIENT_FORWARDING_POLICY``
  - ``INSPECTION_POLICY``
  - ``ISOLATION_POLICY``
  - ``CREDENTIAL_POLICY``
  - ``CAPABILITIES_POLICY``
  - ``CLIENTLESS_SESSION_PROTECTION_POLICY``

- [PR #371](https://github.com/zscaler/terraform-provider-zpa/pull/371) - Fixed ``object_type`` validation for all supported policy types.

## 3.0.0-beta (September, 18 2023)

### Notes

- Release date: **(September, 18 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #355](https://github.com/zscaler/terraform-provider-zpa/pull/355) - Introduced the new resource and datasource `zpa_microtenant_controller`
- [PR #355](https://github.com/zscaler/terraform-provider-zpa/pull/355) - Added support to the new Microtenant Controller feature to the following resources:
  - `zpa_app_connector_controller`,, `zpa_app_connector_group`, `zpa_application_segment`, `zpa_application_segment_browser_access`, `zpa_application_segment_inspection`, `zpa_application_segment_pra`, `zpa_application_server`, `zpa_policy_type`, `zpa_policy_access_rule`, `zpa_policy_access_forwarding_rule`, `zpa_policy_access_timeout_rule`, `zpa_policy_access_inspection_rule`, `zpa_policy_access_isolation_rule`, `zpa_segment_group`, `zpa_server_group`, `zpa_provisioning_key`, `zpa_machine_group`, `zpa_service_edge_group`, `zpa_service_edge_controller`

⚠️ **WARNING:**: The new attribute ``microtenant_id`` is optional. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

⚠️ **WARNING:**: This feature is in limited availability and requires additional license. To learn more, contact Zscaler Support or your local account team.

- [PR #356](https://github.com/zscaler/terraform-provider-zpa/pull/356) - Added support to the following new ZPA Cloud Browser Isolation resources and datasources:

- Resources
  - `zpa_cloud_browser_isolation_banner` - Cloud Browser Isolation Banner Controller
  - `zpa_cloud_browser_isolation_certificate` - Cloud Browser Isolation Certificate Controller
  - `zpa_cloud_browser_isolation_external_profile` - Cloud Browser Isolation External Profile Controller

- Data Sources
  - `zpa_cloud_browser_isolation_banner` - Cloud Browser Isolation Banner Controller
  - `zpa_cloud_browser_isolation_certificate` - Cloud Browser Isolation Certificate Controller
  - `zpa_cloud_browser_isolation_external_profile` - Cloud Browser Isolation External Profile Controller
  - `zpa_cloud_browser_isolation_region` - Cloud Browser Isolation Regions
  - `zpa_cloud_browser_isolation_zpa_profile` - Cloud Browser Isolation ZPA Profile

  ⚠️ **WARNING:**: Cloud Browser Isolation (CBI) is a licensed feature flag. Please contact Zscaler support or your local account team for details.

- [PR #363](https://github.com/zscaler/terraform-provider-zpa/pull/363) - Added support for `COUNTRY_CODE` object type within the `zpa_policy_access_rule` resource. The provider validates the use of proper 2 letter country codes [ISO3166 By Alpha2Code](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) - Issue [#361](https://github.com/zscaler/terraform-provider-zpa/issues/361)

- [PR #366](https://github.com/zscaler/terraform-provider-zpa/pull/366) - Added ISO3166 Alpha2Code for ``country_code`` validation on `zpa_app_connector_groups` and `zpa_service_edge_group` resources

## 2.83.0-beta (September, 5 2023) - Beta

### Notes

- Release date: **(September, 5 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #350](https://github.com/zscaler/terraform-provider-zpa/pull/350) - Update provider to [Zscaler SDK GO](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v1.8.0-beta) v1.8.0-beta. This version provides caching  mechanism, which aims to enhance the provider performance as well as decrease the number of API calls being made to the ZPA API.

### Fixes

- [PR #350](https://github.com/zscaler/terraform-provider-zpa/pull/350) - Fixed drift within Access Policy Condition to ensure update is performed when adding and removing application segments.

## 2.82.4-beta (August, 18 2023) - Beta

### Notes

- Release date: **(August, 18 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #348](https://github.com/zscaler/terraform-provider-zpa/pull/348)

⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of this resource for all ZPA policy types.

## 2.82.3-beta (August, 18 2023) - Beta

### Notes

- Release date: **(August, 18 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #347](https://github.com/zscaler/terraform-provider-zpa/pull/347)

⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of this resource for all ZPA policy types.

## 2.82.2-beta (August, 17 2023) - Beta

### Notes

- Release date: **(August, 17 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #345](https://github.com/zscaler/terraform-provider-zpa/pull/345)
  1. ``zpa_policy_access_rule_reorder`` Added check to prevent ``order <= 0``
  2. ``zpa_policy_access_rule_reorder`` Added check to prevent non-contigous (gaps) in rule order numbers
  3. ``zpa_policy_access_rule_reorder`` Added check to prevent rule order number to be greater than the total number of rules being configured.

⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of this resource for all ZPA policy types.

## 2.82.1-beta (August, 16 2023) - Beta

### Notes

- Release date: **(August, 16 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #344](https://github.com/zscaler/terraform-provider-zpa/pull/344)
  1. Implemented a new resource ``zpa_policy_access_rule_reorder`` to support Access policy rule reorder in a more efficient way.

⚠️ **WARNING:**: The attribute ``rule_order`` is now deprecated in favor of this resource for all ZPA policy types.

## 2.81.0 (August, 1 2023)

### Notes

- Release date: **(August, 1 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #334](https://github.com/zscaler/terraform-provider-zpa/pull/326) - Added support to ZPA ``GOVUS`` environment. Issue [#333](https://github.com/zscaler/terraform-provider-zpa/issues/333)

## 2.8.0 (July, 5 2023)

### Notes

- Release date: **(July, 5 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #325](https://github.com/zscaler/terraform-provider-zpa/pull/325) - Added new attribute ``waf_disabled`` to resource ``zpa_app_connector_group``
- [PR #326](https://github.com/zscaler/terraform-provider-zpa/pull/326) - Added support to ZPA ``QA`` environment

### Fixes

- [PR #319](https://github.com/zscaler/terraform-provider-zpa/pull/319) - Fixed links to Zenith Community demo videos in the documentation
- [PR #321](https://github.com/zscaler/terraform-provider-zpa/pull/321) - Fixed resource ``zpa_server_group``due to panic when set attribute ``dynamic_discovery`` to false.
- [PR #323](https://github.com/zscaler/terraform-provider-zpa/pull/323) - Fixed attribute ``server_groups`` in all ``zpa_application_segment`` resources due to server group ID reorder, which caused drift behavior. Issue [#322](https://github.com/zscaler/terraform-provider-zpa/issues/322)

## 2.7.9 (June, 10 2023)

### Notes

- Release date: **(June, 10 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- Updated to Zscaler-SDK-GO v1.5.5. The update improves search mechanisms for ZPA resources, to ensure streamline upstream GET API requests and responses using ``search`` parameter. Notice that not all current API endpoints support the search parameter, in which case, all resources will be returned.

## 2.7.8 (June, 3 2023)

### Notes

- Release date: **(June, 3 2023)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #311](https://github.com/zscaler/terraform-provider-zpa/pull/311) Fixed ZPA resource ``Service Edge Group`` and ``Service Edge Controller`` Struct to support attribute ``publish_ips``.
- [PR #314](https://github.com/zscaler/terraform-provider-zpa/pull/314) Fixed ``rhs`` attribute within the ``GetPolicyConditionsSchema``function to prevent invalid new value inconsistency issue.

## 2.7.7 (May, 23 2023)

### Notes

- Release date: **(May, 23 2023)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #309](https://github.com/zscaler/terraform-provider-zpa/pull/309) Updated provider to Zscaler SDK GO v1.5.2. The update added exception handling within the ZPA API Client to deal with simultaneous DB requests, which were affecting the ZPA Policy Access rule order creation.

⚠️ **WARNING:** Due to API restrictions, we recommend to limit the number of requests to ONE, when configuring the following resources:

- ``zpa_policy_access_rule``
- ``zpa_policy_inspection_rule``
- ``zpa_policy_timeout_rule``
- ``zpa_policy_forwarding_rule``
- ``zpa_policy_isolation_rule``
  - Internal References:
    - [ET-53585](https://jira.corp.zscaler.com/browse/ET-53585)
    - [ET-48860](https://confluence.corp.zscaler.com/display/ET/ET-48860+incorrect+rules+order)

Terraform uses goroutines to speed up deployment, but the number of parallel
operations it launches may exceed [what is recommended](https://help.zscaler.com/zpa/about-rate-limiting).
When configuring ZPA Policies we recommend to limit the number of concurrent API calls to **ONE**. This limit ensures that there is no performance impact during the provisioning of large Terraform configurations involving access policy creation.

This recommendation applies to the following resources:

- ``zpa_policy_access_rule``
- ``zpa_policy_inspection_rule``
- ``zpa_policy_timeout_rule``
- ``zpa_policy_forwarding_rule``
- ``zpa_policy_isolation_rule``

In order to accomplish this, we recommend setting the [parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at this limit to prevent performance impacts.

## 2.7.6 (May, 20 2023)

### Notes

- Release date: **(May, 20 2023)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #306](https://github.com/zscaler/terraform-provider-zpa/pull/306) Fix resource ``zpa_policy_forwarding_rule`` to ensure updates are executed during resource rule modifications.
- [PR #307](https://github.com/zscaler/terraform-provider-zpa/pull/307) Fix resource ``zpa_policy_timeout_rule`` to ensure updates are executed during resource rule modifications.
- [PR #308](https://github.com/zscaler/terraform-provider-zpa/pull/308) Fix the following access rule resources to ensure updates are executed during resource rule modifications:
  * ``zpa_policy_inspection_rule``
  * ``zpa_policy_isolation_rule``

## 2.7.5 (May, 18 2023)

### Notes

- Release date: **(May, 18 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #304](https://github.com/zscaler/terraform-provider-zpa/pull/304) Fix attribute ``select_connector_close_to_app`` by setting schema attribute to ``ForceNew`` across all application segments to ensure proper resource update when UDP port is set and ``select``_connector_close_to_app`` is switched to false.

## 2.7.4 (May, 13 2023)

### Notes

- Release date: **(May, 13 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #301](https://github.com/zscaler/terraform-provider-zpa/pull/301) Improve scim values searching

## 2.7.3 (May, 11 2023)

### Notes

- Release date: **(May, 11 2023)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #298](https://github.com/zscaler/terraform-provider-zpa/pull/285) Fixed issue with empty IDs in the resource ``zpa_service_edge_groups``
- [PR #298](https://github.com/zscaler/terraform-provider-zpa/pull/285) Fix Service Edge Group Trusted Networks  for resource ``zpa_service_edge_groups``

## 2.7.2 (April, 28 2023)

### Notes

- Release date: **(April, 28 2023)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #285](https://github.com/zscaler/terraform-provider-zpa/pull/285) Allow empty server group attribute in ``server_group`` attribute within an application segment
- [PR #291](https://github.com/zscaler/terraform-provider-zpa/pull/291) Added function to support detaching objects from all policy types prior to destroy operation.

### Enhacements

- [PR #292](https://github.com/zscaler/terraform-provider-zpa/pull/292) Added validation to application segments on attributes ``select_closest_app_connector`` to ensure no UDP port configuration is submitted. By default only TCP ports are supported when this attribute is set to ``true``.

## 2.7.1 (April, 12 2023)

### Notes

- Release date: **(April, 12 2023)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #279](https://github.com/zscaler/terraform-provider-zpa/pull/279) Update to Zscaler-SDK-GO 1.4.0 to support long Terraform runs and improve exponential backoff mechanism.
- [PR #280](https://github.com/zscaler/terraform-provider-zpa/pull/280) Added function to support detaching objects from Access policies prior to destroy operation.
- [PR #281](https://github.com/zscaler/terraform-provider-zpa/pull/281) Fixed browser access acceptance test to prevent port overlap and lingering resources
- [PR #285](https://github.com/zscaler/terraform-provider-zpa/pull/285) Make ``server_group`` attribute in the application segment optional to support UI behavior Issue[#283](https://github.com/zscaler/terraform-provider-zpa/issues/283)

⚠️ **WARNING:** In order to improve performance during long long Terraform runs involving ``zpa_application_segment`` resource, the Provider no longer performs pre-check on port overlaps. For this reason, we advise that Terraform configuration is checked properly during coding to ensure application segments with the same domain and ports are not conflicting. The port overlap pre-check remains in place for all other application segment types.

## 2.7.0 (March, 23 2023)

### Notes

- Release date: **(March, 23 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #272](https://github.com/zscaler/terraform-provider-zpa/pull/272) The ZPA Terraform Provider API Client, will now support long runs, that exceeds the 3600 seconds token validity. Terraform will automatically request a new API bearer token at that time in order to continue the resource provisioning. This enhacement will prevent long pipeline runs from being interrupted.

- [PR #272](https://github.com/zscaler/terraform-provider-zpa/pull/272) Update provider to Zscaler-SDK-GO v1.3.0

- [PR #272](https://github.com/zscaler/terraform-provider-zpa/pull/272) The SDK now supports authentication to ZPA DEV environment.

### Bug Fix

- [PR #271](https://github.com/zscaler/terraform-provider-zpa/pull/271) Added deprecate message to ``zpa_segment_group`` under the following attributes:
  - ``policy_migrated``: "The `policy_migrated` field is now deprecated for the resource `zpa_segment_group`, please remove this attribute to prevent configuration drifts"
  - ``tcp_keep_alive_enabled``: "The `tcp_keep_alive_enabled` field is now deprecated for the resource `zpa_segment_group`, please replace all uses of this within the `zpa_application_segment`resources with the attribute `tcp_keep_alive`".

  Both the above attributes can be safely removed without impact to production configuration; however, they are still supported for backwards compatibity purposes. [#270](https://github.com/zscaler/terraform-provider-zpa/issues/270)

## 2.6.6 (March, 20 2023)

### Notes

- Release date: **(March, 20 2023)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #268](https://github.com/zscaler/terraform-provider-zpa/pull/268) Fixed provider crashing when flattening IDP controller user metadata function Issue [#267](https://github.com/zscaler/terraform-provider-zpa/issues/267)

- [PR #268](https://github.com/zscaler/terraform-provider-zpa/pull/268) Added new ZPA IDP Controller attributes to data source. The following new attributes have been added:
  - ``login_hint``
  - ``force_auth``
  - ``enable_arbitrary_auth_domains``

## 2.6.5 (March, 19 2023)

### Notes

- Release date: **(March, 19 2023)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #262](https://github.com/zscaler/terraform-provider-zpa/pull/262) SCIM Group Search Pagination Issue affecting the following resource:
  - ``zpa_scim_groups``

## 2.6.4 (March, 16 2023)

### Notes

- Release date: **(March, 16 2023)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #263](https://github.com/zscaler/terraform-provider-zpa/pull/263) (fix) Added missing new object_type ``PLATFORM`` validation for access policy resources

## 2.6.3 (March, 7 2023)

### Notes

- Release date: **(March, 7 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #257](https://github.com/zscaler/terraform-provider-zpa/pull/257) Added the new ZPA Application Segment attributes for the following resources:
  - ``zpa_application_segment``, ``zpa_application_segment_browser_access``, ``zpa_application_segment_inspection``, ``zpa_application_segment_pra``
    - ``tcp_keep_alive``
    - ``is_incomplete_dr_config``
    - ``use_in_dr_mode``
    - ``select_connector_close_to_app``

  - ``zpa_app_connector_group``
    - ``use_in_dr_mode``

## 2.6.2 (March, 1 2023)

### Notes

- Release date: **(March, 1 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #251](https://github.com/zscaler/terraform-provider-zpa/pull/251) - Added new action ``REQUIRE_APPROVAL`` to ``zpa_policy_access_rule`` - [Issue [#250](https://github.com/zscaler/terraform-provider-zpa/issues/250)]

## 2.6.1 (February, 15 2023)

### Notes

- Release date: **(February, 15 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #242](https://github.com/zscaler/terraform-provider-zpa/pull/242) - Added new data source and resources below:
  - ``zpa_isolation_profile`` - This data source gets all isolation profiles for the specified customer. The Isolation Profile ID can then be referenced in a ``zpa_policy_isolation_rule`` when the ``action`` is set to ``ISOLATE``
  - ``zpa_policy_isolation_rule`` - This resource, creates an Isolation Rule. Notice that in order to create an isolation policy the ZPA tenant must be licensed accordingly. ``zpa_policy_isolation_rule`` when the ``action`` is set to ``ISOLATE``

### Bug Fix

- [PR #244](https://github.com/zscaler/terraform-provider-zpa/pull/244) - Fixed ``zpa_server_groups`` resource ``servers`` attribute to support typeSet instead of typeList.
- [PR #244](https://github.com/zscaler/terraform-provider-zpa/pull/244) - Fixed ``zpa_app_connector_group`` resource ``connectors`` attribute to support typeSet instead of typeList.

## 2.6.0 (February, 15 2023)

### Notes

- Release date: **(February, 15 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #242](https://github.com/zscaler/terraform-provider-zpa/pull/242) - Added new data source and resources below:
  - ``zpa_isolation_profile`` - This data source gets all isolation profiles for the specified customer. The Isolation Profile ID can then be referenced in a ``zpa_policy_isolation_rule`` when the ``action`` is set to ``ISOLATE``
  - ``zpa_policy_isolation_rule`` - This resource, creates an Isolation Rule. Notice that in order to create an isolation policy the ZPA tenant must be licensed accordingly. ``zpa_policy_isolation_rule`` when the ``action`` is set to ``ISOLATE``

### Bug Fix

- [PR #244](https://github.com/zscaler/terraform-provider-zpa/pull/244) - Fixed ``zpa_server_groups`` resource ``servers`` attribute to support typeSet instead of typeList.

## 2.5.6 (January, 24 2023)

### Notes

- Release date: **(January, 24 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #238](https://github.com/zscaler/terraform-provider-zpa/pull/238) - Added new log_type (``zpn_pbroker_comprehensive_stats``) attribute to ``zpa_lss_config_log_type_formats`` and ``zpa_lss_config_controller``.

### Notes

- Release date: **(January, 16 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #232](https://github.com/zscaler/terraform-provider-zpa/pull/232) - Added new ZPA Inspection control parameters

  - ZPA Inspection Profile: ``web_socket_controls``
  - ZPA Custom Inspection Control:
    - ``control_type``: The following values are supported:
      - ``WEBSOCKET_PREDEFINED``, ``WEBSOCKET_CUSTOM``, ``ZSCALER``, ``CUSTOM``, ``PREDEFINED``

    - ``protocol_type``: The following values are supported:
      - ``HTTP``, ``WEBSOCKET_CUSTOM``, ``ZSCALER``, ``CUSTOM``, ``PREDEFINED``

### Fixes

- [PR #234](https://github.com/zscaler/terraform-provider-zpa/pull/234) - Removed Segment Group detachment function, so it can use the new ``force_delete`` parameter when removing application segments from a segment group.

## 2.5.3 (January, 2 2023)

### Notes

- Release date: **(January, 2 2023)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #224](https://github.com/zscaler/terraform-provider-zpa/pull/224) Implemented longitude/latitude math function validation for more accurancy when configuring ``zpa_app_connector_group`` resources.


# 2.5.2 (December, 02 2022)

### Notes

- Release date: **(December, 02 2022)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #223](https://github.com/zscaler/zscaler-sdk-go/pull/223) Fixed pagination issue with ZPA endpoints

## 2.5.1 (November, 30 2022)

### Notes

- Release date: **(November, 30 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #219](https://github.com/zscaler/terraform-provider-zia/pull/219) Added ForceNew helper to ``zpa_policy_timeout_rule`` parameters ``reauth_idle_timeout`` and ``reauth_timeout``. Changing the values will cause the resource to be recreated on the fly.
- [PR #219](https://github.com/zscaler/terraform-provider-zia/pull/219) Added missing ``ip_anchored`` parameter to ``resource_zpa_application_segment_browser_access``
- [PR #220](https://github.com/zscaler/terraform-provider-zia/pull/220) Udated provider to Zscaler-SDK-Go v0.3.2 to ensure pagination works correctly when more than 500 items on a list.

## 2.5.0 (November, 27 2022)

### Notes

- Release date: **(November, 27 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #217](https://github.com/zscaler/terraform-provider-zia/pull/217) Fixed Read/Update/Delete functions to allow automatic recreation of resources, that have been manually deleted via the UI.
- [PR #217](https://github.com/zscaler/terraform-provider-zia/pull/217) Updated provider to zscaler-sdk-go v0.2.2

## 2.4.1 (November, 9 2022)

### Notes

- Release date: **(November, 9 2022)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #208](https://github.com/zscaler/terraform-provider-zpa/pull/208) - Implemented TCP/UDP Port overlap check and duplicated domain validation for ``zpa_application_segment_browser_access``
- [PR #209](https://github.com/zscaler/terraform-provider-zpa/pull/209) - Implemented TCP/UDP Port overlap check and duplicated domain validation for ``zpa_application_segment_pra``.
- [PR #210](https://github.com/zscaler/terraform-provider-zpa/pull/210) - Implemented TCP/UDP Port overlap check and duplicated domain validation for ``zpa_application_segment_inspection``.

### Bug Fixes

- [PR #206](https://github.com/zscaler/terraform-provider-zpa/pull/206) - Fix TCP/UDP port overlap check issue
- [PR #207](https://github.com/zscaler/terraform-provider-zpa/pull/207) - Fix duplicated domain_name entries during TCP/UDP port overlap issues

## 2.4.0 (October, 24 2022)

### Notes

- Release date: **(October, 24 2022)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #188](https://github.com/zscaler/terraform-provider-zpa/pull/188) - feat(new parameters added to App Connector Group resource TCPQuick*
  - The following new App Connector Group parameters have been added:
  - tcpQuickAckApp - Whether TCP Quick Acknowledgement is enabled or disabled for the application.
  - tcpQuickAckAssistant - Whether TCP Quick Acknowledgement is enabled or disabled for the application.
  - tcpQuickAckReadAssistant - Whether TCP Quick Acknowledgement is enabled or disabled for the application.
  - UseInDrMode
- [PR #188](https://github.com/zscaler/terraform-provider-zpa/pull/188) - Upgrade to zscaler-sdk-go v0.0.12 to support new App Connector Group parameters ``TCPQuick*`` and ``UseInDrMode``
- [PR #190](https://github.com/zscaler/terraform-provider-zpa/pull/190) - Added ZPA Terraform Provider Video Series link in the documentation, leading to [Zenith Community Portal](https://community.zscaler.com/tag/devops)
- [PR #194](https://github.com/zscaler/terraform-provider-zpa/pull/194) - Updated Provider to Zscaler-SDK-GO v0.1.1
- [PR #196](https://github.com/zscaler/terraform-provider-zpa/pull/196) - Renamed ``zpa_browser_access`` resource and data source to ``zpa_application_segment_browser_access`` for better distinction with other application segment resources. The use of the previous resource name is still supported; however, a warning message will be displayed after the apply process to inform about the change.
- [PR #196](https://github.com/zscaler/terraform-provider-zpa/pull/196) - Fixed ``zpa_application_segment_browser_access`` ``clientless_apps`` inner parameters, which were not being updated during PUT method.
- [PR #197](https://github.com/zscaler/terraform-provider-zpa/pull/197) - Updated ``zpa_service_edge_group`` parameter ``is_public`` to accept a value of Bool (true or false) instead of the current String values of (DEFAULT, TRUE or FALSE) for easier configuration. The Provider will convert the input value to string during run-time.
- [PR #201](https://github.com/zscaler/terraform-provider-zpa/pull/201) - Added ``zpa_app_connector_controller`` resource to allow app connector resource management and bulk delete action for app connector deproviosioning.
- [PR #202](https://github.com/zscaler/terraform-provider-zpa/pull/202) - Included validation function in the ``zpa_app_connector_group`` resource for the parameters ``version_profile_name`` and ``version_profile_id``. Users can now use ``version_profile_name`` with one of the following values: ``Default``, ``Previous Default``, ``New Release``

### Bug Fixes

- [PR #181](https://github.com/zscaler/terraform-provider-zpa/pull/181) - Added Support to ZPA Preview Cloud and updated to zscaler-sdk-go v0.0.9
- [PR #193](https://github.com/zscaler/terraform-provider-zpa/pull/193) - Fixed rule order in access policies, when Zscaler Deception rule exists.
- [PR #198](https://github.com/zscaler/terraform-provider-zpa/pull/198) - Due to Golang update the function ``ConfigureFunc`` used to configure the provider was deprecated; hence, the ZPA Terraform Provider was updated to use the ``ConfigureContextFunc`` instead.
- [PR #199](https://github.com/zscaler/terraform-provider-zpa/pull/199) - Fix application segment tcp/udp port conflict. The provider will issue an error message when 2 application segments have conflicting domain_name, tcp/udp ports
- [PR #200](https://github.com/zscaler/terraform-provider-zpa/pull/200) - Implemented new application segment parameter ``force_delete`` to ensure dependency removal prior to delete action.

## 2.3.2 (September, 2 2022)

### Notes

- Release date: **(September, 2 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

Fixed authentication issue when specifying zpa_cloud="PRODUCTION"
## 2.3.1 (August 17 2022)

### Notes

- Release date: **(August, 17 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #169](https://github.com/zscaler/terraform-provider-zpa/pull/169) Fixed policy rule order, where the rule order in the UI didn't correspond to the desired order set in HCL. Issue [[#166](https://github.com/zscaler/terraform-provider-zpa/issues/166)]
- [PR #170](https://github.com/zscaler/terraform-provider-zpa/pull/170) Fixed special character encoding, where certain symbols caused Terraform to indicate potential configuration drifts. Issue [[#149](https://github.com/zscaler/terraform-provider-zpa/issues/149)]
- [PR #171](https://github.com/zscaler/terraform-provider-zpa/pull/171) Fixed policy configuration attributes where i.e SCIM_GROUPs were causing drifts without changes have been performed. Issue [[#165](https://github.com/zscaler/terraform-provider-zpa/issues/165)]
- [PR #175](https://github.com/zscaler/terraform-provider-zpa/pull/175) Fixed application segment drifts caused by tcp & udp ports.
- [PR #176](https://github.com/zscaler/terraform-provider-zpa/pull/176) Fixed application segment PRA drifts caused by tcp & udp ports.
- [PR #177](https://github.com/zscaler/terraform-provider-zpa/pull/177) Fixed application segment Inspection drifts caused by tcp & udp ports.

## 2.3.0 (August, 17 2022)

### Notes

- Release date: **(August, 17 2022)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #161](https://github.com/zscaler/terraform-provider-zpa/pull/161) Integrated newly created Zscaler GO SDK. Models are now centralized in the repository [zscaler-sdk-go](https://github.com/zscaler/zscaler-sdk-go)

## 2.2.2 (July, 19 2022)

### Notes

- Release date: **(July, 19 2022)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #159](https://github.com/zscaler/terraform-provider-zpa/pull/159) Added Terraform UserAgent for Backend API tracking

## 2.2.1 (July, 6 2022)

### Notes

- Release date: **(July, 6 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- Fix: Fixed authentication mechanism variables for ZPA Beta and GOV

### Documentation

1. Fixed application segment documentation and examples

## 2.2.0 (June 30, 2022)

### Notes:
- Release date: **(June, 30 2022)**
- Supported Terraform version: **v1.x**

### New Features

1. The provider now supports the following ZPA Privileged Remote Access (PRA) features:

- **zpa_application_segment_pra** - The resource supports enabling Priviledged Remote Access Application Segment ``SECURE_REMOTE_ACCESS``option for `RDP` and `SSH` via the ``app_types`` parameter. PR [#133](https://github.com/zscaler/terraform-provider-zpa/pull/133)

2. The provider now supports the following ZPA Inspection features:
- **zpa_inspection_custom_controls** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
- **zpa_inpection_predefined_controls** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
- **zpa_inspection_all_predefined_controls** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
- **zpa_inspection_profile** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
- **zpa_policy_access_inspection_rule** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
- **zpa_application_segment_inspection** - The resource supports enabling `INSPECT` for `HTTP` and `HTTPS` via the `app_types` parameter. PR [#135](https://github.com/zscaler/terraform-provider-zpa/pull/135)

4. Implemented a new Application Segment resource parameter ``select_connector_close_to_app``. The parameter can only be set for TCP based applications. PR [#137](https://github.com/zscaler/terraform-provider-zpa/pull/137)

### Enhancements

- Added support to `scim_attribute_header` to support policy access SCIM criteria based on SCIM attribute values.  Issue [#146](https://github.com/zscaler/terraform-provider-zpa/issues/146) / PR [#147]((https://github.com/zscaler/terraform-provider-zpa/pull/147))

- ZPA Beta Cloud: The provider now supports authentication via environment variables or static credentials to ZPA Beta Cloud. For authentication instructions please refer to the documentation page [here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.md) PR [#136](https://github.com/zscaler/terraform-provider-zpa/pull/136)

- ZPA Gov Cloud: The provider now supports authentication via environment variables or static credentials to ZPA Gov Cloud. For authentication instructions please refer to the documentation page [here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.md) PR [#145](https://github.com/zscaler/terraform-provider-zpa/pull/145)

### Bug Fixes
- Fix: Fixed update function on **zpa_app_server_controller** resource to ensure desired state is enforced in the upstream resource. Issue [#128](https://github.com/zscaler/terraform-provider-zpa/issues/128)
- Fix: Fixed `enabled` parameter on **zpa_app_connector_group** resource by removing default action from resource schema. Issue [#128](https://github.com/zscaler/terraform-provider-zpa/issues/128)
- Fix: Fixed Golangci linter and upgraded to golangci-lint-action@v3

### Documentation
1. Added release notes guide to documentation PR [#140](https://github.com/zscaler/terraform-provider-zpa/pull/140)
2. Fixed documentation misspellings

## 2.1.5 (May, 18 2022)

### Notes:
- Release date: **(May, 18 2022)**
- Supported Terraform version: **v1.x**

### Annoucements:

The Terraform Provider for Zscaler Private Access (ZPA) is now officially hosted under Zscaler's GitHub account and published in the Terraform Registry. For more details, visit the Zscaler Community Article [Here](https://community.zscaler.com/t/zpa-and-zia-terraform-providers-now-verified/16675)
Administrators who used previous versions of the provider, and followed instructions to install the binary as a custom provider, must update their provider block as such:

```terraform
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "2.1.5"
    }
  }
}
provider "zpa" {}

```

### Enhancements:
- Documentation: Updated documentation to comply with Terraform registry formatting. #125
- ``zpa_posture_profile`` Updated search mechanism to support posture profile name search without the Zscaler cloud name. PR #123
- ``zpa_trusted_network`` Updated search mechanism to support trusted network name search without the Zscaler cloud name. PR #123

### Bug Fixes

- Fixed ``zpa_application_segment`` to support updates on ``tcp_port_ranges``, ``udp_port_ranges`` and ``tcp_port_range``, ``udp_port_range`` Issue #103

## 2.1.3 (May, 18 2022)

### Notes:

- Release date: **(May, 18 2022)**
- Supported Terraform version: **v1.x**

### Annoucements:

The Terraform Provider for Zscaler Private Access (ZPA) is now officially hosted under Zscaler's GitHub account and published in the Terraform Registry. For more details, visit the Zscaler Community Article [Here](https://community.zscaler.com/t/zpa-and-zia-terraform-providers-now-verified/16675)
Administrators who used previous versions of the provider, and followed instructions to install the binary as a custom provider, must update their provider block as such:

```terraform
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "2.1.3"
    }
  }
}
provider "zpa" {}

```

## 2.1.2 (May, 6 2022)

### Notes:
- Release date: **(May, 6 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- Fix: tcp and udp ports were not being updated during changes, requiring the application segment resource to be fully destroyed and rebuilt. Implemented ``ForceNew`` in the the ``zpa_application_segment`` resource parameters: ``tcp_port_range``, ``udp_port_range``, ``tcp_port_ranges``, ``udp_port_ranges``. This behavior instructs Terraform to first destroy and then recreate the resource if any of the attributes change in the configuration, as opposed to trying to update the existing resource. The destruction of the resource does not impact attached resources such as server groups, segment groups or policies.

## 2.1.1 (April, 27 2022)

### Notes:
- Release date: **(April, 27 2022)**
- Supported Terraform version: **v1.x**

### Enhancements:

1. Refactored and added new acceptance tests for better statement coverage. These tests are considered best practice and were added to routinely verify that the ZPA Terraform Plugin produces the expected outcome. [PR#88], [PR#96], [PR#98], [PR#99]

2. Support explicitly empty port ranges. Allow optional use of Attributes as Blocks syntax for ``zpa_application_segment`` {tcp,udp}_port_range blocks, allowing clean specification of "no port ranges" in dynamic contexts. [PR#97](https://github.com/zscaler/terraform-provider-zpa/pull/97) Thanks @isometry

### Deprecations

1. Deprecated all legacy policy set controller endpoints: ``/policySet/global``, ``/policySet/reauth``, ``/policySet/bypass`` [PR#88](https://github.com/zscaler/terraform-provider-zpa/pull/88)

2. Deprecated all references to ZPA private API gateway. [PR#87](https://github.com/zscaler/terraform-provider-zpa/pull/87)

## 2.1.0 (March, 05 2022)

### Notes

- Release date: **(March, 05 2022)**
- Supported Terraform version: **v1.x**

### Enhancements:

1. Refactored and added new acceptance tests. These tests are considered best practice and were added to routinely verify that the ZPA Terraform Plugin produces the expected outcome. [PR#xx](https://github.com/zscaler/terraform-provider-zpa/pull/xx)

- ``data_source_zpa_app_connector_controller_test``
- ``data_source_zpa_app_connector_group_test``
- ``data_source_zpa_app_server_controller_test``
- ``data_source_zpa_application_segment_test``
- ``data_source_zpa_ba_certificate_test``
- ``data_source_zpa_browser_access_test``
- ``data_source_zpa_cloud_connector_group_test``
- ``data_source_zpa_customer_version_profile_test``
- ``data_source_zpa_enrollement_cert_test``
- ``data_source_zpa_idp_controller_test``
- ``data_source_zpa_lss_config_client_types_test``
- ``data_source_zpa_lss_config_log_types_formats_test``
- ``data_source_zpa_lss_config_status_codes_test``
- ``data_source_zpa_machine_group_test``
- ``data_source_zpa_posture_profile_test``
- ``data_source_zpa_segment_group_test``
- ``data_source_zpa_server_group_test``
- ``data_source_zpa_trusted_network_test``
- ``resource_zpa_app_connector_group_test``
- ``resource_zpa_app_server_controller_test``
- ``resource_zpa_application_segment_test``
- ``resource_zpa_segment_group_test``
- ``resource_zpa_server_group_test``
- ``resource_zpa_service_edge_group_test``
- ``resource_zpa_policy_access_rule_test``
- ``resource_zpa_policy_access_timeout_rule_test``
- ``resource_zpa_policy_access_forwarding_rule_test``

### Bug Fixes

- Fix: Acceptance Tests for ``zpa_browser_access_test``
- Fix: Consolidate Policy Type resources
- Fix: Refactor ZPA API Client

## 2.0.7 (February, 17 2022)

### Notes:
- Release date: **(February, 17 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- ZPA-50: Fixed and removed deprecated arguments from ``zpa_application_segments`` data source and resource :wrench:
- ZPA-50: Fixed ``zpa_posture_profile`` and ``zpa_trusted_networks`` acceptance tests to include ZIA cloud name :wrench:

### Enhancements:

- ZPA-51: Updated common ``NetworkPorts`` flatten and expand functions for better optimization and global use across multiple application segment resources. This update affects the following resources: ``data_source_zpa_application_segment``, ``data_source_zpa_browser_access`` and ``resource_zpa_application_segment``, ``resource_source_zpa_browser_access`` :rocket:

## 2.0.6 (February, 3 2022)

### Notes:
- Release date: **(February, 3 2022)**
- Supported Terraform version: **v1.x**

### New Data Sources:
- Added new data source for ``zpa_app_connector_controller`` resource. [PR#62](https://github.com/zscaler/terraform-provider-zpa/pull/62)
- Added new data source for ``zpa_service_edge_controller`` resource. [PR#63](https://github.com/zscaler/terraform-provider-zpa/pull/63)

### New Acceptance Tests:
These tests are considered best practice and were added to routinely verify that the ZPA Terraform Plugin produces the expected outcome. [PR#64](https://github.com/zscaler/terraform-provider-zpa/pull/64)

- ``data_source_zpa_app_connector_controller_test``
- ``data_source_zpa_app_connector_group_test``
- ``data_source_zpa_app_server_controller_test``
- ``data_source_zpa_application_segment_test``
- ``data_source_zpa_ba_certificate_test``
- ``data_source_zpa_browser_access_test``
- ``data_source_zpa_cloud_connector_group_test``
- ``data_source_zpa_customer_version_profile_test``
- ``data_source_zpa_enrollement_cert_test``
- ``data_source_zpa_idp_controller_test``
- ``data_source_zpa_lss_config_client_types_test``
- ``data_source_zpa_lss_config_log_types_formats_test``
- ``data_source_zpa_lss_config_status_codes_test``
- ``data_source_zpa_machine_group_test``
- ``data_source_zpa_posture_profile_test``
- ``data_source_zpa_segment_group_test``
- ``data_source_zpa_server_group_test``
- ``data_source_zpa_trusted_network_test``
- ``resource_zpa_app_connector_group_test``
- ``resource_zpa_app_server_controller_test``
- ``resource_zpa_application_segment_test``
- ``resource_zpa_segment_group_test``
- ``resource_zpa_server_group_test``
- ``resource_zpa_service_edge_group_test``

## 2.0.5 (December, 20 2021)

### Notes:
- Release date: **(December, 20 2021)**
- Supported Terraform version: **v1.x**

### Enhancements

- The provider now supports the ability to import policy access resources via its `name` and/or `id` property to support easier migration of existing ZPA resources via `terraform import` command.
- The  following policy access resources are supported:
  - resource_zpa_policy_access_rule - [PR#51](https://github.com/zscaler/terraform-provider-zpa/issues/51)] :rocket:
  - resource_zpa_policy_access_timeout_rule - [PR#51](https://github.com/zscaler/terraform-provider-zpa/pull/51) :rocket:
  - resource_zpa_policy_access_forwarding_rule - [PR#51](https://github.com/zscaler/terraform-provider-zpa/pull/51) :rocket:

- The provider now supports policy access creation to be associated with Cloud Connector Group resource
  - resource_zpa_policy_access_rule - [PR#54](https://github.com/zscaler/terraform-provider-zpa/pull/54) :rocket:
  - Added new `client_type` to support access, forward, and timeout policy creation. The following new types have been added:
  - zpn_client_type_ip_anchoring, zpn_client_type_browser_isolation, zpn_client_type_machine_tunnel and zpn_client_type_edge_connector. [PR#57](https://github.com/zscaler/terraform-provider-zpa/issues/57)] :rocket:

- Updated the following examples for more accuracy:
  - resource_zpa_policy_access_rule
  - resource_zpa_app_connector_group

### Bug Fixes:
- Fixed pagination issues with all resources where only the default pagesize was being returned. [PR#52](https://github.com/zscaler/terraform-provider-zpa/pull/52) :wrench:
- Fixed issue where Terraform showed that resources had been modified even though nothing had been changed in the upstream resources.[PR#54](https://github.com/zscaler/terraform-provider-zpa/pull/54) :wrench:

## 2.0.4 (December, 6 2021)

### Notes:
- Release date: **(December, 6 2021)**
- Supported Terraform version: **v1.x**

### New Data Sources:
- Added new data source for ``zpa_browser_access`` resource.

### Enhancements:
- The provider now supports the ability to import resources via its `name` and/or `id` property to support easier migration of existing ZPA resources via `terraform import` command.
This capability is currently available to the following resources:
- resource_zpa_app_connector_group - Issue ([#29](https://github.com/zscaler/terraform-provider-zpa/issues/29))
- resource_zpa_app_server_controller - [PR#42](https://github.com/zscaler/terraform-provider-zpa/pull/42) :rocket:
- resource_zpa_application_segment - [PR#42](https://github.com/zscaler/terraform-provider-zpa/pull/42) :rocket:
- resource_zpa_segment_group - [PR#42](https://github.com/zscaler/terraform-provider-zpa/pull/42) :rocket:
- resource_zpa_server_group - [PR#42](https://github.com/zscaler/terraform-provider-zpa/pull/42) :rocket:
- resource_zpa_service_edge_group - [PR#42](https://github.com/zscaler/terraform-provider-zpa/pull/42) :rocket:
- resource_zpa_provisioning_key - [PR#45](https://github.com/zscaler/terraform-provider-zpa/pull/45) :rocket:
- resource_zpa_browser_access - [PR#48](https://github.com/zscaler/terraform-provider-zpa/pull/48) :rocket:
- zpa_lss_config_controller - [PR#48](https://github.com/zscaler/terraform-provider-zpa/pull/48) :rocket:

Note: To import resources not currently supported, the resource numeric ID is required.

### Bug Fixes:
- Fixed [INFO] and [Error] message in ``data_source_zpa_lss_config_controller`` [PR#43](https://github.com/zscaler/terraform-provider-zpa/pull/43) 🔧

## 2.0.3 (November, 21 2021)

### Notes:
- Release date: **(November, 21 2021)**
- Supported Terraform version: **v1.x**

###  Dependabot Updates:
- Dependabot updates [PR#33](https://github.com/zscaler/terraform-provider-zpa/pull/33/) Bump github.com/hashicorp/terraform-plugin-docs from 0.5.0 to 0.5.1 #33
- Dependabot updates [PR#34](https://github.com/zscaler/terraform-provider-zpa/pull/34) Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.8.0 to 2.9.0

## 2.0.2 (November, 7 2021)

### Notes:
- Release date: **(November, 7 2021)**
- Supported Terraform version: **v1.x**

### Enhancements:
- Added custom validation function ``ValidateStringFloatBetween`` to ``resource_zpa_app_connector_group`` to validate ``longitude`` and ``latitude`` parameters. [ZPA-17](https://github.com/zscaler/terraform-provider-zpa/pull/17).
- Added custom validation function ``ValidateStringFloatBetween`` to ``resource_zpa_service_edge_group`` to validate ``longitude`` and ``latitude`` parameters. [ZPA-18](https://github.com/zscaler/terraform-provider-zpa/pull/18).

## 2.0.1 (November, 4 2021)

### Notes:
- Release date: **(November, 4 2021)**
- Supported Terraform version: **v1.x**

### Bug Fixes:
- Fixed issue where provider authentication parameters for hard coded credentials was not working:
- Changed the following variable names: ``client_id``, ``client_secret`` and ``customerid`` to ``zpa_client_id``, ``zpa_client_secret`` and ``zpa_customer_id``.

## 2.0.0 (November, 3 2021)

### Notes:
- Release date: **(November, 3 2021)**
- Supported Terraform version: **v1.x**

- New management APIs are now available to manage App Connectors, App Connector Groups, Service Edges, Service Edge Groups, and Log Streaming Service (LSS) configurations.
- New prerequisite APIs for enrollment certificates, provisioning keys, and to get version profiles, client types, status codes, and LSS formats are added.
- A new API to reorder policy rules is added.
- The endpoints to get all browser access (BA) certificates, IdPs, posture profiles, trusted networks, and SAML attributes are now deprecated, and new APIs with pagination are provided.
- API endpoints specific to a policy (global/reauth/bypass) are deprecated and replaced by a generic API that takes policyType as a parameter.
- The port range configuration for the application segment has been enhanced for more readability. The tcpPortRanges and udpPortRanges fields are deprecated and replaced with tcpPortRange and udpPortRange.

### Features:
### New Resources:
- New Resource: ``resource_zpa_app_connector_group`` 🆕
- New Resource: ``resource_zpa_service_edge_group`` 🆕
- New Resource: ``resource_zpa_provisioning_key`` 🆕
- New Resource: ``resource_zpa_lss_config_controller`` 🆕

### New Data Sources:
- New Data Source: ``data_source_zpa_enrollement_cert`` 🆕
- New Data Source: ``data_source_zpa_customer_version_profile`` 🆕
- New Data Source: ``data_source_zpa_lss_config_controller`` 🆕
- New Data Source: ``data_source_zpa_lss_config_log_types_formats`` 🆕
- New Data Source: ``data_source_zpa_lss_config_status_codes`` 🆕
- New Data Source: ``data_source_zpa_lss_config_client_types`` 🆕
- New Data Source: ``data_source_zpa_policy_type`` 🆕

### Enhacements:
1. A new API to reorder policy rules is added. This update affects the following resources:
    - ``resource_zpa_policy_access_rule`` :rocket:
    - ``resource_zpa_policy_access_timeout_rule`` :rocket:
    - ``resource_zpa_policy_access_forwarding_rule`` :rocket:
2. Updated the following data sources to V2 API to support pagination:
    - ``data_source_zpa_idp_controller`` :rocket:
    - ``data_source_zpa_saml_attribute``:rocket:
    - ``data_source_zpa_scim_attribute_header`` :rocket:
    - ``data_source_zpa_trusted_network`` :rocket:
    - ``data_source_zpa_posture_profile`` :rocket:
    - ``data_source_zpa_ba_certificate`` :rocket:
    - ``data_source_zpa_machine_group`` :rocket:
3. Added additional validations to ``bypass_type`` parameter in ``resource_zpa_browser_access``. :rocket:
4. The port range configuration for the application segment has been enhanced for more readability. This update affects the following resources:
    - ``resource_zpa_application_segment`` :rocket:
    - ``resource_zpa_browser_access`` :rocket:

### Deprecations:
- API endpoints specific to a policy (global/reauth/bypass) are deprecated and replaced by a generic API that takes policyType as a parameter.

1. Deprecated ``data_source_zpa_global_forwarding_policy`` and ``data_source_zpa_global_timeout_policy`` and replaced with ``data_source_zpa_policy_type`` 💥

2. Deprecated ``data_source_zpa_global_access_policy`` and renamed with ``data_source_zpa_policy_type`` 💥

3. Deprecated ``tcp_port_ranges`` and ``udp_port_ranges`` fields are deprecated and replaced with ``tcp_port_range`` and ``udp_port_range``. The values will be kept in Terraform schema until next provider update for backwards compatibility. 💥

## 1.0.0 (September, 23 2021)

### Notes:
- Release date: **(September, 23 2021)**
- Supported Terraform version: **v1.x**

### Initial Release
#### New Resources:
- New Resource: ``resource_zpa_app_server_controller`` 🆕
- New Resource: ``resource_zpa_application_segment`` 🆕
- New Resource: ``resource_zpa_browser_access`` 🆕
- New Resource: ``resource_zpa_policy_access_forwarding_rule`` 🆕
- New Resource: ``resource_zpa_policy_access_rule`` 🆕
- New Resource: ``resource_zpa_policy_access_timeout_rule`` 🆕
- New Resource: ``resource_zpa_segment_group`` 🆕
- New Resource: ``resource_zpa_server_group`` 🆕

#### New Data Sources:
- New Data Source: ``data_source_zpa_app_connector_group`` 🆕
- New Data Source: ``data_source_zpa_app_server_controller`` 🆕
- New Data Source: ``data_source_zpa_application_segment`` 🆕
- New Data Source: ``data_source_zpa_ba_certificate`` 🆕
- New Data Source: ``data_source_zpa_cloud_connector_group`` 🆕
- New Data Source: ``data_source_zpa_global_access_policy`` 🆕
- New Data Source: ``data_source_zpa_global_forwarding_policy`` 🆕
- New Data Source: ``data_source_zpa_global_timeout_policy`` 🆕
- New Data Source: ``data_source_zpa_idp_controller`` 🆕
- New Data Source: ``data_source_zpa_machine_group`` 🆕
- New Data Source: ``data_source_zpa_posture_profile`` 🆕
- New Data Source: ``data_source_zpa_saml_attribute`` 🆕
- New Data Source: ``data_source_zpa_scim_attribute_header`` 🆕
- New Data Source: ``data_source_zpa_scim_group`` 🆕
- New Data Source: ``data_source_zpa_segment_group`` 🆕
- New Data Source: ``data_source_zpa_server_group`` 🆕
- New Data Source: ``data_source_zpa_trusted_network`` 🆕
