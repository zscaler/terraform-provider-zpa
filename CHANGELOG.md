# Changelog

## 2.0.7 (February 17, 2022)

## BUG Fixes

- ZPA-50: Fixed and removed deprecated arguments from ``zpa_application_segments`` data source and resource [PR#70](https://github.com/willguibr/terraform-provider-zpa/pull/70) :wrench:
- ZPA-50: Fixed ``zpa_posture_profile`` and ``zpa_trusted_networks`` acceptance tests to include ZIA cloud name [PR#70](https://github.com/willguibr/terraform-provider-zpa/pull/70) :wrench:

## Enhancements

- ZPA-51: Updated common ``NetworkPorts`` flatten and expand functions for better optimization and global use across multiple application segment resources. This update affects the following resources: ``data_source_zpa_application_segment``, ``data_source_zpa_browser_access`` and ``resource_zpa_application_segment``, ``resource_source_zpa_browser_access`` [PR#71](https://github.com/willguibr/terraform-provider-zpa/pull/71) :rocket:
- ZPA-52: Added validation helpers to ``NetworkPortsSchema`` and ``ip_address`` arguments to prevent wrong user data input during apply [PR#72](https://github.com/willguibr/terraform-provider-zpa/pull/72) :rocket

## 2.0.6 (February 3, 2022)

## New Data Sources

- Added new data source for ``zpa_app_connector_controller`` resource. [PR#62](https://github.com/willguibr/terraform-provider-zpa/pull/62)
- Added new data source for ``zpa_service_edge_controller`` resource. [PR#63](https://github.com/willguibr/terraform-provider-zpa/pull/63)

## New Acceptance Tests

These tests are considered best practice and were added to routinely verify that the ZPA Terraform Plugin produces the expected outcome. [PR#64](https://github.com/willguibr/terraform-provider-zpa/pull/64)

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

## 2.0.5 (December 20, 2021)

ENHANCEMENTS:

- The provider now supports the ability to import policy access resources via its `name` and/or `id` property to support easier migration of existing ZPA resources via `terraform import` command.
- The  following policy access resources are supported:
  - resource_zpa_policy_access_rule - [PR#51](https://github.com/willguibr/terraform-provider-zpa/issues/51)] :rocket:
  - resource_zpa_policy_access_timeout_rule - [PR#51](https://github.com/willguibr/terraform-provider-zpa/pull/51) :rocket:
  - resource_zpa_policy_access_forwarding_rule - [PR#51](https://github.com/willguibr/terraform-provider-zpa/pull/51) :rocket:

- The provider now supports policy access creation to be associated with Cloud Connector Group resource
  - resource_zpa_policy_access_rule - [PR#54](https://github.com/willguibr/terraform-provider-zpa/pull/54) :rocket:
  - Added new `client_type` to support access, forward, and timeout policy creation. The following new types have been added:
  - zpn_client_type_ip_anchoring, zpn_client_type_browser_isolation, zpn_client_type_machine_tunnel and zpn_client_type_edge_connector. [PR#57](https://github.com/willguibr/terraform-provider-zpa/issues/57)] :rocket:

- Updated the following examples for more accuracy:
  - resource_zpa_policy_access_rule
  - resource_zpa_app_connector_group

## Bug Fixes

- Fixed pagination issues with all resources where only the default pagesize was being returned. [PR#52](https://github.com/willguibr/terraform-provider-zpa/pull/52) :wrench:
- Fixed issue where Terraform showed that resources had been modified even though nothing had been changed in the upstream resources.[PR#54](https://github.com/willguibr/terraform-provider-zpa/pull/54) :wrench:

## 2.0.4 (December 6, 2021)

## New Data Source

- Added new data source for ``zpa_browser_access`` resource.

## Enhancement

- The provider now supports the ability to import resources via its `name` and/or `id` property to support easier migration of existing ZPA resources via `terraform import` command.
This capability is currently available to the following resources:
  - resource_zpa_app_connector_group - Issue [[#29](https://github.com/willguibr/terraform-provider-zpa/issues/29)]
  - resource_zpa_app_server_controller - [PR#42](https://github.com/willguibr/terraform-provider-zpa/pull/42) :rocket:
  - resource_zpa_application_segment - [PR#42](https://github.com/willguibr/terraform-provider-zpa/pull/42) :rocket:
  - resource_zpa_segment_group - [PR#42](https://github.com/willguibr/terraform-provider-zpa/pull/42) :rocket:
  - resource_zpa_server_group - [PR#42](https://github.com/willguibr/terraform-provider-zpa/pull/42) :rocket:
  - resource_zpa_service_edge_group - [PR#42](https://github.com/willguibr/terraform-provider-zpa/pull/42) :rocket:
  - resource_zpa_provisioning_key - [PR#45](https://github.com/willguibr/terraform-provider-zpa/pull/45) :rocket:
  - resource_zpa_browser_access - [PR#48](https://github.com/willguibr/terraform-provider-zpa/pull/48) :rocket:
  - zpa_lss_config_controller - [PR#48](https://github.com/willguibr/terraform-provider-zpa/pull/48) :rocket:

Note: To import resources not currently supported, the resource numeric ID is required.

BUG FIXES

- Fixed [INFO] and [Error] message in ``data_source_zpa_lss_config_controller`` [PR#43](https://github.com/willguibr/terraform-provider-zpa/pull/43) 🔧

# 2.0.3 (November 21, 2021)

DEPENDABOT UPDATES:

- Dependabot updates [PR#33](https://github.com/willguibr/terraform-provider-zpa/pull/33/) Bump github.com/hashicorp/terraform-plugin-docs from 0.5.0 to 0.5.1 #33
- Dependabot updates [PR#34](https://github.com/willguibr/terraform-provider-zpa/pull/34) Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.8.0 to 2.9.0

## 2.0.2 (November 7, 2021)

ENHANCEMENTS:

- Added custom validation function ``ValidateStringFloatBetween`` to ``resource_zpa_app_connector_group`` to validate ``longitude`` and ``latitude`` parameters. [ZPA-17](https://github.com/willguibr/terraform-provider-zpa/pull/17).
- Added custom validation function ``ValidateStringFloatBetween`` to ``resource_zpa_service_edge_group`` to validate ``longitude`` and ``latitude`` parameters. [ZPA-18](https://github.com/willguibr/terraform-provider-zpa/pull/18).

## 2.0.1 (November 4, 2021)

## Bug Fixes

- Fixed issue where provider authentication parameters for hard coded credentials was not working:
- Changed the following variable names: ``client_id``, ``client_secret`` and ``customerid`` to ``zpa_client_id``, ``zpa_client_secret`` and ``zpa_customer_id``.

# 2.0.0 (November 3, 2021)

## Notes

- New management APIs are now available to manage App Connectors, App Connector Groups, Service Edges, Service Edge Groups, and Log Streaming Service (LSS) configurations.
- New prerequisite APIs for enrollment certificates, provisioning keys, and to get version profiles, client types, status codes, and LSS formats are added.
- A new API to reorder policy rules is added.
- The endpoints to get all browser access (BA) certificates, IdPs, posture profiles, trusted networks, and SAML attributes are now deprecated, and new APIs with pagination are provided.
- API endpoints specific to a policy (global/reauth/bypass) are deprecated and replaced by a generic API that takes policyType as a parameter.
- The port range configuration for the application segment has been enhanced for more readability. The tcpPortRanges and udpPortRanges fields are deprecated and replaced with tcpPortRange and udpPortRange.

### Features

#### New Management Resources

- New Resource: ``resource_zpa_app_connector_group`` 🆕
- New Resource: ``resource_zpa_service_edge_group`` 🆕
- New Resource: ``resource_zpa_provisioning_key`` 🆕
- New Resource: ``resource_zpa_lss_config_controller`` 🆕

#### New Management Data Sources

- New Data Source: ``data_source_zpa_enrollement_cert`` 🆕
- New Data Source: ``data_source_zpa_customer_version_profile`` 🆕
- New Data Source: ``data_source_zpa_lss_config_controller`` 🆕
- New Data Source: ``data_source_zpa_lss_config_log_types_formats`` 🆕
- New Data Source: ``data_source_zpa_lss_config_status_codes`` 🆕
- New Data Source: ``data_source_zpa_lss_config_client_types`` 🆕
- New Data Source: ``data_source_zpa_policy_type`` 🆕

### Enhancements

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

### Deprecations

- API endpoints specific to a policy (global/reauth/bypass) are deprecated and replaced by a generic API that takes policyType as a parameter.

1. Deprecated ``data_source_zpa_global_forwarding_policy`` and ``data_source_zpa_global_timeout_policy`` and replaced with ``data_source_zpa_policy_type`` 💥

2. Deprecated ``data_source_zpa_global_access_policy`` and renamed with ``data_source_zpa_policy_type`` 💥

3. Deprecated ``tcp_port_ranges`` and ``udp_port_ranges`` fields are deprecated and replaced with ``tcp_port_range`` and ``udp_port_range``. The values will be kept in Terraform schema until next provider update for backwards compatibility. 💥

## 1.0.0 (September 23, 2021)

### Initial Release

#### RESOURCE FEATURES

- New Resource: ``resource_zpa_app_server_controller`` 🆕
- New Resource: ``resource_zpa_application_segment`` 🆕
- New Resource: ``resource_zpa_browser_access`` 🆕
- New Resource: ``resource_zpa_policy_access_forwarding_rule`` 🆕
- New Resource: ``resource_zpa_policy_access_rule`` 🆕
- New Resource: ``resource_zpa_policy_access_timeout_rule`` 🆕
- New Resource: ``resource_zpa_segment_group`` 🆕
- New Resource: ``resource_zpa_server_group`` 🆕

DATA SOURCE FEATURES

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
