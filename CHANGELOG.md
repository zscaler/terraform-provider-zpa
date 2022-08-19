# Changelog

## 2.3.1

### Notes

- Release date: **(August 30 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #169](https://github.com/zscaler/terraform-provider-zpa/pull/161) Fixed policy rule order, where the rule order in the UI didn't correspond to the desired order set in HCL. Issue [[#166](https://github.com/zscaler/terraform-provider-zpa/issues/166)]
- [PR #170](https://github.com/zscaler/terraform-provider-zpa/pull/161) Fixed special character encoding, where certain symbols caused Terraform to indicate potential configuration drifts. Issue [[#149](https://github.com/zscaler/terraform-provider-zpa/issues/149)]

## 2.3.0

### Notes

- Release date: **(August 17 2022)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #161](https://github.com/zscaler/terraform-provider-zpa/pull/161) Integrated newly created Zscaler GO SDK. Models are now centralized in the repository [zscaler-sdk-go](https://github.com/zscaler/zscaler-sdk-go)

## 2.2.2

### Notes

- Release date: **(July 19 2022)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #159](https://github.com/zscaler/terraform-provider-zpa/pull/159) Added Terraform UserAgent for Backend API tracking

## 2.2.1

### Notes

- Release date: **(July 6 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- Fix: Fixed authentication mechanism variables for ZPA Beta and GOV

### Documentation

1. Fixed application segment documentation and examples

## 2.2.0

### Notes

- Supported Terraform version: **v1.x**

### New Features

1. The provider now supports the following ZPA Privileged Remote Access (PRA) features:

**zpa_application_segment_pra** - The resource supports enabling `SECURE_REMOTE_ACCESS` for RDP and SSH via the `app_types` parameter. PR [#133](https://github.com/zscaler/terraform-provider-zpa/pull/133)

2. The provider now supports the following ZPA Inspection features:
**zpa_inspection_custom_controls** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
**zpa_inpection_predefined_controls** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
**zpa_inspection_profile** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
**zpa_policy_access_inspection_rule** PR[#134](https://github.com/zscaler/terraform-provider-zpa/pull/134)
**zpa_application_segment_inspection** - The resource supports enabling `INSPECT` for HTTP and HTTPS via the `app_types` parameter. PR [#135](https://github.com/zscaler/terraform-provider-zpa/pull/135)

4. Implemented a new Application Segment resource parameter ``select_connector_close_to_app``. The parameter can only be set for TCP based applications. PR [#137](https://github.com/zscaler/terraform-provider-zpa/pull/137)

### Enhancements

- Added support to `scim_attribute_header` to support policy access SCIM criteria based on SCIM attribute values.  Issue [#146](https://github.com/zscaler/terraform-provider-zpa/issues/146) / PR [#147]((https://github.com/zscaler/terraform-provider-zpa/pull/147))

- ZPA Beta Cloud: The provider now supports authentication via environment variables or static credentials to ZPA Beta Cloud. For authentication instructions please refer to the documentation page [here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.md) PR [#136](https://github.com/zscaler/terraform-provider-zpa/pull/136)

- ZPA Gov Cloud: The provider now supports authentication via environment variables or static credentials to ZPA Gov Cloud. For authentication instructions please refer to the documentation page [here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.md) PR [#145](https://github.com/zscaler/terraform-provider-zpa/pull/145)

### Bug Fixes
- Fix: Fixed update function on **zpa_app_server_controller** resource to ensure desired state is enforced in the upstream resource. Issue [#128](https://github.com/zscaler/terraform-provider-zpa/issues/128)
- Fix: Fixed Golangci linter

### Documentation
1. Added release notes guide to documentation PR #140
2. Fixed documentation misspellings

## 2.1.5 (May, 18 2022)

### Notes

- Supported Terraform version: **v1.x**

### Annoucements

The Terraform Provider for Zscaler Private Access (ZPA) is now officially hosted under Zscaler's GitHub account and published in the Terraform Registry. For more details, visit the Zscaler Community Article [Here](https://community.zscaler.com/t/zpa-and-zia-terraform-providers-now-verified/16675)
Administrators who used previous versions of the provider, and followed instructions to install the binary as a custom provider, must update their provider block as such:

```hcl
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

### Enhancements

- Documentation: Updated documentation to comply with Terraform registry formatting. #125
- ``zpa_posture_profile`` Updated search mechanism to support posture profile name search without the Zscaler cloud name. PR #123
- ``zpa_trusted_network`` Updated search mechanism to support trusted network name search without the Zscaler cloud name. PR #123

### Bug Fixes

- Fixed ``zpa_application_segment`` to support updates on ``tcp_port_ranges``, ``udp_port_ranges`` and ``tcp_port_range``, ``udp_port_range`` Issue #103

## 2.1.3 (May, 18 2022)

### Notes

- Supported Terraform version: **v1.x**

### Announcements

The Terraform Provider for Zscaler Private Access (ZPA) is now officially hosted under Zscaler's GitHub account and published in the Terraform Registry. For more details, visit the Zscaler Community Article [Here](https://community.zscaler.com/t/zpa-and-zia-terraform-providers-now-verified/16675)
Administrators who used previous versions of the provider, and followed instructions to install the binary as a custom provider, must update their provider block as such:

```hcl
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

## 2.1.2 (May 6, 2022)

### Notes

- Supported Terraform version: **v1.x**

### Bug Fixes

- Fix: tcp and udp ports were not being updated during changes, requiring the application segment resource to be fully destroyed and rebuilt. Implemented ``ForceNew`` in the the ``zpa_application_segment`` resource parameters: ``tcp_port_range``, ``udp_port_range``, ``tcp_port_ranges``, ``udp_port_ranges``. This behavior instructs Terraform to first destroy and then recreate the resource if any of the attributes change in the configuration, as opposed to trying to update the existing resource. The destruction of the resource does not impact attached resources such as server groups, segment groups or policies.

## 2.1.1 (April 27, 2022)

### Notes

- Supported Terraform version: **v1.x**

### Enhancements

1. Refactored and added new acceptance tests for better statement coverage. These tests are considered best practice and were added to routinely verify that the ZPA Terraform Plugin produces the expected outcome. [PR#88], [PR#96], [PR#98], [PR#99]

2. Support explicitly empty port ranges. Allow optional use of Attributes as Blocks syntax for ``zpa_application_segment`` {tcp,udp}_port_range blocks, allowing clean specification of "no port ranges" in dynamic contexts. [PR#97](https://github.com/zscaler/terraform-provider-zpa/pull/97) Thanks @isometry

### Deprecations

1. Deprecated all legacy policy set controller endpoints: ``/policySet/global``, ``/policySet/reauth``, ``/policySet/bypass`` [PR#88](https://github.com/zscaler/terraform-provider-zpa/pull/88)

2. Deprecated all references to ZPA private API gateway. [PR#87](https://github.com/zscaler/terraform-provider-zpa/pull/87)

## 2.1.0 (March 05, 2022)

### Notes

- Supported Terraform version: **v1.x**

### Enhancements

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

## 2.0.7 (February 17, 2022)

### Notes

- Supported Terraform version: **v1.x**

### Bug Fixes

- ZPA-50: Fixed and removed deprecated arguments from ``zpa_application_segments`` data source and resource :wrench:
- ZPA-50: Fixed ``zpa_posture_profile`` and ``zpa_trusted_networks`` acceptance tests to include ZIA cloud name :wrench:

### Enhancements

- ZPA-51: Updated common ``NetworkPorts`` flatten and expand functions for better optimization and global use across multiple application segment resources. This update affects the following resources: ``data_source_zpa_application_segment``, ``data_source_zpa_browser_access`` and ``resource_zpa_application_segment``, ``resource_source_zpa_browser_access`` :rocket:

## 2.0.6 (February 3, 2022)

### Notes

- Supported Terraform version: **v1.x**

### New Data Sources

- Added new data source for ``zpa_app_connector_controller`` resource. [PR#62](https://github.com/zscaler/terraform-provider-zpa/pull/62)
- Added new data source for ``zpa_service_edge_controller`` resource. [PR#63](https://github.com/zscaler/terraform-provider-zpa/pull/63)

### New Acceptance Tests

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

## 2.0.5 (December 20, 2021)

### Notes

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

### Bug Fixes

- Fixed pagination issues with all resources where only the default pagesize was being returned. [PR#52](https://github.com/zscaler/terraform-provider-zpa/pull/52) :wrench:
- Fixed issue where Terraform showed that resources had been modified even though nothing had been changed in the upstream resources.[PR#54](https://github.com/zscaler/terraform-provider-zpa/pull/54) :wrench:

## 2.0.4 (December 6, 2021)

### Notes

- Supported Terraform version: **v1.x**

### New Data Sources

- Added new data source for ``zpa_browser_access`` resource.

### Enhancements

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

### Bug Fixes

- Fixed [INFO] and [Error] message in ``data_source_zpa_lss_config_controller`` [PR#43](https://github.com/zscaler/terraform-provider-zpa/pull/43) 🔧

## 2.0.3 (November 21, 2021)

### Notes

- Supported Terraform version: **v1.x**

### Dependabot Updates

- Dependabot updates [PR#33](https://github.com/zscaler/terraform-provider-zpa/pull/33/) Bump github.com/hashicorp/terraform-plugin-docs from 0.5.0 to 0.5.1 #33
- Dependabot updates [PR#34](https://github.com/zscaler/terraform-provider-zpa/pull/34) Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.8.0 to 2.9.0

## 2.0.2 (November 7, 2021)

### Notes

- Supported Terraform version: **v1.x**

### Enhancements

- Added custom validation function ``ValidateStringFloatBetween`` to ``resource_zpa_app_connector_group`` to validate ``longitude`` and ``latitude`` parameters. [ZPA-17](https://github.com/zscaler/terraform-provider-zpa/pull/17).
- Added custom validation function ``ValidateStringFloatBetween`` to ``resource_zpa_service_edge_group`` to validate ``longitude`` and ``latitude`` parameters. [ZPA-18](https://github.com/zscaler/terraform-provider-zpa/pull/18).

## 2.0.1 (November 4, 2021)

### Notes

- Supported Terraform version: **v1.x**

### Bug Fixes

- Fixed issue where provider authentication parameters for hard coded credentials was not working:
- Changed the following variable names: ``client_id``, ``client_secret`` and ``customerid`` to ``zpa_client_id``, ``zpa_client_secret`` and ``zpa_customer_id``.

## 2.0.0 (November 3, 2021)

### Notes

- Supported Terraform version: **v1.x**

- New management APIs are now available to manage App Connectors, App Connector Groups, Service Edges, Service Edge Groups, and Log Streaming Service (LSS) configurations.
- New prerequisite APIs for enrollment certificates, provisioning keys, and to get version profiles, client types, status codes, and LSS formats are added.
- A new API to reorder policy rules is added.
- The endpoints to get all browser access (BA) certificates, IdPs, posture profiles, trusted networks, and SAML attributes are now deprecated, and new APIs with pagination are provided.
- API endpoints specific to a policy (global/reauth/bypass) are deprecated and replaced by a generic API that takes policyType as a parameter.
- The port range configuration for the application segment has been enhanced for more readability. The tcpPortRanges and udpPortRanges fields are deprecated and replaced with tcpPortRange and udpPortRange.

### Features

### New Resources

- New Resource: ``resource_zpa_app_connector_group`` 🆕
- New Resource: ``resource_zpa_service_edge_group`` 🆕
- New Resource: ``resource_zpa_provisioning_key`` 🆕
- New Resource: ``resource_zpa_lss_config_controller`` 🆕

### New Data Sources

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

### Notes

- Supported Terraform version: **v1.x**

### Initial Release

#### Resource Features

- New Resource: ``resource_zpa_app_server_controller`` 🆕
- New Resource: ``resource_zpa_application_segment`` 🆕
- New Resource: ``resource_zpa_browser_access`` 🆕
- New Resource: ``resource_zpa_policy_access_forwarding_rule`` 🆕
- New Resource: ``resource_zpa_policy_access_rule`` 🆕
- New Resource: ``resource_zpa_policy_access_timeout_rule`` 🆕
- New Resource: ``resource_zpa_segment_group`` 🆕
- New Resource: ``resource_zpa_server_group`` 🆕

### Data Source Features

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
