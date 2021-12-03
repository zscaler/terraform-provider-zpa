## 2.0.4 (December 3, 2021)

## Enhancement

- The provider now supports the ability to import resources via its `name` and/or `id` property to support easier migration of existing ZPA resources via `terraform import` command.
This capability is currently available to the following resources:
**resource_zpa_app_connector_group** - Issue [[#29](https://github.com/willguibr/terraform-provider-zpa/issues/29)]
**resource_zpa_app_server_controller** - :rocket:
**resource_zpa_application_segment** - :rocket:
**resource_zpa_segment_group** - :rocket:
**resource_zpa_server_group** - :rocket:
**resource_zpa_service_edge_group** - :rocket:

Note: To import resources not currently supported, the resource numeric ID is required.

# 2.0.3 (November 21, 2021)

## Dependabot Updates

- Dependabot updates [PR#33](https://github.com/willguibr/terraform-provider-zpa/pull/33/) Bump github.com/hashicorp/terraform-plugin-docs from 0.5.0 to 0.5.1 #33
- Dependabot updates [PR#34](https://github.com/willguibr/terraform-provider-zpa/pull/34) Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.8.0 to 2.9.0

## 2.0.2 (November 7, 2021)

## Enhancement

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
